package export

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/bshort/monotreme/internal/util"
	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/server/profile"
	"github.com/bshort/monotreme/store"
)

const (
	// Constants from auth.go for JWT token validation
	Issuer                  = "monotreme"
	KeyID                   = "v1"
	AccessTokenAudienceName = "user.access-token"
)

// ClaimsMessage represents JWT claims
type ClaimsMessage struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type ExportService struct {
	Profile *profile.Profile
	Store   *store.Store
	Secret  string
}

func NewExportService(profile *profile.Profile, store *store.Store, secret string) *ExportService {
	return &ExportService{
		Profile: profile,
		Store:   store,
		Secret:  secret,
	}
}

func (es *ExportService) RegisterRoutes(e *echo.Echo) {
	e.GET("/export/shortcuts.html", es.handleShortcutsExport)
}

func (es *ExportService) handleShortcutsExport(c echo.Context) error {
	ctx := context.Background()

	// Get access token from query parameter
	accessToken := c.QueryParam("token")
	if accessToken == "" {
		return c.String(http.StatusUnauthorized, "Access token required. Use: /export/shortcuts.html?token=YOUR_ACCESS_TOKEN")
	}

	// Authenticate the user
	userID, err := es.authenticateUser(ctx, accessToken)
	if err != nil {
		return c.String(http.StatusUnauthorized, fmt.Sprintf("Authentication failed: %s", err.Error()))
	}

	// Get user's shortcuts only
	shortcuts, err := es.Store.ListShortcuts(ctx, &store.FindShortcut{
		CreatorID: &userID,
	})
	if err != nil {
		return errors.Wrap(err, "failed to get shortcuts")
	}

	// Sort shortcuts by creation time (most recent first)
	sort.Slice(shortcuts, func(i, j int) bool {
		return shortcuts[i].CreatedTs > shortcuts[j].CreatedTs
	})

	// Get user info for the export
	user, err := es.Store.GetUser(ctx, &store.FindUser{
		ID: &userID,
	})
	if err != nil {
		return errors.Wrap(err, "failed to get user")
	}

	htmlContent, err := es.generateBookmarkHTML(shortcuts, user)
	if err != nil {
		return errors.Wrap(err, "failed to generate bookmark HTML")
	}

	// Set headers for file download
	timestamp := time.Now().Format("1_2_06") // M_D_YY format
	filename := fmt.Sprintf("monotreme_bookmarks_%s.html", timestamp)
	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.String(http.StatusOK, htmlContent)
}

func (es *ExportService) authenticateUser(ctx context.Context, accessToken string) (int32, error) {
	if accessToken == "" {
		return 0, errors.New("access token not found")
	}

	// Parse and validate JWT token
	claims := &ClaimsMessage{}
	_, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, errors.Errorf("unexpected access token signing method=%v, expect %v", t.Header["alg"], jwt.SigningMethodHS256)
		}
		if kid, ok := t.Header["kid"].(string); ok {
			if kid == KeyID {
				return []byte(es.Secret), nil
			}
		}
		return nil, errors.Errorf("unexpected access token kid=%v", t.Header["kid"])
	})
	if err != nil {
		return 0, errors.Wrap(err, "invalid or expired access token")
	}

	// Validate audience
	audienceValid := false
	for _, audience := range claims.Audience {
		if audience == AccessTokenAudienceName {
			audienceValid = true
			break
		}
	}
	if !audienceValid {
		return 0, errors.Errorf("invalid access token audience")
	}

	// Get user ID from claims
	userID, err := util.ConvertStringToInt32(claims.Subject)
	if err != nil {
		return 0, errors.Wrapf(err, "malformed user ID in access token: %s", claims.Subject)
	}

	// Verify user exists and is active
	user, err := es.Store.GetUser(ctx, &store.FindUser{
		ID: &userID,
	})
	if err != nil {
		return 0, errors.Wrapf(err, "failed to find user with ID: %d", userID)
	}
	if user == nil {
		return 0, errors.Errorf("user not found with ID: %d", userID)
	}
	if user.RowStatus == storepb.RowStatus_ARCHIVED {
		return 0, errors.Errorf("user account has been deactivated")
	}

	// Verify access token exists in user's token list
	accessTokens, err := es.Store.GetUserAccessTokens(ctx, userID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get user access tokens")
	}

	tokenValid := false
	for _, userAccessToken := range accessTokens {
		if userAccessToken.AccessToken == accessToken {
			tokenValid = true
			break
		}
	}
	if !tokenValid {
		return 0, errors.New("access token not found in user's token list")
	}

	return userID, nil
}

func (es *ExportService) generateBookmarkHTML(shortcuts []*storepb.Shortcut, user *store.User) (string, error) {
	ctx := context.Background()

	htmlHeader := fmt.Sprintf(`<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>`)

	htmlFooter := `</DL><p>`

	// Get user's collections
	collections, err := es.Store.ListCollections(ctx, &store.FindCollection{
		CreatorID: &user.ID,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to get collections")
	}

	// Create a map of shortcut ID to shortcut for quick lookup
	shortcutMap := make(map[int32]*storepb.Shortcut)
	for _, shortcut := range shortcuts {
		shortcutMap[shortcut.Id] = shortcut
	}

	// Track which shortcuts are in collections
	shortcutsInCollections := make(map[int32]bool)

	content := ""

	// Generate collections with their shortcuts
	for _, collection := range collections {
		if len(collection.ShortcutIds) == 0 {
			continue // Skip empty collections
		}

		// Collection header
		content += fmt.Sprintf(`    <DT><H3 ADD_DATE="%d" LAST_MODIFIED="%d">%s</H3>
    <DL><p>
`, collection.CreatedTs, collection.UpdatedTs, escapeHTML(collection.Title))

		// Add shortcuts in this collection
		for _, shortcutID := range collection.ShortcutIds {
			if shortcut, exists := shortcutMap[shortcutID]; exists {
				shortcutsInCollections[shortcutID] = true
				content += fmt.Sprintf(`        <DT><A HREF="%s" ADD_DATE="%d">%s</A>
`, escapeHTML(shortcut.Link), shortcut.CreatedTs, escapeHTML(shortcut.Title))
			}
		}

		// Collection footer
		content += `    </DL><p>
`
	}

	// Add uncollected shortcuts in a default folder
	uncollectedShortcuts := []*storepb.Shortcut{}
	for _, shortcut := range shortcuts {
		if !shortcutsInCollections[shortcut.Id] {
			uncollectedShortcuts = append(uncollectedShortcuts, shortcut)
		}
	}

	if len(uncollectedShortcuts) > 0 {
		content += fmt.Sprintf(`    <DT><H3 ADD_DATE="%d" LAST_MODIFIED="%d">Uncollected Shortcuts</H3>
    <DL><p>
`, time.Now().Unix(), time.Now().Unix())

		for _, shortcut := range uncollectedShortcuts {
			content += fmt.Sprintf(`        <DT><A HREF="%s" ADD_DATE="%d">%s</A>
`, escapeHTML(shortcut.Link), shortcut.CreatedTs, escapeHTML(shortcut.Title))
		}

		content += `    </DL><p>
`
	}

	return htmlHeader + content + htmlFooter, nil
}

func escapeHTML(s string) string {
	// Basic HTML escaping for bookmark files
	s = fmt.Sprintf("%s", s)
	// More comprehensive escaping would be needed for production
	return s
}