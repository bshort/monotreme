package frontend

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/bshort/monotreme/internal/util"
	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/server/common"
	"github.com/bshort/monotreme/server/profile"
	"github.com/bshort/monotreme/store"
)

//go:embed dist
var embeddedFiles embed.FS

const (
	headerMetadataPlaceholder = "<!-- monotreme.metadata -->"
)

type FrontendService struct {
	Profile *profile.Profile
	Store   *store.Store
}

func NewFrontendService(profile *profile.Profile, store *store.Store) *FrontendService {
	return &FrontendService{
		Profile: profile,
		Store:   store,
	}
}

func (s *FrontendService) getShortcutPrefix(ctx context.Context) string {
	shortcutRelatedSetting, err := s.Store.GetWorkspaceSetting(ctx, &store.FindWorkspaceSetting{
		Key: storepb.WorkspaceSettingKey_WORKSPACE_SETTING_SHORTCUT_RELATED,
	})
	if err != nil || shortcutRelatedSetting == nil {
		return "s" // Default fallback
	}
	prefix := shortcutRelatedSetting.GetShortcutRelated().GetShortcutPrefix()
	if prefix == "" {
		return "s" // Default fallback
	}
	return prefix
}

func (s *FrontendService) isShortcutRoute(c echo.Context) bool {
	path := c.Path()
	if path == "/" {
		return false
	}

	// Split path into segments
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) != 2 {
		return false
	}

	// Check if the first segment matches the current shortcut prefix
	ctx := c.Request().Context()
	currentPrefix := s.getShortcutPrefix(ctx)
	isShortcut := segments[0] == currentPrefix

	// Log for debugging
	slog.Info("isShortcutRoute check",
		slog.String("path", path),
		slog.String("prefix", segments[0]),
		slog.String("current", currentPrefix),
		slog.Bool("isShortcut", isShortcut))

	return isShortcut
}

