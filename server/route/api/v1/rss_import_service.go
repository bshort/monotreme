package v1

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	apiv1pb "github.com/bshort/monotreme/proto/gen/api/v1"
	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

func (s *APIV1Service) ListRssFeeds(ctx context.Context, _ *apiv1pb.ListRssFeedsRequest) (*apiv1pb.ListRssFeedsResponse, error) {
	userID := ctx.Value(userIDContextKey).(int32)

	rssFeeds, err := s.Store.ListRssFeeds(ctx, &store.FindRssFeed{
		CreatorID: &userID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list RSS feeds: %v", err)
	}

	response := &apiv1pb.ListRssFeedsResponse{}
	for _, rssFeed := range rssFeeds {
		response.RssFeeds = append(response.RssFeeds, convertRssFeedFromStore(rssFeed))
	}

	return response, nil
}

func (s *APIV1Service) GetRssFeed(ctx context.Context, request *apiv1pb.GetRssFeedRequest) (*apiv1pb.RssFeed, error) {
	userID := ctx.Value(userIDContextKey).(int32)

	rssFeed, err := s.Store.GetRssFeed(ctx, &store.FindRssFeed{
		ID:        &request.Id,
		CreatorID: &userID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get RSS feed: %v", err)
	}
	if rssFeed == nil {
		return nil, status.Errorf(codes.NotFound, "RSS feed not found")
	}

	return convertRssFeedFromStore(rssFeed), nil
}

func (s *APIV1Service) CreateRssFeed(ctx context.Context, request *apiv1pb.CreateRssFeedRequest) (*apiv1pb.RssFeed, error) {
	userID := ctx.Value(userIDContextKey).(int32)

	if request.RssFeed == nil {
		return nil, status.Errorf(codes.InvalidArgument, "RSS feed is required")
	}

	if request.RssFeed.Url == "" {
		return nil, status.Errorf(codes.InvalidArgument, "RSS feed URL is required")
	}

	// Validate RSS feed URL by attempting to fetch it
	if err := validateRssFeedURL(request.RssFeed.Url); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid RSS feed URL: %v", err)
	}

	// Set defaults
	if request.RssFeed.ImportFrequencyHours <= 0 {
		request.RssFeed.ImportFrequencyHours = 24 // Default to daily
	}

	create := &store.RssFeed{
		CreatorID:             userID,
		Title:                 request.RssFeed.Title,
		URL:                   request.RssFeed.Url,
		Description:           request.RssFeed.Description,
		AutoImport:            request.RssFeed.AutoImport,
		ImportFrequencyHours:  request.RssFeed.ImportFrequencyHours,
		DefaultTags:           request.RssFeed.DefaultTags,
		DefaultVisibility:     convertVisibilityToStore(request.RssFeed.DefaultVisibility),
		ShortcutPrefix:        request.RssFeed.ShortcutPrefix,
		IsActive:              true,
	}

	rssFeed, err := s.Store.CreateRssFeed(ctx, create)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create RSS feed: %v", err)
	}

	return convertRssFeedFromStore(rssFeed), nil
}

func (s *APIV1Service) UpdateRssFeed(ctx context.Context, request *apiv1pb.UpdateRssFeedRequest) (*apiv1pb.RssFeed, error) {
	userID := ctx.Value(userIDContextKey).(int32)

	if request.RssFeed == nil {
		return nil, status.Errorf(codes.InvalidArgument, "RSS feed is required")
	}

	// Check if the feed exists and belongs to the user
	existingFeed, err := s.Store.GetRssFeed(ctx, &store.FindRssFeed{
		ID:        &request.RssFeed.Id,
		CreatorID: &userID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get RSS feed: %v", err)
	}
	if existingFeed == nil {
		return nil, status.Errorf(codes.NotFound, "RSS feed not found")
	}

	// Validate URL if it's being updated
	if request.RssFeed.Url != "" && request.RssFeed.Url != existingFeed.URL {
		if err := validateRssFeedURL(request.RssFeed.Url); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid RSS feed URL: %v", err)
		}
	}

	update := &store.UpdateRssFeed{
		ID:                   request.RssFeed.Id,
		Title:                &request.RssFeed.Title,
		URL:                  &request.RssFeed.Url,
		Description:          &request.RssFeed.Description,
		AutoImport:           &request.RssFeed.AutoImport,
		ImportFrequencyHours: &request.RssFeed.ImportFrequencyHours,
		DefaultTags:          request.RssFeed.DefaultTags,
		ShortcutPrefix:       &request.RssFeed.ShortcutPrefix,
		IsActive:             &request.RssFeed.IsActive,
	}

	if request.RssFeed.DefaultVisibility != apiv1pb.Visibility_VISIBILITY_UNSPECIFIED {
		visibility := convertVisibilityToStore(request.RssFeed.DefaultVisibility)
		update.DefaultVisibility = &visibility
	}

	rssFeed, err := s.Store.UpdateRssFeed(ctx, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update RSS feed: %v", err)
	}

	return convertRssFeedFromStore(rssFeed), nil
}

