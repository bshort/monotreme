package v1

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/google/uuid"

	v1pb "github.com/bshort/monotreme/proto/gen/api/v1"
	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/server/service/license"
	"github.com/bshort/monotreme/store"
)

func (s *APIV1Service) ListCollections(ctx context.Context, _ *v1pb.ListCollectionsRequest) (*v1pb.ListCollectionsResponse, error) {
	collections, err := s.Store.ListCollections(ctx, &store.FindCollection{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get collection list, err: %v", err)
	}

	convertedCollections := []*v1pb.Collection{}
	for _, collection := range collections {
		convertedCollections = append(convertedCollections, convertCollectionFromStore(collection))
	}

	response := &v1pb.ListCollectionsResponse{
		Collections: convertedCollections,
	}
	return response, nil
}

func (s *APIV1Service) GetCollection(ctx context.Context, request *v1pb.GetCollectionRequest) (*v1pb.Collection, error) {
	collection, err := s.Store.GetCollection(ctx, &store.FindCollection{
		ID: &request.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get collection by name: %v", err)
	}
	if collection == nil {
		return nil, status.Errorf(codes.NotFound, "collection not found")
	}

	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil && collection.Visibility != storepb.Visibility_PUBLIC {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}
	return convertCollectionFromStore(collection), nil
}

func (s *APIV1Service) GetCollectionByName(ctx context.Context, request *v1pb.GetCollectionByNameRequest) (*v1pb.Collection, error) {
	collection, err := s.Store.GetCollection(ctx, &store.FindCollection{
		Name: &request.Name,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get collection by name: %v", err)
	}
	if collection == nil {
		return nil, status.Errorf(codes.NotFound, "collection not found")
	}

	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil && collection.Visibility != storepb.Visibility_PUBLIC {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}
	return convertCollectionFromStore(collection), nil
}

func (s *APIV1Service) CreateCollection(ctx context.Context, request *v1pb.CreateCollectionRequest) (*v1pb.Collection, error) {
	if request.Collection.Name == "" || request.Collection.Title == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name and title are required")
	}

	if !s.LicenseService.IsFeatureEnabled(license.FeatureTypeUnlimitedCollections) {
		collections, err := s.Store.ListCollections(ctx, &store.FindCollection{})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get collection list, err: %v", err)
		}
		collectionsLimit := int(s.LicenseService.GetSubscription().CollectionsLimit)
		if len(collections) >= collectionsLimit {
			return nil, status.Errorf(codes.PermissionDenied, "Maximum number of collections %d reached", collectionsLimit)
		}
	}

	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	collectionCreate := &storepb.Collection{
		CreatorId:   user.ID,
		Name:        request.Collection.Name,
		Title:       request.Collection.Title,
		Description: request.Collection.Description,
		ShortcutIds: request.Collection.ShortcutIds,
		Visibility:  convertVisibilityToStorepb(request.Collection.Visibility),
	}
	collection, err := s.Store.CreateCollection(ctx, collectionCreate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create collection, err: %v", err)
	}

	return convertCollectionFromStore(collection), nil
}

func (s *APIV1Service) UpdateCollection(ctx context.Context, request *v1pb.UpdateCollectionRequest) (*v1pb.Collection, error) {
	if request.UpdateMask == nil || len(request.UpdateMask.Paths) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "updateMask is required")
	}

	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	collection, err := s.Store.GetCollection(ctx, &store.FindCollection{
		ID: &request.Collection.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get collection by name: %v", err)
	}
	if collection == nil {
		return nil, status.Errorf(codes.NotFound, "collection not found")
	}
	if collection.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	update := &store.UpdateCollection{
		ID: collection.Id,
	}
	for _, path := range request.UpdateMask.Paths {
		switch path {
		case "name":
			update.Name = &request.Collection.Name
		case "title":
			update.Title = &request.Collection.Title
		case "description":
			update.Description = &request.Collection.Description
		case "shortcut_ids":
			update.ShortcutIDs = request.Collection.ShortcutIds
		case "visibility":
			visibility := convertVisibilityToStorepb(request.Collection.Visibility)
			update.Visibility = &visibility
		}
	}
	collection, err = s.Store.UpdateCollection(ctx, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update collection, err: %v", err)
	}

	return convertCollectionFromStore(collection), nil
}