func (s *FrontendService) Serve(ctx context.Context, e *echo.Echo) {
	rawIndexHTML := getRawIndexHTML()

	// Add middleware to handle shortcut/collection routes BEFORE static middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()
			urlPath := c.Request().URL.Path
			method := c.Request().Method
			requestURI := c.Request().RequestURI
			slog.Info("Middleware check", slog.String("path", path), slog.String("urlPath", urlPath), slog.String("method", method))

			// Add debug header to test if middleware is working
			c.Response().Header().Set("X-Debug-Middleware", "active")
			c.Response().Header().Set("X-Debug-Path", path)
			c.Response().Header().Set("X-Debug-URL-Path", urlPath)
			c.Response().Header().Set("X-Debug-Method", method)
			c.Response().Header().Set("X-Debug-RequestURI", requestURI)

			// Use URL.Path if c.Path() is empty
			if path == "" {
				path = urlPath
			}

			// Only handle GET requests
			if method != "GET" {
				c.Response().Header().Set("X-Debug-Skip", "non-GET")
				return next(c)
			}

			// Split path into segments
			segments := strings.Split(strings.Trim(path, "/"), "/")
			c.Response().Header().Set("X-Debug-Segments", fmt.Sprintf("%d", len(segments)))
			if len(segments) != 2 {
				c.Response().Header().Set("X-Debug-Skip", "not-2-segments")
				return next(c)
			}

			prefix := segments[0]
			name := segments[1]
			c.Response().Header().Set("X-Debug-Prefix", prefix)
			c.Response().Header().Set("X-Debug-Name", name)
			ctx := c.Request().Context()

			// Handle collection routes
			if prefix == "c" {
				collection, err := s.Store.GetCollection(ctx, &store.FindCollection{
					Name: &name,
				})
				if err == nil && collection != nil {
					indexHTML := strings.ReplaceAll(rawIndexHTML, headerMetadataPlaceholder, generateCollectionMetadata(collection).String())
					return c.HTML(http.StatusOK, indexHTML)
				}
			}

			// Handle shortcut routes
			currentPrefix := s.getShortcutPrefix(ctx)
			slog.Info("Checking shortcut middleware", slog.String("current", currentPrefix), slog.String("requested", prefix))
			c.Response().Header().Set("X-Debug-Current-Prefix", currentPrefix)

			if prefix == currentPrefix {
				c.Response().Header().Set("X-Debug-Prefix-Match", "true")
				shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
					Name: &name,
				})
				c.Response().Header().Set("X-Debug-Shortcut-Error", fmt.Sprintf("%v", err))
				if err == nil && shortcut != nil {
					c.Response().Header().Set("X-Debug-Shortcut-Found", "true")
					// Create shortcut view activity.
					if err := s.createShortcutViewActivity(ctx, c.Request(), shortcut); err != nil {
						slog.Warn("failed to create shortcut view activity", slog.String("error", err.Error()))
					}

					// Copy query parameters to the target URL
					targetURL := shortcut.Link
					if c.Request().URL.RawQuery != "" {
						separator := "?"
						if strings.Contains(targetURL, "?") {
							separator = "&"
						}
						targetURL = fmt.Sprintf("%s%s%s", targetURL, separator, c.Request().URL.RawQuery)
					}

					// Redirect to the shortcut's target URL
					return c.Redirect(http.StatusFound, targetURL)
				} else {
					c.Response().Header().Set("X-Debug-Shortcut-Found", "false")
					// Log attempted access to non-existent shortcut
					if err := s.createShortcutNotFoundActivity(ctx, c.Request(), name); err != nil {
						slog.Warn("failed to create shortcut not found activity", slog.String("error", err.Error()))
					}
					slog.Info("Shortcut not found", slog.String("name", name), slog.String("path", path))
				}
			} else {
				c.Response().Header().Set("X-Debug-Prefix-Match", "false")
			}

			// Continue to next middleware (static files)
			return next(c)
		}
	})

	// Use echo static middleware to serve the built dist folder.
	// Reference: https://github.com/labstack/echo/blob/master/middleware/static.go
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Filesystem: getFileSystem("dist"),
		Skipper: func(c echo.Context) bool {
			// Skip static serving for API routes and shortcut routes
			path := c.Path()
			if util.HasPrefixes(path, "/api", "/monotreme.api.v1") {
				return true
			}

			// Skip static serving for potential shortcut/collection routes
			segments := strings.Split(strings.Trim(path, "/"), "/")
			if len(segments) == 2 {
				prefix := segments[0]
				// Check if this could be a shortcut route (s prefix) or collection route (c prefix)
				if prefix == "s" || prefix == "c" {
					return true
				}
			}

			return false
		},
	}))

	assetsGroup := e.Group("assets")
	// Use echo gzip middleware to compress the response.
	// Reference: https://echo.labstack.com/docs/middleware/gzip
	assetsGroup.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			// Skip gzip for API routes
			return util.HasPrefixes(c.Path(), "/api", "/monotreme.api.v1")
		},
		Level: 5,
	}))
	assetsGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderCacheControl, "max-age=31536000, immutable")
			return next(c)
		}
	})
	assetsGroup.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Filesystem: getFileSystem("dist/assets"),
		Skipper: func(c echo.Context) bool {
			// Skip static serving for API routes
			return util.HasPrefixes(c.Path(), "/api", "/monotreme.api.v1")
		},
	}))
}

// Routes are now handled by middleware in the Serve method

func (s *FrontendService) createShortcutViewActivity(ctx context.Context, request *http.Request, shortcut *storepb.Shortcut) error {
	ip := getReadUserIP(request)
	referer := request.Header.Get("Referer")
	userAgent := request.Header.Get("User-Agent")
	params := map[string]*storepb.ActivityShorcutViewPayload_ValueList{}
	for key, values := range request.URL.Query() {
		params[key] = &storepb.ActivityShorcutViewPayload_ValueList{Values: values}
	}
	payload := &storepb.ActivityShorcutViewPayload{
		ShortcutId: shortcut.Id,
		Ip:         ip,
		Referer:    referer,
		UserAgent:  userAgent,
		Params:     params,
	}
	payloadStr, err := protojson.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal activity payload")
	}
	activity := &store.Activity{
		CreatorID: common.BotID,
		Type:      store.ActivityShortcutView,
		Level:     store.ActivityInfo,
		Payload:   string(payloadStr),
	}
	_, err = s.Store.CreateActivity(ctx, activity)
	if err != nil {
		return errors.Wrap(err, "Failed to create activity")
	}
	return nil
}

