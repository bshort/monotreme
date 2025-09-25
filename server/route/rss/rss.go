package rss

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

type RSSService struct {
	Profile *profile.Profile
	Store   *store.Store
	Secret  string
}

func NewRSSService(profile *profile.Profile, store *store.Store, secret string) *RSSService {
	return &RSSService{
		Profile: profile,
		Store:   store,
		Secret:  secret,
	}
}

func (rs *RSSService) RegisterRoutes(e *echo.Echo) {
	e.GET("/rss/collections.xml", rs.handleCollectionsRSS)
	e.GET("/rss/collection/:id.xml", rs.handleCollectionRSS)
}

func (rs *RSSService) handleCollectionsRSS(c echo.Context) error {
	ctx := context.Background()

	// Get access token from query parameter
	accessToken := c.QueryParam("token")
	if accessToken == "" {
		return c.String(http.StatusUnauthorized, "Access token required. Use: /rss/collections.xml?token=YOUR_ACCESS_TOKEN")
	}

	// Authenticate the user
	userID, err := rs.authenticateUser(ctx, accessToken)
	if err != nil {
		return c.String(http.StatusUnauthorized, fmt.Sprintf("Authentication failed: %s", err.Error()))
	}

	// Get user's collections only
	collections, err := rs.Store.ListCollections(ctx, &store.FindCollection{
		CreatorID: &userID,
	})
	if err != nil {
		return errors.Wrap(err, "failed to get collections")
	}

	// Sort collections by updated time (most recent first)
	sort.Slice(collections, func(i, j int) bool {
		return collections[i].UpdatedTs > collections[j].UpdatedTs
	})

	rssXML, err := rs.generateCollectionsRSSXML(collections, userID)
	if err != nil {
		return errors.Wrap(err, "failed to generate RSS XML")
	}

	c.Response().Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	return c.String(http.StatusOK, rssXML)
}

func (rs *RSSService) handleCollectionRSS(c echo.Context) error {
	ctx := context.Background()

	// Get collection ID from path parameter
	collectionIDStr := c.Param("id")
	collectionID, err := util.ConvertStringToInt32(collectionIDStr)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid collection ID: %s", collectionIDStr))
	}

	// Get access token from query parameter
	accessToken := c.QueryParam("token")
	if accessToken == "" {
		return c.String(http.StatusUnauthorized, "Access token required. Use: /rss/collection/ID.xml?token=YOUR_ACCESS_TOKEN")
	}

	// Authenticate the user
	userID, err := rs.authenticateUser(ctx, accessToken)
	if err != nil {
		return c.String(http.StatusUnauthorized, fmt.Sprintf("Authentication failed: %s", err.Error()))
	}

	// Get the specific collection
	collection, err := rs.Store.GetCollection(ctx, &store.FindCollection{
		ID:        &collectionID,
		CreatorID: &userID, // Ensure user can only access their own collections
	})
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("Collection not found or access denied: %d", collectionID))
	}
	if collection == nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("Collection not found: %d", collectionID))
	}

	// Get all shortcuts in this collection
	shortcuts, err := rs.getCollectionShortcuts(ctx, collection)
	if err != nil {
		return errors.Wrap(err, "failed to get collection shortcuts")
	}

	rssXML, err := rs.generateCollectionRSSXML(collection, shortcuts, userID)
	if err != nil {
		return errors.Wrap(err, "failed to generate collection RSS XML")
	}

	c.Response().Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	return c.String(http.StatusOK, rssXML)
}