func (s *APIV1Service) DeleteCollection(ctx context.Context, request *v1pb.DeleteCollectionRequest) (*emptypb.Empty, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	collection, err := s.Store.GetCollection(ctx, &store.FindCollection{
		ID: &request.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get collection by name: %v", err)
	}
	if collection == nil {
		return nil, status.Errorf(codes.NotFound, "collection not found")
	}
	if collection.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	err = s.Store.DeleteCollection(ctx, &store.DeleteCollection{
		ID: collection.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete collection, err: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *APIV1Service) ImportBookmarks(ctx context.Context, request *v1pb.ImportBookmarksRequest) (*v1pb.ImportBookmarksResponse, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil || user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Only admin users can import bookmarks")
	}

	// Parse the HTML content to extract bookmark collections and shortcuts
	bookmarkData, err := parseBookmarksHTML(request.HtmlContent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse bookmark HTML: %v", err)
	}

	// Log the detected format for debugging
	var formatStr string
	switch bookmarkData.Format {
	case FormatChrome:
		formatStr = "Chrome/Chromium"
	case FormatFirefox:
		formatStr = "Firefox"
	default:
		formatStr = "Unknown"
	}
	fmt.Printf("Detected bookmark format: %s\n", formatStr)

	var createdCollections []*v1pb.Collection
	totalShortcuts := int32(0)
	var shortcutsCreated, shortcutsUpdated, collectionsCreated, collectionsUpdated int32

	// Create collections and shortcuts
	for _, collection := range bookmarkData.Collections {
		// Create or update shortcuts for this collection first
		var shortcutIDs []int32
		for _, bookmark := range collection.Bookmarks {
			shortcutName := generateShortcutName(bookmark.Title, bookmark.URL)

			// Check if a shortcut with this name already exists
			existingShortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
				Name: &shortcutName,
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to check existing shortcut: %v", err)
			}

			var resultShortcut *storepb.Shortcut
			if existingShortcut != nil {
				// Update existing shortcut
				updatedShortcut, err := s.Store.UpdateShortcut(ctx, &store.UpdateShortcut{
					ID:          existingShortcut.Id,
					Link:        &bookmark.URL,
					Title:       &bookmark.Title,
					Description: stringPtr("Updated from bookmark import"),
				})
				if err != nil {
					return nil, status.Errorf(codes.Internal, "failed to update existing shortcut: %v", err)
				}
				resultShortcut = updatedShortcut
				shortcutsUpdated++
			} else {
				// Create new shortcut
				newShortcut := &storepb.Shortcut{
					CreatorId:   user.ID,
					Name:        shortcutName,
					Link:        bookmark.URL,
					Title:       bookmark.Title,
					Description: "",
					Tags:        []string{},
					Visibility:  storepb.Visibility_WORKSPACE,
					OgMetadata:  &storepb.OpenGraphMetadata{},
					Uuid:        uuid.New().String(),
				}
				createdShortcut, err := s.Store.CreateShortcut(ctx, newShortcut)
				if err != nil {
					return nil, status.Errorf(codes.Internal, "failed to create shortcut: %v", err)
				}
				resultShortcut = createdShortcut
				shortcutsCreated++
			}
			totalShortcuts++
			shortcutIDs = append(shortcutIDs, resultShortcut.Id)
		}

		// Create or update the collection
		collectionName := generateCollectionName(collection.Title)

		// Check if a collection with this name already exists
		existingCollection, err := s.Store.GetCollection(ctx, &store.FindCollection{
			Name: &collectionName,
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to check existing collection: %v", err)
		}

		var resultCollection *storepb.Collection
		if existingCollection != nil {
			// Update existing collection by adding new shortcuts to it
			combinedShortcutIDs := combineShortcutIDs(existingCollection.ShortcutIds, shortcutIDs)
			updatedCollection, err := s.Store.UpdateCollection(ctx, &store.UpdateCollection{
				ID:          existingCollection.Id,
				Title:       &collection.Title,
				Description: stringPtr(fmt.Sprintf("Updated bookmark collection: %s", collection.Title)),
				ShortcutIDs: combinedShortcutIDs,
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update existing collection: %v", err)
			}
			resultCollection = updatedCollection
			collectionsUpdated++
		} else {
			// Create new collection
			collectionCreate := &storepb.Collection{
				CreatorId:   user.ID,
				Name:        collectionName,
				Title:       collection.Title,
				Description: fmt.Sprintf("Imported bookmark collection: %s", collection.Title),
				ShortcutIds: shortcutIDs,
				Visibility:  storepb.Visibility_WORKSPACE,
			}
			createdCollection, err := s.Store.CreateCollection(ctx, collectionCreate)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to create collection: %v", err)
			}
			resultCollection = createdCollection
			collectionsCreated++
		}
		createdCollections = append(createdCollections, convertCollectionFromStore(resultCollection))
	}

	return &v1pb.ImportBookmarksResponse{
		Collections:        createdCollections,
		TotalShortcuts:     totalShortcuts,
		TotalCollections:   int32(len(createdCollections)),
		ShortcutsCreated:   shortcutsCreated,
		ShortcutsUpdated:   shortcutsUpdated,
		CollectionsCreated: collectionsCreated,
		CollectionsUpdated: collectionsUpdated,
	}, nil
}

type BookmarkCollection struct {
	Title     string
	Bookmarks []Bookmark
}

type Bookmark struct {
	Title string
	URL   string
}

type BookmarkData struct {
	Collections []BookmarkCollection
	Format      BookmarkFormat
}

type BookmarkFormat int

const (
	FormatUnknown BookmarkFormat = iota
	FormatChrome
	FormatFirefox
)

func parseBookmarksHTML(htmlContent string) (*BookmarkData, error) {
	// Detect the bookmark format
	format := detectBookmarkFormat(htmlContent)

	var collections []BookmarkCollection
	var err error

	switch format {
	case FormatFirefox:
		collections, err = parseFirefoxBookmarks(htmlContent)
	case FormatChrome:
		collections, err = parseChromeBookmarks(htmlContent)
	default:
		// Try Chrome format as fallback
		collections, err = parseChromeBookmarks(htmlContent)
		if err != nil {
			// If Chrome parsing fails, try Firefox format
			collections, err = parseFirefoxBookmarks(htmlContent)
			format = FormatFirefox
		} else {
			format = FormatChrome
		}
	}

	if err != nil {
		return nil, err
	}

	return &BookmarkData{
		Collections: collections,
		Format:      format,
	}, nil
}

func detectBookmarkFormat(htmlContent string) BookmarkFormat {
	// Firefox bookmarks typically have specific patterns
	firefoxPatterns := []string{
		"LAST_MODIFIED=",           // Firefox timestamps
		"PERSONAL_TOOLBAR_FOLDER=", // Firefox toolbar attribute
		"Mozilla Firefox",          // Firefox user agent in comments
		"Bookmarks Menu",           // Firefox default folder name
		"SHORTCUTURL=",            // Firefox shortcut URLs
	}

	// Chrome/Chromium bookmarks patterns
	chromePatterns := []string{
		"Bookmarks bar",        // Chrome toolbar name
		"Other bookmarks",      // Chrome default folder
		"Chromium",            // Chromium user agent
		"Chrome",              // Chrome user agent
	}

	firefoxScore := 0
	chromeScore := 0

	htmlLower := strings.ToLower(htmlContent)

	// Count Firefox-specific patterns
	for _, pattern := range firefoxPatterns {
		if strings.Contains(htmlLower, strings.ToLower(pattern)) {
			firefoxScore++
		}
	}

	// Count Chrome-specific patterns
	for _, pattern := range chromePatterns {
		if strings.Contains(htmlLower, strings.ToLower(pattern)) {
			chromeScore++
		}
	}

	// Also check for specific attributes
	if strings.Contains(htmlContent, `PERSONAL_TOOLBAR_FOLDER="true"`) {
		firefoxScore += 2
	}
	if strings.Contains(htmlContent, `ADD_DATE=`) && strings.Contains(htmlContent, `LAST_MODIFIED=`) {
		firefoxScore++
	}

	if firefoxScore > chromeScore {
		return FormatFirefox
	} else if chromeScore > 0 {
		return FormatChrome
	}

	return FormatUnknown
}

func parseChromeBookmarks(htmlContent string) ([]BookmarkCollection, error) {
	h3Regex := regexp.MustCompile(`<H3[^>]*>([^<]+)</H3>`)
	linkRegex := regexp.MustCompile(`<A[^>]+HREF="([^"]+)"[^>]*>([^<]+)</A>`)

	var collections []BookmarkCollection
	lines := strings.Split(htmlContent, "\n")

	var currentCollection *BookmarkCollection
	var inCollection bool
	var dlDepth int

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Track DL nesting depth
		if strings.Contains(line, "<DL>") {
			dlDepth++
		}
		if strings.Contains(line, "</DL>") {
			dlDepth--
			if dlDepth <= 1 {
				inCollection = false
			}
		}

		// Check for H3 (collection title)
		if h3Match := h3Regex.FindStringSubmatch(line); h3Match != nil {
			// Save previous collection if exists
			if currentCollection != nil && len(currentCollection.Bookmarks) > 0 {
				collections = append(collections, *currentCollection)
			}
			// Start new collection
			currentCollection = &BookmarkCollection{
				Title:     strings.TrimSpace(h3Match[1]),
				Bookmarks: []Bookmark{},
			}
			inCollection = true
			continue
		}

		// Check for bookmark links within a collection
		if inCollection && currentCollection != nil {
			if linkMatch := linkRegex.FindStringSubmatch(line); linkMatch != nil {
				bookmark := Bookmark{
					URL:   strings.TrimSpace(linkMatch[1]),
					Title: strings.TrimSpace(linkMatch[2]),
				}
				currentCollection.Bookmarks = append(currentCollection.Bookmarks, bookmark)
			}
		}
	}

	// Don't forget the last collection
	if currentCollection != nil && len(currentCollection.Bookmarks) > 0 {
		collections = append(collections, *currentCollection)
	}

	return collections, nil
}

func parseFirefoxBookmarks(htmlContent string) ([]BookmarkCollection, error) {
	// Firefox has a slightly different structure and may use nested folders differently
	h3Regex := regexp.MustCompile(`<H3[^>]*>([^<]+)</H3>`)
	linkRegex := regexp.MustCompile(`<A[^>]+HREF="([^"]+)"[^>]*>([^<]+)</A>`)

	var collections []BookmarkCollection
	lines := strings.Split(htmlContent, "\n")

	var currentCollection *BookmarkCollection
	var inCollection bool
	var dlStack []string // Track nested DL elements

	for _, line := range lines {
		originalLine := line
		line = strings.TrimSpace(line)

		// Track DL stack for Firefox's more complex nesting
		if strings.Contains(line, "<DL>") {
			dlStack = append(dlStack, "DL")
		}
		if strings.Contains(line, "</DL>") && len(dlStack) > 0 {
			dlStack = dlStack[:len(dlStack)-1]
			// End collection if we're back to root level
			if len(dlStack) <= 1 {
				inCollection = false
			}
		}

		// Check for H3 (collection title) - Firefox may have more attributes
		if h3Match := h3Regex.FindStringSubmatch(line); h3Match != nil {
			title := strings.TrimSpace(h3Match[1])

			// Skip certain Firefox system folders
			if isFirefoxSystemFolder(title) {
				continue
			}

			// Save previous collection if exists
			if currentCollection != nil && len(currentCollection.Bookmarks) > 0 {
				collections = append(collections, *currentCollection)
			}

			// Start new collection
			currentCollection = &BookmarkCollection{
				Title:     title,
				Bookmarks: []Bookmark{},
			}
			inCollection = true
			continue
		}

		// Check for bookmark links within a collection
		if inCollection && currentCollection != nil && len(dlStack) >= 2 {
			if linkMatch := linkRegex.FindStringSubmatch(originalLine); linkMatch != nil {
				url := strings.TrimSpace(linkMatch[1])
				title := strings.TrimSpace(linkMatch[2])

				// Skip empty or invalid bookmarks
				if url != "" && title != "" && !strings.HasPrefix(url, "place:") {
					bookmark := Bookmark{
						URL:   url,
						Title: title,
					}
					currentCollection.Bookmarks = append(currentCollection.Bookmarks, bookmark)
				}
			}
		}
	}

	// Don't forget the last collection
	if currentCollection != nil && len(currentCollection.Bookmarks) > 0 {
		collections = append(collections, *currentCollection)
	}

	return collections, nil
}

func isFirefoxSystemFolder(title string) bool {
	systemFolders := []string{
		"Recently Bookmarked",
		"Recent Tags",
		"Most Visited",
		"Getting Started",
		"Mozilla Firefox",
	}

	for _, folder := range systemFolders {
		if strings.EqualFold(title, folder) {
			return true
		}
	}
	return false
}

func generateCollectionName(title string) string {
	// Convert title to a valid collection name (lowercase, replace spaces with hyphens, etc.)
	name := strings.ToLower(title)
	name = strings.ReplaceAll(name, " ", "-")
	name = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(name, "")
	if name == "" {
		name = "imported-collection"
	}
	return name
}

func generateShortcutName(title, url string) string {
	// Generate a short name for the shortcut
	if title != "" {
		name := strings.ToLower(title)
		name = strings.ReplaceAll(name, " ", "-")
		name = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(name, "")
		if len(name) > 30 {
			name = name[:30]
		}
		if name != "" {
			return name
		}
	}

	// Fallback: use domain from URL
	if strings.Contains(url, "://") {
		parts := strings.Split(url, "://")
		if len(parts) > 1 {
			domain := strings.Split(parts[1], "/")[0]
			domain = strings.ReplaceAll(domain, "www.", "")
			if domain != "" {
				return domain
			}
		}
	}

	return "bookmark"
}

// Helper function to convert string to pointer
func stringPtr(s string) *string {
	return &s
}

// Helper function to combine shortcut IDs without duplicates
func combineShortcutIDs(existing []int32, new []int32) []int32 {
	// Create a map to track existing IDs
	idMap := make(map[int32]bool)
	var result []int32

	// Add existing IDs
	for _, id := range existing {
		if !idMap[id] {
			idMap[id] = true
			result = append(result, id)
		}
	}

	// Add new IDs that don't already exist
	for _, id := range new {
		if !idMap[id] {
			idMap[id] = true
			result = append(result, id)
		}
	}

	return result
}

func convertCollectionFromStore(collection *storepb.Collection) *v1pb.Collection {
	return &v1pb.Collection{
		Id:          collection.Id,
		CreatorId:   collection.CreatorId,
		CreatedTime: timestamppb.New(time.Unix(collection.CreatedTs, 0)),
		UpdatedTime: timestamppb.New(time.Unix(collection.UpdatedTs, 0)),
		Name:        collection.Name,
		Title:       collection.Title,
		Description: collection.Description,
		ShortcutIds: collection.ShortcutIds,
		Visibility:  convertVisibilityFromStorepb(collection.Visibility),
	}
}
