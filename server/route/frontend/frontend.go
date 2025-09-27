package frontend

import (
	"context"
	"embed"
	"fmt"
	"html"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
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

	// Add route for public user shortcuts display
	e.GET("/:username/shortcuts", s.handlePublicShortcuts)

	// Add route for public user collections display
	e.GET("/:username/collections", s.handlePublicCollections)

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

func (s *FrontendService) handlePublicShortcuts(c echo.Context) error {
	ctx := c.Request().Context()
	username := c.Param("username")
	filter := c.QueryParam("filter") // Get filter parameter: "public", "private", "all"

	// Find user by nickname
	users, err := s.Store.ListUsers(ctx, &store.FindUser{
		Nickname: &username,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to find user",
		})
	}
	if len(users) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	user := users[0]

	// Determine which shortcuts to fetch based on filter
	var visibilityList []storepb.Visibility
	var filterDescription string

	switch filter {
	case "private":
		visibilityList = []storepb.Visibility{storepb.Visibility_WORKSPACE}
		filterDescription = "Private shortcuts"
	case "all":
		visibilityList = []storepb.Visibility{storepb.Visibility_PUBLIC, storepb.Visibility_WORKSPACE}
		filterDescription = "All shortcuts"
	case "public":
		fallthrough
	default:
		visibilityList = []storepb.Visibility{storepb.Visibility_PUBLIC}
		filterDescription = "Public shortcuts"
		filter = "public" // Normalize default case
	}

	// Get shortcuts for this user based on filter
	shortcuts, err := s.Store.ListShortcuts(ctx, &store.FindShortcut{
		CreatorID: &user.ID,
		VisibilityList: visibilityList,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch shortcuts",
		})
	}

	// Create HTML response
	html := s.generatePublicShortcutsHTML(username, shortcuts, filter, filterDescription)
	return c.HTML(http.StatusOK, html)
}