func (s *APIV1Service) DeleteRssFeed(ctx context.Context, request *apiv1pb.DeleteRssFeedRequest) (*emptypb.Empty, error) {
	userID := ctx.Value(userIDContextKey).(int32)

	// Check if the feed exists and belongs to the user
	existingFeed, err := s.Store.GetRssFeed(ctx, &store.FindRssFeed{
		ID:        &request.Id,
		CreatorID: &userID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get RSS feed: %v", err)
	}
	if existingFeed == nil {
		return nil, status.Errorf(codes.NotFound, "RSS feed not found")
	}

	if err := s.Store.DeleteRssFeed(ctx, &store.DeleteRssFeed{ID: request.Id}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete RSS feed: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *APIV1Service) TriggerRssFeedImport(ctx context.Context, request *apiv1pb.TriggerRssFeedImportRequest) (*apiv1pb.RssFeedImportResponse, error) {
	userID := ctx.Value(userIDContextKey).(int32)

	// Check if the feed exists and belongs to the user
	rssFeed, err := s.Store.GetRssFeed(ctx, &store.FindRssFeed{
		ID:        &request.Id,
		CreatorID: &userID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get RSS feed: %v", err)
	}
	if rssFeed == nil {
		return nil, status.Errorf(codes.NotFound, "RSS feed not found")
	}

	// Import RSS feed
	importResult, err := s.importRssFeed(ctx, rssFeed)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to import RSS feed: %v", err)
	}

	return importResult, nil
}

// RSS parsing structures
type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
}

func validateRssFeedURL(url string) error {
	// Basic URL validation
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("URL must start with http:// or https://")
	}

	// Try to fetch the RSS feed
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to fetch RSS feed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("RSS feed returned status %d", resp.StatusCode)
	}

	// Try to parse as XML to validate it's a valid feed
	var rssFeed RSSFeed
	if err := xml.NewDecoder(resp.Body).Decode(&rssFeed); err != nil {
		return errors.Wrap(err, "invalid RSS feed format")
	}

	return nil
}

func (s *APIV1Service) importRssFeed(ctx context.Context, rssFeed *store.RssFeed) (*apiv1pb.RssFeedImportResponse, error) {
	// Fetch RSS feed
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(rssFeed.URL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch RSS feed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("RSS feed returned status %d", resp.StatusCode)
	}

	// Parse RSS feed
	var rss RSSFeed
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return nil, errors.Wrap(err, "failed to parse RSS feed")
	}

	response := &apiv1pb.RssFeedImportResponse{
		ImportedCount:     0,
		ImportedShortcuts: []string{},
		Errors:            []string{},
	}

	// Process each item
	for _, item := range rss.Channel.Items {
		// Check if we've already processed this item
		if item.GUID == "" {
			item.GUID = item.Link // Use link as GUID if no GUID provided
		}

		exists, err := s.Store.CheckRssFeedItemExists(ctx, rssFeed.ID, item.GUID)
		if err != nil {
			response.Errors = append(response.Errors, fmt.Sprintf("Failed to check if item exists: %v", err))
			continue
		}
		if exists {
			continue // Skip already processed items
		}

		// Generate shortcut name
		shortcutName := generateRssShortcutName(item.Title, rssFeed.ShortcutPrefix)

		// Check if shortcut name already exists
		existingShortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
			Name:      &shortcutName,
			CreatorID: &rssFeed.CreatorID,
		})
		if err == nil && existingShortcut != nil {
			// Make the name unique by appending timestamp
			shortcutName = fmt.Sprintf("%s-%d", shortcutName, time.Now().Unix())
		}

		// Create shortcut
		shortcut := &storepb.Shortcut{
			CreatorId:   rssFeed.CreatorID,
			Name:        shortcutName,
			Link:        item.Link,
			Title:       item.Title,
			Description: item.Description,
			Tags:        rssFeed.DefaultTags,
			Visibility:  rssFeed.DefaultVisibility,
		}

		createdShortcut, err := s.Store.CreateShortcut(ctx, shortcut)
		if err != nil {
			response.Errors = append(response.Errors, fmt.Sprintf("Failed to create shortcut for '%s': %v", item.Title, err))

			// Record failed import
			s.Store.CreateRssFeedItem(ctx, &store.RssFeedItem{
				RssFeedID:     rssFeed.ID,
				ItemGUID:      item.GUID,
				ItemLink:      item.Link,
				ItemTitle:     item.Title,
				ItemDescription: item.Description,
				ImportSuccess: false,
				ImportError:   err.Error(),
			})
			continue
		}

		// Record successful import
		s.Store.CreateRssFeedItem(ctx, &store.RssFeedItem{
			RssFeedID:     rssFeed.ID,
			ItemGUID:      item.GUID,
			ItemLink:      item.Link,
			ItemTitle:     item.Title,
			ItemDescription: item.Description,
			ShortcutID:    &createdShortcut.Id,
			ImportSuccess: true,
		})

		response.ImportedCount++
		response.ImportedShortcuts = append(response.ImportedShortcuts, shortcutName)
	}

	// Update RSS feed last import time
	now := time.Now()
	s.Store.UpdateRssFeed(ctx, &store.UpdateRssFeed{
		ID:             rssFeed.ID,
		LastImportTime: &now,
		TotalImported:  &[]int32{rssFeed.TotalImported + response.ImportedCount}[0],
	})

	return response, nil
}