func (s *FrontendService) createShortcutNotFoundActivity(ctx context.Context, request *http.Request, shortcutName string) error {
	ip := getReadUserIP(request)
	referer := request.Header.Get("Referer")
	userAgent := request.Header.Get("User-Agent")
	params := map[string]*storepb.ActivityShorcutViewPayload_ValueList{}
	for key, values := range request.URL.Query() {
		params[key] = &storepb.ActivityShorcutViewPayload_ValueList{Values: values}
	}

	// Create a payload for shortcut not found activity
	// We'll use a similar structure but with ShortcutId = 0 to indicate not found
	payload := &storepb.ActivityShorcutViewPayload{
		ShortcutId: 0, // 0 indicates shortcut not found
		Ip:         ip,
		Referer:    referer,
		UserAgent:  userAgent,
		Params:     params,
	}
	payloadStr, err := protojson.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal shortcut not found activity payload")
	}

	// Create activity with a custom message indicating the shortcut name that was attempted
	activity := &store.Activity{
		CreatorID: common.BotID,
		Type:      store.ActivityShortcutView, // Reuse same type for consistency
		Level:     store.ActivityWarn,         // Use Warn level to distinguish from successful visits
		Payload:   string(payloadStr),
	}
	_, err = s.Store.CreateActivity(ctx, activity)
	if err != nil {
		return errors.Wrap(err, "Failed to create shortcut not found activity")
	}
	return nil
}

func getReadUserIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func getFileSystem(path string) http.FileSystem {
	fs, err := fs.Sub(embeddedFiles, path)
	if err != nil {
		panic(err)
	}

	return http.FS(fs)
}

func generateShortcutMetadata(shortcut *storepb.Shortcut) *Metadata {
	metadata := getDefaultMetadata()
	title, description := shortcut.Title, shortcut.Description
	if shortcut.OgMetadata != nil {
		if shortcut.OgMetadata.Title != "" {
			title = shortcut.OgMetadata.Title
		}
		if shortcut.OgMetadata.Description != "" {
			description = shortcut.OgMetadata.Description
		}
		metadata.ImageURL = shortcut.OgMetadata.Image
	}
	metadata.Title = title
	metadata.Description = description
	return metadata
}

func generateCollectionMetadata(collection *storepb.Collection) *Metadata {
	metadata := getDefaultMetadata()
	metadata.Title = collection.Title
	metadata.Description = collection.Description
	return metadata
}

func getRawIndexHTML() string {
	bytes, _ := embeddedFiles.ReadFile("dist/index.html")
	return string(bytes)
}

type Metadata struct {
	Title       string
	Description string
	ImageURL    string
}

func getDefaultMetadata() *Metadata {
	return &Metadata{
		Title: "Monotreme",
	}
}

func (m *Metadata) String() string {
	metadataList := []string{
		fmt.Sprintf(`<title>%s</title>`, m.Title),
		fmt.Sprintf(`<meta name="description" content="%s" />`, m.Description),
		fmt.Sprintf(`<meta property="og:title" content="%s" />`, m.Title),
		fmt.Sprintf(`<meta property="og:description" content="%s" />`, m.Description),
		fmt.Sprintf(`<meta property="og:image" content="%s" />`, m.ImageURL),
		`<meta property="og:type" content="website" />`,
		// Twitter related fields.
		fmt.Sprintf(`<meta property="twitter:title" content="%s" />`, m.Title),
		fmt.Sprintf(`<meta property="twitter:description" content="%s" />`, m.Description),
		fmt.Sprintf(`<meta property="twitter:image" content="%s" />`, m.ImageURL),
	}
	return strings.Join(metadataList, "\n")
}
