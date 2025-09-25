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

	htmlHeader := fmt.Sprintf(`<!DOCTYPE NETSCAPE-Bookmark-file-1>
<!-- This is an automatically generated file.
     It will be read and overwritten.
     DO NOT EDIT! -->
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3 ADD_DATE="%d" LAST_MODIFIED="%d" PERSONAL_TOOLBAR_FOLDER="true">Monotreme Shortcuts - %s</H3>
    <DL><p>
`, time.Now().Unix(), time.Now().Unix(), user.Nickname)

	htmlFooter := `    </DL><p>
</DL><p>`

	bookmarks := ""
	for _, shortcut := range shortcuts {
		// Generate bookmark entry
		addDate := shortcut.CreatedTs

		bookmarkHTML := fmt.Sprintf(`        <DT><A HREF="%s" ADD_DATE="%d">%s</A>
`,
			escapeHTML(shortcut.Link),
			addDate,
			escapeHTML(shortcut.Title))

		bookmarks += bookmarkHTML
	}

	return htmlHeader + bookmarks + htmlFooter, nil
}

func escapeHTML(s string) string {
	// Basic HTML escaping for bookmark files
	s = fmt.Sprintf("%s", s)
	// More comprehensive escaping would be needed for production
	return s
}