func (s *FrontendService) generatePublicShortcutsHTML(username string, shortcuts []*storepb.Shortcut, currentFilter string, filterDescription string) string {
	ctx := context.Background() // We need context for getShortcutPrefix
	escapedUsername := html.EscapeString(username)
	htmlContent := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + escapedUsername + `'s Shortcuts</title>
    <style>
        body {
            font-family: system-ui, -apple-system, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 2rem;
            background-color: #f8fafc;
        }
        .header {
            margin-bottom: 2rem;
            text-align: center;
        }
        .header h1 {
            color: #1e293b;
            margin-bottom: 0.5rem;
        }
        .header p {
            color: #64748b;
        }
        .filter-container {
            display: flex;
            justify-content: center;
            gap: 0.5rem;
            margin-bottom: 2rem;
            flex-wrap: wrap;
        }
        .filter-button {
            padding: 0.5rem 1rem;
            border: 2px solid #e2e8f0;
            background: white;
            color: #64748b;
            text-decoration: none;
            border-radius: 8px;
            font-weight: 500;
            transition: all 0.2s ease;
            display: inline-block;
        }
        .filter-button:hover {
            border-color: #3b82f6;
            color: #3b82f6;
            text-decoration: none;
            transform: translateY(-1px);
        }
        .filter-button.active {
            background: #3b82f6;
            border-color: #3b82f6;
            color: white;
        }
        .filter-button.active:hover {
            background: #2563eb;
            border-color: #2563eb;
            color: white;
            transform: translateY(-1px);
        }
        .shortcuts-grid {
            display: grid;
            gap: 1rem;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        }
        .shortcut-card {
            background: white;
            border-radius: 8px;
            padding: 1.5rem;
            border: 1px solid #e2e8f0;
            transition: all 0.2s ease;
            display: flex;
            gap: 1rem;
            text-decoration: none;
            color: inherit;
        }
        .shortcut-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 25px rgba(0,0,0,0.1);
            border-color: #3b82f6;
            text-decoration: none;
            color: inherit;
        }
        .shortcut-favicon {
            width: 32px;
            height: 32px;
            flex-shrink: 0;
            border-radius: 4px;
            background: #f1f5f9;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .shortcut-favicon img {
            width: 24px;
            height: 24px;
            border-radius: 2px;
        }
        .shortcut-favicon-placeholder {
            width: 24px;
            height: 24px;
            background: #cbd5e1;
            border-radius: 2px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 12px;
            color: #64748b;
            font-weight: 600;
        }
        .shortcut-favicon-placeholder.hidden {
            display: none;
        }
        .shortcut-content {
            flex: 1;
            min-width: 0;
        }
        .shortcut-title {
            font-size: 1.1rem;
            font-weight: 600;
            color: #1e293b;
            margin-bottom: 0.5rem;
        }
        .shortcut-description {
            font-size: 0.9rem;
            color: #64748b;
            margin-bottom: 0.75rem;
            line-height: 1.5;
        }
        .shortcut-link {
            font-size: 0.85rem;
            color: #3b82f6;
            word-break: break-all;
            display: block;
            margin-bottom: 0.75rem;
        }
        .shortcut-name {
            font-size: 0.85rem;
            font-weight: 500;
            color: #475569;
            background: #f8fafc;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            display: inline-block;
        }
        .tags {
            margin-top: 0.5rem;
        }
        .tag {
            display: inline-block;
            background: #e0e7ff;
            color: #3730a3;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            font-size: 0.75rem;
            margin-right: 0.5rem;
            margin-bottom: 0.25rem;
        }
        .no-shortcuts {
            text-align: center;
            color: #64748b;
            font-style: italic;
            grid-column: 1 / -1;
            padding: 3rem;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>` + escapedUsername + `'s Shortcuts</h1>
        <p>` + html.EscapeString(filterDescription) + `</p>
    </div>

    <div class="filter-container">`

	// Add filter buttons with active state
	if currentFilter == "public" {
		htmlContent += `<a href="?filter=public" class="filter-button active">Public</a>`
	} else {
		htmlContent += `<a href="?filter=public" class="filter-button">Public</a>`
	}

	if currentFilter == "private" {
		htmlContent += `<a href="?filter=private" class="filter-button active">Private</a>`
	} else {
		htmlContent += `<a href="?filter=private" class="filter-button">Private</a>`
	}

	if currentFilter == "all" {
		htmlContent += `<a href="?filter=all" class="filter-button active">All</a>`
	} else {
		htmlContent += `<a href="?filter=all" class="filter-button">All</a>`
	}

	htmlContent += `</div>

    <div class="shortcuts-grid">`

	if len(shortcuts) == 0 {
		emptyMessage := "No shortcuts available"
		switch currentFilter {
		case "public":
			emptyMessage = "No public shortcuts available"
		case "private":
			emptyMessage = "No private shortcuts available"
		case "all":
			emptyMessage = "No shortcuts available"
		}
		htmlContent += `<div class="no-shortcuts">` + emptyMessage + `</div>`
	} else {
		shortcutPrefix := s.getShortcutPrefix(ctx)
		for _, shortcut := range shortcuts {
			shortcutURL := fmt.Sprintf("/%s/%s", shortcutPrefix, shortcut.Name)
			htmlContent += `<a href="` + html.EscapeString(shortcutURL) + `" target="_blank" class="shortcut-card">
                <div class="shortcut-favicon">`

			// Check if there's a custom icon or use favicon from the domain
			placeholder := "?"
			if len(shortcut.Name) > 0 {
				placeholder = strings.ToUpper(string([]rune(shortcut.Name)[0]))
			}

			if shortcut.CustomIcon != "" {
				htmlContent += `<img src="` + html.EscapeString(shortcut.CustomIcon) + `" alt="Icon" onerror="this.style.display='none'; this.nextSibling.classList.remove('hidden');">
					<div class="shortcut-favicon-placeholder hidden">` + placeholder + `</div>`
			} else {
				// Try to get favicon from the domain
				faviconURL := s.getFaviconURL(shortcut.Link)
				if faviconURL != "" {
					htmlContent += `<img src="` + html.EscapeString(faviconURL) + `" alt="Favicon" onerror="this.style.display='none'; this.nextSibling.classList.remove('hidden');">
					<div class="shortcut-favicon-placeholder hidden">` + placeholder + `</div>`
				} else {
					htmlContent += `<div class="shortcut-favicon-placeholder">` + placeholder + `</div>`
				}
			}

			htmlContent += `</div>
                <div class="shortcut-content">`

			// 1. Title
			if shortcut.Title != "" {
				htmlContent += `<div class="shortcut-title">` + html.EscapeString(shortcut.Title) + `</div>`
			}

			// 2. Description
			if shortcut.Description != "" {
				htmlContent += `<div class="shortcut-description">` + html.EscapeString(shortcut.Description) + `</div>`
			}

			// 3. Link (display only, no longer clickable since the whole card is clickable)
			htmlContent += `<div class="shortcut-link">` + html.EscapeString(shortcut.Link) + `</div>`

			// 4. Shortcut name (at the bottom)
			htmlContent += `<div class="shortcut-name">` + html.EscapeString(shortcut.Name) + `</div>`

			// Tags
			if len(shortcut.Tags) > 0 {
				htmlContent += `<div class="tags">`
				for _, tag := range shortcut.Tags {
					if tag != "" {
						htmlContent += `<span class="tag">` + html.EscapeString(tag) + `</span>`
					}
				}
				htmlContent += `</div>`
			}

			htmlContent += `</div></a>`
		}
	}

	htmlContent += `    </div>
</body>
</html>`
	return htmlContent
}

