package swagger

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/bshort/monotreme/internal/util"
	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

const (
	// Constants for JWT token validation
	Issuer                     = "monotreme"
	KeyID                      = "v1"
	AccessTokenAudienceName    = "user.access-token"
	AccessTokenCookieName      = "access-token"
)

// ClaimsMessage represents JWT claims
type ClaimsMessage struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type SwaggerService struct {
	swaggerSpec string
	store       *store.Store
	secret      string
}

// NewSwaggerService creates a new SwaggerService.
func NewSwaggerService(swaggerSpec string, store *store.Store, secret string) *SwaggerService {
	return &SwaggerService{
		swaggerSpec: swaggerSpec,
		store:       store,
		secret:      secret,
	}
}

// RegisterRoutes registers Swagger UI and spec endpoints
func (s *SwaggerService) RegisterRoutes(e *echo.Echo) {
	// Apply authentication middleware to all swagger routes
	swaggerGroup := e.Group("")
	swaggerGroup.Use(s.authMiddleware)

	// Serve the swagger spec
	swaggerGroup.GET("/api/v1/swagger.yaml", s.serveSwaggerSpec)

	// Serve Swagger UI
	swaggerGroup.GET("/api-docs", s.serveSwaggerUI)

	// Redirect /api-docs/ to /api-docs for convenience
	swaggerGroup.GET("/api-docs/", s.redirectToSwaggerUI)
}

// serveSwaggerSpec serves the raw swagger specification
func (s *SwaggerService) serveSwaggerSpec(c echo.Context) error {
	return c.String(http.StatusOK, s.swaggerSpec)
}

// serveSwaggerUI serves the Swagger UI interface
func (s *SwaggerService) serveSwaggerUI(c echo.Context) error {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Monotreme API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui.css" />
    <link rel="icon" href="/monotreme.png" type="image/png">
   <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/api/v1/swagger.yaml',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}

// redirectToSwaggerUI redirects /api-docs/ to /api-docs
func (s *SwaggerService) redirectToSwaggerUI(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/api-docs")
}

// authMiddleware validates user authentication before allowing access to swagger
func (s *SwaggerService) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// Try to get token from cookie first
		accessToken := ""
		if cookie, err := c.Cookie(AccessTokenCookieName); err == nil {
			accessToken = cookie.Value
		}

		// If no cookie, try Authorization header
		if accessToken == "" {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					accessToken = parts[1]
				}
			}
		}

		// If no token found, return unauthorized with login prompt
		if accessToken == "" {
			return s.unauthorizedResponse(c)
		}

		// Authenticate the user
		_, err := s.authenticateUser(ctx, accessToken)
		if err != nil {
			return s.unauthorizedResponse(c)
		}

		// User is authenticated, continue to the next handler
		return next(c)
	}
}

// authenticateUser validates the access token and returns the user ID
func (s *SwaggerService) authenticateUser(ctx context.Context, accessToken string) (int32, error) {
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
				return []byte(s.secret), nil
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
	user, err := s.store.GetUser(ctx, &store.FindUser{
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
	accessTokens, err := s.store.GetUserAccessTokens(ctx, userID)
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

// unauthorizedResponse returns a user-friendly unauthorized response
func (s *SwaggerService) unauthorizedResponse(c echo.Context) error {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation - Login Required</title>
    <link rel="icon" href="/monotreme.png" type="image/png">
    <style>
        body {
            font-family: system-ui, -apple-system, sans-serif;
            max-width: 600px;
            margin: 0 auto;
            padding: 2rem;
            background-color: #f8fafc;
            color: #1e293b;
        }
        .container {
            background: white;
            border-radius: 8px;
            padding: 2rem;
            border: 1px solid #e2e8f0;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            text-align: center;
        }
        h1 {
            color: #1e293b;
            margin-bottom: 1rem;
        }
        p {
            color: #64748b;
            line-height: 1.6;
            margin-bottom: 1.5rem;
        }
        .login-button {
            background: #3b82f6;
            color: white;
            padding: 0.75rem 1.5rem;
            border: none;
            border-radius: 6px;
            font-size: 1rem;
            text-decoration: none;
            display: inline-block;
            transition: background-color 0.2s;
        }
        .login-button:hover {
            background: #2563eb;
            text-decoration: none;
            color: white;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔒 Authentication Required</h1>
        <p>You need to be logged in to access the API documentation.</p>
        <p>Please log in to your Monotreme account first, then return to this page.</p>
        <a href="/" class="login-button">Go to Login</a>
    </div>
</body>
</html>`
	return c.HTML(http.StatusUnauthorized, html)
}