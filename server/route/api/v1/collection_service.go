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

	var createdCollections []*v1pb.Collection
	totalShortcuts := int32(0)

	// Create collections and shortcuts
	for _, collection := range bookmarkData.Collections {
		// Create shortcuts for this collection first
		var shortcutIDs []int32
		for _, bookmark := range collection.Bookmarks {
			shortcut := &storepb.Shortcut{
				CreatorId:   user.ID,
				Name:        generateShortcutName(bookmark.Title, bookmark.URL),
				Link:        bookmark.URL,
				Title:       bookmark.Title,
				Description: "",
				Tags:        []string{},
				Visibility:  storepb.Visibility_WORKSPACE,
				OgMetadata:  &storepb.OpenGraphMetadata{},
				Uuid:        uuid.New().String(),
			}
			createdShortcut, err := s.Store.CreateShortcut(ctx, shortcut)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to create shortcut: %v", err)
			}
			shortcutIDs = append(shortcutIDs, createdShortcut.Id)
			totalShortcuts++
		}

		// Create the collection
		collectionCreate := &storepb.Collection{
			CreatorId:   user.ID,
			Name:        generateCollectionName(collection.Title),
			Title:       collection.Title,
			Description: fmt.Sprintf("Imported bookmark collection: %s", collection.Title),
			ShortcutIds: shortcutIDs,
			Visibility:  storepb.Visibility_WORKSPACE,
		}
		createdCollection, err := s.Store.CreateCollection(ctx, collectionCreate)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create collection: %v", err)
		}
		createdCollections = append(createdCollections, convertCollectionFromStore(createdCollection))
	}

	return &v1pb.ImportBookmarksResponse{
		Collections:      createdCollections,
		TotalShortcuts:   totalShortcuts,
		TotalCollections: int32(len(createdCollections)),
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
}

func parseBookmarksHTML(htmlContent string) (*BookmarkData, error) {
	// Regular expressions to parse the HTML bookmark structure
	h3Regex := regexp.MustCompile(`<H3[^>]*>([^<]+)</H3>`)
	linkRegex := regexp.MustCompile(`<A[^>]+HREF="([^"]+)"[^>]*>([^<]+)</A>`)

	var collections []BookmarkCollection
	lines := strings.Split(htmlContent, "\n")

	var currentCollection *BookmarkCollection
	var inCollection bool

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for H3 (collection title)
		if h3Match := h3Regex.FindStringSubmatch(line); h3Match != nil {
			// Save previous collection if exists
			if currentCollection != nil {
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

		// Check for closing DL tag to end collection
		if strings.Contains(line, "</DL>") && inCollection {
			inCollection = false
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
	if currentCollection != nil {
		collections = append(collections, *currentCollection)
	}

	return &BookmarkData{Collections: collections}, nil
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