func (s *FrontendService) handlePublicCollections(c echo.Context) error {
	ctx := c.Request().Context()
	username := c.Param("username")

	// Find user by nickname
	users, err := s.Store.ListUsers(ctx, &store.FindUser{
		Nickname: &username,
	})
	if err != nil {
		slog.Error("Failed to find user", "error", err, "username", username)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to find user",
		})
	}
	if len(users) == 0 {
		slog.Info("User not found", "username", username)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	user := users[0]
	slog.Info("Found user", "username", username, "userID", user.ID)

	// Get ALL collections for this user first to see what's available
	allCollections, err := s.Store.ListCollections(ctx, &store.FindCollection{
		CreatorID: &user.ID,
	})
	if err != nil {
		slog.Error("Failed to fetch all collections", "error", err, "userID", user.ID)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch collections",
		})
	}
	slog.Info("Found collections for user", "userID", user.ID, "totalCollections", len(allCollections))

	// Get ALL collections for this user (temporarily for debugging)
	collections, err := s.Store.ListCollections(ctx, &store.FindCollection{
		CreatorID: &user.ID,
		// VisibilityList: []storepb.Visibility{storepb.Visibility_PUBLIC}, // Commented out for debugging
	})
	if err != nil {
		slog.Error("Failed to fetch collections", "error", err, "userID", user.ID)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch collections",
		})
	}
	slog.Info("Found all collections for user", "userID", user.ID, "collections", len(collections))

	// For each collection, get the public shortcuts
	collectionsWithShortcuts := make([]CollectionWithShortcuts, 0)
	for _, collection := range collections {
		slog.Info("Processing collection", "collectionID", collection.Id, "collectionTitle", collection.Title, "shortcutIDs", collection.ShortcutIds)
		shortcuts := make([]*storepb.Shortcut, 0)

		// Get shortcuts by IDs and filter for public ones
		for _, shortcutID := range collection.ShortcutIds {
			shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
				ID: &shortcutID,
			})
			if err != nil {
				slog.Warn("Failed to get shortcut", "shortcutID", shortcutID, "error", err)
				continue // Skip if shortcut not found
			}
			slog.Info("Found shortcut", "shortcutID", shortcutID, "name", shortcut.Name, "visibility", shortcut.Visibility.String())
			// Include ALL shortcuts for debugging (normally would filter for public only)
			shortcuts = append(shortcuts, shortcut)
			slog.Info("Added shortcut to collection", "shortcutName", shortcut.Name, "visibility", shortcut.Visibility.String())
		}

		slog.Info("Collection processed", "collectionTitle", collection.Title, "publicShortcuts", len(shortcuts))
		collectionsWithShortcuts = append(collectionsWithShortcuts, CollectionWithShortcuts{
			Collection: collection,
			Shortcuts:  shortcuts,
		})
	}

	// Create HTML response
	slog.Info("Generating HTML for collections", "username", username, "collectionsCount", len(collectionsWithShortcuts))
	html := s.generatePublicCollectionsHTML(username, collectionsWithShortcuts)
	return c.HTML(http.StatusOK, html)
}