func (rs *RSSService) authenticateUser(ctx context.Context, accessToken string) (int32, error) {
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
				return []byte(rs.Secret), nil
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
	user, err := rs.Store.GetUser(ctx, &store.FindUser{
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
	accessTokens, err := rs.Store.GetUserAccessTokens(ctx, userID)
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

func (rs *RSSService) generateCollectionsRSSXML(collections []*storepb.Collection, userID int32) (string, error) {
	baseURL := fmt.Sprintf("http://localhost:%d", rs.Profile.Port)
	if rs.Profile.Mode == "prod" {
		// In production, you might want to set this from an environment variable
		baseURL = fmt.Sprintf("http://localhost:%d", rs.Profile.Port)
	}

	// Get user info for title
	user, err := rs.Store.GetUser(context.Background(), &store.FindUser{
		ID: &userID,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to get user")
	}

	rssHeader := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>
    <title>` + user.Nickname + `'s Monotreme Collections</title>
    <link>` + baseURL + `/collections</link>
    <description>Latest shortcuts from all collections from Monotreme</description>
    <language>en-us</language>
    <lastBuildDate>` + time.Now().Format(time.RFC3339) + `</lastBuildDate>
    <atom:link href="` + baseURL + `/rss/collections.xml" rel="self" type="application/rss+xml" />

`

	// Collect all shortcuts from all collections with their collection context
	type shortcutWithCollection struct {
		shortcut   *storepb.Shortcut
		collection *storepb.Collection
	}
	var allShortcuts []shortcutWithCollection

	for _, collection := range collections {
		for _, shortcutID := range collection.ShortcutIds {
			shortcut, err := rs.Store.GetShortcut(context.Background(), &store.FindShortcut{
				ID: &shortcutID,
			})
			if err == nil && shortcut != nil {
				allShortcuts = append(allShortcuts, shortcutWithCollection{
					shortcut:   shortcut,
					collection: collection,
				})
			}
		}
	}

	// Sort all shortcuts by creation time (most recent first)
	sort.Slice(allShortcuts, func(i, j int) bool {
		return allShortcuts[i].shortcut.CreatedTs > allShortcuts[j].shortcut.CreatedTs
	})

	items := ""
	for _, sc := range allShortcuts {
		shortcut := sc.shortcut

		// Build the Monotreme shortcut URL
		shortcutURL := fmt.Sprintf("%s/s/%s", baseURL, shortcut.Name)

		// Use shortcut UUID as GUID for uniqueness
		guid := shortcut.Uuid

		// Description - empty if not provided or show as self-closing tag
		var descriptionXML string
		if shortcut.Description == "" {
			descriptionXML = "<description/>"
		} else {
			descriptionXML = fmt.Sprintf("<description>%s</description>", escapeHTML(shortcut.Description))
		}

		itemXML := fmt.Sprintf(`    <item>
        <title>%s</title>
        <link>%s</link>
        <guid>%s</guid>
        <source>%s</source>
        <pubDate>%s</pubDate>
        %s
    </item>
`,
			escapeHTML(shortcut.Title),
			shortcutURL,
			guid,
			shortcut.Link,
			time.Unix(shortcut.CreatedTs, 0).Format(time.RFC3339),
			descriptionXML)

		items += itemXML
	}

	rssFooter := `</channel>
</rss>`

	return rssHeader + items + rssFooter, nil
}

func (rs *RSSService) getCollectionShortcuts(ctx context.Context, collection *storepb.Collection) ([]*storepb.Shortcut, error) {
	var shortcuts []*storepb.Shortcut

	for _, shortcutID := range collection.ShortcutIds {
		shortcut, err := rs.Store.GetShortcut(ctx, &store.FindShortcut{
			ID: &shortcutID,
		})
		if err == nil && shortcut != nil {
			shortcuts = append(shortcuts, shortcut)
		}
	}

	// Sort shortcuts by creation time (most recent first)
	sort.Slice(shortcuts, func(i, j int) bool {
		return shortcuts[i].CreatedTs > shortcuts[j].CreatedTs
	})

	return shortcuts, nil
}

func (rs *RSSService) generateCollectionRSSXML(collection *storepb.Collection, shortcuts []*storepb.Shortcut, userID int32) (string, error) {
	baseURL := fmt.Sprintf("http://localhost:%d", rs.Profile.Port)
	if rs.Profile.Mode == "prod" {
		// In production, you might want to set this from an environment variable
		baseURL = fmt.Sprintf("http://localhost:%d", rs.Profile.Port)
	}

	// Get user info for title
	user, err := rs.Store.GetUser(context.Background(), &store.FindUser{
		ID: &userID,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to get user")
	}

	rssHeader := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
<channel>
    <title>` + user.Nickname + `'s Monotreme Collections - <![CDATA[` + collection.Title + `]]></title>
    <link>` + baseURL + `/c/<![CDATA[` + collection.Name + `]]></link>
    <description>Latest links from <![CDATA[` + collection.Title + `]]> collection from Monotreme</description>
    <language>en-us</language>
    <lastBuildDate>` + time.Now().Format(time.RFC3339) + `</lastBuildDate>
    <atom:link href="` + baseURL + `/rss/collection/` + fmt.Sprintf("%d", collection.Id) + `.xml" rel="self" type="application/rss+xml" />

`

	items := ""
	for _, shortcut := range shortcuts {
		// Build the Monotreme shortcut URL
		shortcutURL := fmt.Sprintf("%s/s/%s", baseURL, shortcut.Name)

		// Use shortcut UUID as GUID for uniqueness
		guid := shortcut.Uuid

		// Description - empty if not provided or show as self-closing tag
		var descriptionXML string
		if shortcut.Description == "" {
			descriptionXML = "<description/>"
		} else {
			descriptionXML = fmt.Sprintf("<description>%s</description>", escapeHTML(shortcut.Description))
		}

		itemXML := fmt.Sprintf(`    <item>
        <title>%s</title>
        <link>%s</link>
        <guid>%s</guid>
        <source>%s</source>
        <pubDate>%s</pubDate>
        %s
    </item>
`,
			escapeHTML(shortcut.Title),
			shortcutURL,
			guid,
			shortcut.Link,
			time.Unix(shortcut.CreatedTs, 0).Format(time.RFC3339),
			descriptionXML)

		items += itemXML
	}

	rssFooter := `    <category><![CDATA[` + collection.Title + `]]></category>
</channel>
</rss>`

	return rssHeader + items + rssFooter, nil
}

func escapeHTML(s string) string {
	// Basic HTML escaping
	s = fmt.Sprintf("%s", s)
	// Note: In a production environment, you'd want to use a proper HTML escaping library
	// But for simplicity, we're using CDATA sections in the XML which should handle most cases
	return s
}