func generateRssShortcutName(title string, prefix string) string {
	// Clean and format the title for use as shortcut name
	name := strings.ToLower(title)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	// Remove special characters except hyphens
	var result strings.Builder
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}
	name = result.String()

	// Remove multiple consecutive hyphens
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}

	// Trim hyphens from start and end
	name = strings.Trim(name, "-")

	// Limit length
	if len(name) > 50 {
		name = name[:50]
		name = strings.Trim(name, "-")
	}

	// Add prefix if provided
	if prefix != "" {
		name = prefix + "-" + name
	}

	// Ensure name is not empty
	if name == "" {
		name = "rss-item-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	return name
}

func convertRssFeedFromStore(rssFeed *store.RssFeed) *apiv1pb.RssFeed {
	result := &apiv1pb.RssFeed{
		Id:                   rssFeed.ID,
		CreatorId:            rssFeed.CreatorID,
		Title:                rssFeed.Title,
		Url:                  rssFeed.URL,
		Description:          rssFeed.Description,
		AutoImport:           rssFeed.AutoImport,
		ImportFrequencyHours: rssFeed.ImportFrequencyHours,
		DefaultTags:          rssFeed.DefaultTags,
		DefaultVisibility:    convertVisibilityFromStore(rssFeed.DefaultVisibility),
		ShortcutPrefix:       rssFeed.ShortcutPrefix,
		IsActive:             rssFeed.IsActive,
		LastError:            rssFeed.LastError,
		TotalImported:        rssFeed.TotalImported,
	}

	if rssFeed.CreatedTs != nil {
		result.CreatedTime = timestampFromUnix(rssFeed.CreatedTs.Unix())
	}
	if rssFeed.UpdatedTs != nil {
		result.UpdatedTime = timestampFromUnix(rssFeed.UpdatedTs.Unix())
	}
	if rssFeed.LastImportTs != nil {
		result.LastImportTime = timestampFromUnix(rssFeed.LastImportTs.Unix())
	}

	return result
}

func convertVisibilityToStore(visibility apiv1pb.Visibility) storepb.Visibility {
	switch visibility {
	case apiv1pb.Visibility_PUBLIC:
		return storepb.Visibility_PUBLIC
	case apiv1pb.Visibility_WORKSPACE:
		return storepb.Visibility_WORKSPACE
	default:
		return storepb.Visibility_WORKSPACE
	}
}

func convertVisibilityFromStore(visibility storepb.Visibility) apiv1pb.Visibility {
	switch visibility {
	case storepb.Visibility_PUBLIC:
		return apiv1pb.Visibility_PUBLIC
	case storepb.Visibility_WORKSPACE:
		return apiv1pb.Visibility_WORKSPACE
	default:
		return apiv1pb.Visibility_WORKSPACE
	}
}

// ImportRssFeed is a public method that can be called by the background runner
func (s *APIV1Service) ImportRssFeed(ctx context.Context, rssFeed *store.RssFeed) error {
	_, err := s.importRssFeed(ctx, rssFeed)
	return err
}