type CollectionWithShortcuts struct {
	Collection *storepb.Collection
	Shortcuts  []*storepb.Shortcut
}

func (s *FrontendService) generatePublicCollectionsHTML(username string, collections []CollectionWithShortcuts) string {
	ctx := context.Background() // We need context for getShortcutPrefix
	escapedUsername := html.EscapeString(username)
	slog.Info("Generating HTML", "username", username, "collectionsToRender", len(collections))
	htmlContent := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + escapedUsername + `'s Collections (Debug)</title>
    <style>
        body {
            font-family: system-ui, -apple-system, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 2rem;
            background-color: #f8fafc;
        }
        .header {
            margin-bottom: 2rem;
            text-align: center;
        }
        .header h1 {
            color: #1e293b;
            margin-bottom: 0.5rem;
        }
        .header p {
            color: #64748b;
        }
        .collections-container {
            display: flex;
            flex-direction: column;
            gap: 3rem;
        }
        .collection-section {
            background: white;
            border-radius: 12px;
            padding: 2rem;
            border: 1px solid #e2e8f0;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .collection-header {
            margin-bottom: 2rem;
            padding-bottom: 1rem;
            border-bottom: 2px solid #e2e8f0;
        }
        .collection-title {
            font-size: 1.5rem;
            font-weight: 700;
            color: #1e293b;
            margin-bottom: 0.5rem;
        }
        .collection-description {
            color: #64748b;
            line-height: 1.6;
        }
        .shortcuts-grid {
            display: grid;
            gap: 1rem;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        }
        .shortcut-card {
            background: #f8fafc;
            border-radius: 8px;
            padding: 1.5rem;
            border: 1px solid #e2e8f0;
            transition: all 0.2s ease;
            display: flex;
            gap: 1rem;
            text-decoration: none;
            color: inherit;
        }
        .shortcut-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
            border-color: #3b82f6;
            text-decoration: none;
            color: inherit;
            background: white;
        }
        .shortcut-favicon {
            width: 32px;
            height: 32px;
            flex-shrink: 0;
            border-radius: 4px;
            background: #f1f5f9;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .shortcut-favicon img {
            width: 24px;
            height: 24px;
            border-radius: 2px;
        }
        .shortcut-favicon-placeholder {
            width: 24px;
            height: 24px;
            background: #cbd5e1;
            border-radius: 2px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 12px;
            color: #64748b;
            font-weight: 600;
        }
        .shortcut-favicon-placeholder.hidden {
            display: none;
        }
        .shortcut-content {
            flex: 1;
            min-width: 0;
        }
        .shortcut-title {
            font-size: 1.1rem;
            font-weight: 600;
            color: #1e293b;
            margin-bottom: 0.5rem;
        }
        .shortcut-description {
            font-size: 0.9rem;
            color: #64748b;
            margin-bottom: 0.75rem;
            line-height: 1.5;
        }
        .shortcut-link {
            font-size: 0.85rem;
            color: #3b82f6;
            word-break: break-all;
            display: block;
            margin-bottom: 0.75rem;
        }
        .shortcut-name {
            font-size: 0.85rem;
            font-weight: 500;
            color: #475569;
            background: #f1f5f9;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            display: inline-block;
        }
        .tags {
            margin-top: 0.5rem;
        }
        .tag {
            display: inline-block;
            background: #e0e7ff;
            color: #3730a3;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            font-size: 0.75rem;
            margin-right: 0.5rem;
            margin-bottom: 0.25rem;
        }
        .no-collections {
            text-align: center;
            color: #64748b;
            font-style: italic;
            padding: 3rem;
            background: white;
            border-radius: 12px;
            border: 1px solid #e2e8f0;
        }
        .no-shortcuts {
            text-align: center;
            color: #64748b;
            font-style: italic;
            padding: 2rem;
            background: #f8fafc;
            border-radius: 8px;
            border: 1px solid #e2e8f0;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>` + escapedUsername + `'s Collections (Debug)</h1>
        <p>All collections and shortcuts (debug mode)</p>
    </div>
    <div class="collections-container">`

	if len(collections) == 0 {
		htmlContent += `<div class="no-collections">No collections found for this user</div>`
	} else {
		shortcutPrefix := s.getShortcutPrefix(ctx)
		for _, collectionWithShortcuts := range collections {
			collection := collectionWithShortcuts.Collection
			shortcuts := collectionWithShortcuts.Shortcuts

			htmlContent += `<div class="collection-section">
                <div class="collection-header">
                    <div class="collection-title">` + html.EscapeString(collection.Title) + `</div>`

			if collection.Description != "" {
				htmlContent += `<div class="collection-description">` + html.EscapeString(collection.Description) + `</div>`
			}

			htmlContent += `</div><div class="shortcuts-grid">`

			if len(shortcuts) == 0 {
				htmlContent += `<div class="no-shortcuts">No public shortcuts in this collection</div>`
			} else {
				for _, shortcut := range shortcuts {
					shortcutURL := fmt.Sprintf("/%s/%s", shortcutPrefix, shortcut.Name)
					htmlContent += `<a href="` + html.EscapeString(shortcutURL) + `" target="_blank" class="shortcut-card">
                        <div class="shortcut-favicon">`

					// Check if there's a custom icon or use favicon from the domain
					placeholder := "?"
					if len(shortcut.Name) > 0 {
						placeholder = strings.ToUpper(string([]rune(shortcut.Name)[0]))
					}

					if shortcut.CustomIcon != "" {
						htmlContent += `<img src="` + html.EscapeString(shortcut.CustomIcon) + `" alt="Icon" onerror="this.style.display='none'; this.nextSibling.classList.remove('hidden');">
							<div class="shortcut-favicon-placeholder hidden">` + placeholder + `</div>`
					} else {
						// Try to get favicon from the domain
						faviconURL := s.getFaviconURL(shortcut.Link)
						if faviconURL != "" {
							htmlContent += `<img src="` + html.EscapeString(faviconURL) + `" alt="Favicon" onerror="this.style.display='none'; this.nextSibling.classList.remove('hidden');">
							<div class="shortcut-favicon-placeholder hidden">` + placeholder + `</div>`
						} else {
							htmlContent += `<div class="shortcut-favicon-placeholder">` + placeholder + `</div>`
						}
					}

					htmlContent += `</div>
                        <div class="shortcut-content">`

					// 1. Title
					if shortcut.Title != "" {
						htmlContent += `<div class="shortcut-title">` + html.EscapeString(shortcut.Title) + `</div>`
					}

					// 2. Description
					if shortcut.Description != "" {
						htmlContent += `<div class="shortcut-description">` + html.EscapeString(shortcut.Description) + `</div>`
					}

					// 3. Link
					htmlContent += `<div class="shortcut-link">` + html.EscapeString(shortcut.Link) + `</div>`

					// 4. Shortcut name
					htmlContent += `<div class="shortcut-name">` + html.EscapeString(shortcut.Name) + `</div>`

					// Tags
					if len(shortcut.Tags) > 0 {
						htmlContent += `<div class="tags">`
						for _, tag := range shortcut.Tags {
							if tag != "" {
								htmlContent += `<span class="tag">` + html.EscapeString(tag) + `</span>`
							}
						}
						htmlContent += `</div>`
					}

					htmlContent += `</div></a>`
				}
			}

			htmlContent += `</div></div>`
		}
	}

	htmlContent += `    </div>
</body>
</html>`
	return htmlContent
}

func (s *FrontendService) getFaviconURL(link string) string {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return ""
	}

	// Only support http and https schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ""
	}

	// Use Google's favicon service which is reliable
	return fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=32", parsedURL.Host)
}

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
