package v1

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mssola/useragent"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1pb "github.com/bshort/monotreme/proto/gen/api/v1"
	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/server/common"
	"github.com/bshort/monotreme/server/service/license"
	"github.com/bshort/monotreme/store"
)

func (s *APIV1Service) ListShortcuts(ctx context.Context, _ *v1pb.ListShortcutsRequest) (*v1pb.ListShortcutsResponse, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}

	// Try to get from Redis cache first
	if s.RedisService != nil && user != nil {
		var cachedResponse v1pb.ListShortcutsResponse
		cacheKey := fmt.Sprintf("shortcuts:user:%d", user.ID)
		cacheErr := s.RedisService.GetJSON(ctx, cacheKey, &cachedResponse)
		if cacheErr == nil {
			// Cache hit
			slog.Info("Redis cache hit", "key", cacheKey, "shortcut_count", len(cachedResponse.Shortcuts))
			return &cachedResponse, nil
		}
	}

	// Cache miss or Redis unavailable - fetch from database
	shortcutList, err := s.Store.ListShortcuts(ctx, &store.FindShortcut{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list shortcuts, err: %v", err)
	}

	// Batch-load all activity view counts at once to avoid N+1 query problem
	activities, err := s.Store.ListActivities(ctx, &store.FindActivity{
		Type:  store.ActivityShortcutView,
		Level: store.ActivityInfo,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list activities, err: %v", err)
	}

	// Build a map of shortcut ID to view count
	viewCountMap := make(map[int32]int32)
	for _, activity := range activities {
		payload := &storepb.ActivityShorcutViewPayload{}
		if err := protojson.Unmarshal([]byte(activity.Payload), payload); err != nil {
			continue // Skip invalid payloads
		}
		viewCountMap[payload.ShortcutId]++
	}

	shortcutMessageList := []*v1pb.Shortcut{}
	for _, shortcut := range shortcutList {
		composedShortcut := s.convertShortcutFromStorepbWithViewCount(shortcut, viewCountMap[shortcut.Id])
		shortcutMessageList = append(shortcutMessageList, composedShortcut)
	}

	response := &v1pb.ListShortcutsResponse{
		Shortcuts: shortcutMessageList,
	}

	// Save to Redis cache with 1 hour expiration
	if s.RedisService != nil && user != nil {
		_ = s.RedisService.SetJSON(ctx, fmt.Sprintf("shortcuts:user:%d", user.ID), response, time.Hour)
	}

	return response, nil
}

func (s *APIV1Service) GetShortcut(ctx context.Context, request *v1pb.GetShortcutRequest) (*v1pb.Shortcut, error) {
	shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
		ID: &request.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get shortcut by id: %v", err)
	}
	if shortcut == nil {
		return nil, status.Errorf(codes.NotFound, "shortcut not found")
	}

	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil && shortcut.Visibility != storepb.Visibility_PUBLIC {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	composedShortcut, err := s.convertShortcutFromStorepb(ctx, shortcut)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert shortcut, err: %v", err)
	}
	return composedShortcut, nil
}

func (s *APIV1Service) GetShortcutByName(ctx context.Context, request *v1pb.GetShortcutByNameRequest) (*v1pb.Shortcut, error) {
	shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
		Name: &request.Name,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get shortcut by name: %v", err)
	}
	if shortcut == nil {
		return nil, status.Errorf(codes.NotFound, "shortcut not found")
	}

	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil && shortcut.Visibility != storepb.Visibility_PUBLIC {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	composedShortcut, err := s.convertShortcutFromStorepb(ctx, shortcut)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert shortcut, err: %v", err)
	}
	return composedShortcut, nil
}

func (s *APIV1Service) CreateShortcut(ctx context.Context, request *v1pb.CreateShortcutRequest) (*v1pb.Shortcut, error) {
	if request.Shortcut.Name == "" || request.Shortcut.Link == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name and link are required")
	}


	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	shortcutCreate := &storepb.Shortcut{
		CreatorId:   user.ID,
		Name:        request.Shortcut.Name,
		Link:        request.Shortcut.Link,
		Title:       request.Shortcut.Title,
		Tags:        request.Shortcut.Tags,
		Description: request.Shortcut.Description,
		Visibility:  convertVisibilityToStorepb(request.Shortcut.Visibility),
		OgMetadata:  &storepb.OpenGraphMetadata{},
		Uuid:        uuid.New().String(),
	}
	if shortcutCreate.Visibility == storepb.Visibility_VISIBILITY_UNSPECIFIED {
		workspaceSetting, err := s.GetWorkspaceSetting(ctx, nil)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get workspace setting, err: %v", err)
		}
		visibility := v1pb.Visibility_WORKSPACE
		if workspaceSetting.DefaultVisibility != v1pb.Visibility_VISIBILITY_UNSPECIFIED {
			visibility = workspaceSetting.DefaultVisibility
		}
		shortcutCreate.Visibility = convertVisibilityToStorepb(visibility)
	}
	if request.Shortcut.OgMetadata != nil {
		shortcutCreate.OgMetadata = &storepb.OpenGraphMetadata{
			Title:       request.Shortcut.OgMetadata.Title,
			Description: request.Shortcut.OgMetadata.Description,
			Image:       request.Shortcut.OgMetadata.Image,
		}
	}
	shortcut, err := s.Store.CreateShortcut(ctx, shortcutCreate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create shortcut, err: %v", err)
	}

	// Process tags: create Tag entities and bookmark_tag relationships
	if err := s.processShortcutTags(ctx, shortcut.Id, user.ID, request.Shortcut.Tags); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process tags, err: %v", err)
	}

	if err := s.createShortcutCreateActivity(ctx, shortcut); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create activity, err: %v", err)
	}

	// Invalidate Redis cache for this user
	if s.RedisService != nil {
		_ = s.RedisService.InvalidateShortcutListCache(ctx, user.ID)
	}

	composedShortcut, err := s.convertShortcutFromStorepb(ctx, shortcut)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert shortcut, err: %v", err)
	}
	return composedShortcut, nil
}

func (s *APIV1Service) UpdateShortcut(ctx context.Context, request *v1pb.UpdateShortcutRequest) (*v1pb.Shortcut, error) {
	if request.UpdateMask == nil || len(request.UpdateMask.Paths) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "updateMask is required")
	}

	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
		ID: &request.Shortcut.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get shortcut by id: %v", err)
	}
	if shortcut == nil {
		return nil, status.Errorf(codes.NotFound, "shortcut not found")
	}
	if shortcut.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	update := &store.UpdateShortcut{
		ID: shortcut.Id,
	}
	tagsUpdated := false
	for _, path := range request.UpdateMask.Paths {
		switch path {
		case "name":
			update.Name = &request.Shortcut.Name
		case "link":
			update.Link = &request.Shortcut.Link
		case "title":
			update.Title = &request.Shortcut.Title
		case "description":
			update.Description = &request.Shortcut.Description
		case "tags":
			tag := strings.Join(request.Shortcut.Tags, " ")
			update.Tag = &tag
			tagsUpdated = true
		case "visibility":
			visibility := convertVisibilityToStorepb(request.Shortcut.Visibility)
			update.Visibility = &visibility
		case "og_metadata":
			if request.Shortcut.OgMetadata != nil {
				update.OpenGraphMetadata = &storepb.OpenGraphMetadata{
					Title:       request.Shortcut.OgMetadata.Title,
					Description: request.Shortcut.OgMetadata.Description,
					Image:       request.Shortcut.OgMetadata.Image,
				}
			}
		case "user_order":
			update.UserOrder = &request.Shortcut.UserOrder
		}
	}
	shortcut, err = s.Store.UpdateShortcut(ctx, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update shortcut, err: %v", err)
	}

	// If tags were updated, update Tag entities and relationships
	if tagsUpdated {
		// Delete existing bookmark_tag relationships
		if err := s.Store.DeleteBookmarkTagsByShortcut(ctx, shortcut.Id); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to delete old tag relationships, err: %v", err)
		}
		// Create new relationships
		if err := s.processShortcutTags(ctx, shortcut.Id, user.ID, request.Shortcut.Tags); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to process tags, err: %v", err)
		}
	}

	// Invalidate Redis cache for this user
	if s.RedisService != nil {
		_ = s.RedisService.InvalidateShortcutListCache(ctx, user.ID)
	}

	composedShortcut, err := s.convertShortcutFromStorepb(ctx, shortcut)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert shortcut, err: %v", err)
	}
	return composedShortcut, nil
}

func (s *APIV1Service) DeleteShortcut(ctx context.Context, request *v1pb.DeleteShortcutRequest) (*emptypb.Empty, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
		ID: &request.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get shortcut by id: %v", err)
	}
	if shortcut == nil {
		return nil, status.Errorf(codes.NotFound, "shortcut not found")
	}
	if shortcut.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	err = s.Store.DeleteShortcut(ctx, &store.DeleteShortcut{
		ID: shortcut.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete shortcut, err: %v", err)
	}

	// Invalidate Redis cache for this user
	if s.RedisService != nil {
		_ = s.RedisService.InvalidateShortcutListCache(ctx, user.ID)
	}

	return &emptypb.Empty{}, nil
}

func (s *APIV1Service) DeleteAllShortcuts(ctx context.Context, _ *v1pb.DeleteAllShortcutsRequest) (*emptypb.Empty, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}

	// Get all shortcuts for the current user
	shortcuts, err := s.Store.ListShortcuts(ctx, &store.FindShortcut{
		CreatorID: &user.ID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list shortcuts, err: %v", err)
	}

	// Delete each shortcut
	for _, shortcut := range shortcuts {
		err = s.Store.DeleteShortcut(ctx, &store.DeleteShortcut{
			ID: shortcut.Id,
		})
		if err != nil {
			slog.Error("Failed to delete shortcut", "shortcut_id", shortcut.Id, "error", err)
			return nil, status.Errorf(codes.Internal, "failed to delete shortcut %d, err: %v", shortcut.Id, err)
		}
	}

	slog.Info("Deleted all shortcuts for user", "user_id", user.ID, "count", len(shortcuts))

	// Invalidate Redis cache for this user
	if s.RedisService != nil {
		_ = s.RedisService.InvalidateShortcutListCache(ctx, user.ID)
	}

	return &emptypb.Empty{}, nil
}

func (s *APIV1Service) GetShortcutAnalytics(ctx context.Context, request *v1pb.GetShortcutAnalyticsRequest) (*v1pb.GetShortcutAnalyticsResponse, error) {
	shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
		ID: &request.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get shortcut by id: %v", err)
	}
	if shortcut == nil {
		return nil, status.Errorf(codes.NotFound, "shortcut not found")
	}

	activityFind := &store.FindActivity{
		Type:              store.ActivityShortcutView,
		PayloadShortcutID: &shortcut.Id,
	}
	// For non-advanced analytics users, we limit the activity to the last 14 days.
	if !s.LicenseService.IsFeatureEnabled(license.FeatureTypeAdvancedAnalytics) {
		createdTsAfter := time.Now().AddDate(0, 0, -14).Unix()
		activityFind.CreatedTsAfter = &createdTsAfter
	}
	activities, err := s.Store.ListActivities(ctx, activityFind)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get activities, err: %v", err)
	}

	referenceMap := make(map[string]int32)
	deviceMap := make(map[string]int32)
	browserMap := make(map[string]int32)
	for _, activity := range activities {
		payload := &storepb.ActivityShorcutViewPayload{}
		if err := protojson.Unmarshal([]byte(activity.Payload), payload); err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to unmarshal payload, err: %v", err))
		}

		if _, ok := referenceMap[payload.Referer]; !ok {
			referenceMap[payload.Referer] = 0
		}
		referenceMap[payload.Referer]++

		ua := useragent.New(payload.UserAgent)
		deviceName := ua.OSInfo().Name
		browserName, _ := ua.Browser()

		if _, ok := deviceMap[deviceName]; !ok {
			deviceMap[deviceName] = 0
		}
		deviceMap[deviceName]++

		if _, ok := browserMap[browserName]; !ok {
			browserMap[browserName] = 0
		}
		browserMap[browserName]++
	}

	response := &v1pb.GetShortcutAnalyticsResponse{
		References: mapToAnalyticsSlice(referenceMap),
		Devices:    mapToAnalyticsSlice(deviceMap),
		Browsers:   mapToAnalyticsSlice(browserMap),
	}
	return response, nil
}

func mapToAnalyticsSlice(m map[string]int32) []*v1pb.GetShortcutAnalyticsResponse_AnalyticsItem {
	analyticsSlice := make([]*v1pb.GetShortcutAnalyticsResponse_AnalyticsItem, 0)
	for key, value := range m {
		analyticsSlice = append(analyticsSlice, &v1pb.GetShortcutAnalyticsResponse_AnalyticsItem{
			Name:  key,
			Count: value,
		})
	}
	slices.SortFunc(analyticsSlice, func(i, j *v1pb.GetShortcutAnalyticsResponse_AnalyticsItem) int {
		return int(i.Count - j.Count)
	})
	return analyticsSlice
}

func (s *APIV1Service) createShortcutCreateActivity(ctx context.Context, shortcut *storepb.Shortcut) error {
	payload := &storepb.ActivityShorcutCreatePayload{
		ShortcutId: shortcut.Id,
	}
	payloadStr, err := protojson.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal activity payload")
	}
	activity := &store.Activity{
		CreatorID: shortcut.CreatorId,
		Type:      store.ActivityShortcutCreate,
		Level:     store.ActivityInfo,
		Payload:   string(payloadStr),
	}
	_, err = s.Store.CreateActivity(ctx, activity)
	if err != nil {
		return errors.Wrap(err, "Failed to create activity")
	}
	return nil
}

func (s *APIV1Service) convertShortcutFromStorepb(ctx context.Context, shortcut *storepb.Shortcut) (*v1pb.Shortcut, error) {
	composedShortcut := &v1pb.Shortcut{
		Id:          shortcut.Id,
		CreatorId:   shortcut.CreatorId,
		CreatedTime: timestamppb.New(time.Unix(shortcut.CreatedTs, 0)),
		UpdatedTime: timestamppb.New(time.Unix(shortcut.UpdatedTs, 0)),
		Name:        shortcut.Name,
		Link:        shortcut.Link,
		Title:       shortcut.Title,
		Tags:        shortcut.Tags,
		Description: shortcut.Description,
		Visibility:  convertVisibilityFromStorepb(shortcut.Visibility),
		Uuid:        shortcut.Uuid,
		UserOrder:   shortcut.UserOrder,
		OgMetadata: &v1pb.Shortcut_OpenGraphMetadata{
			Title:       shortcut.OgMetadata.Title,
			Description: shortcut.OgMetadata.Description,
			Image:       shortcut.OgMetadata.Image,
		},
	}

	activityList, err := s.Store.ListActivities(ctx, &store.FindActivity{
		Type:              store.ActivityShortcutView,
		Level:             store.ActivityInfo,
		PayloadShortcutID: &composedShortcut.Id,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list activities")
	}
	composedShortcut.ViewCount = int32(len(activityList))

	// Populate tag_info with UUIDs and names from Tag entities
	tagInfo, err := s.getTagInfoForShortcut(ctx, shortcut.Id)
	if err != nil {
		// Don't fail the entire request if tag lookup fails, just log it
		return composedShortcut, nil
	}
	composedShortcut.TagInfo = tagInfo

	return composedShortcut, nil
}

// convertShortcutFromStorepbWithViewCount converts a shortcut from storepb to v1pb with a pre-computed view count
// This avoids the N+1 query problem when listing many shortcuts
func (s *APIV1Service) convertShortcutFromStorepbWithViewCount(shortcut *storepb.Shortcut, viewCount int32) *v1pb.Shortcut {
	composedShortcut := &v1pb.Shortcut{
		Id:          shortcut.Id,
		CreatorId:   shortcut.CreatorId,
		CreatedTime: timestamppb.New(time.Unix(shortcut.CreatedTs, 0)),
		UpdatedTime: timestamppb.New(time.Unix(shortcut.UpdatedTs, 0)),
		Name:        shortcut.Name,
		Link:        shortcut.Link,
		Title:       shortcut.Title,
		Tags:        shortcut.Tags,
		Description: shortcut.Description,
		Visibility:  convertVisibilityFromStorepb(shortcut.Visibility),
		Uuid:        shortcut.Uuid,
		UserOrder:   shortcut.UserOrder,
		ViewCount:   viewCount,
		OgMetadata: &v1pb.Shortcut_OpenGraphMetadata{
			Title:       shortcut.OgMetadata.Title,
			Description: shortcut.OgMetadata.Description,
			Image:       shortcut.OgMetadata.Image,
		},
	}
	return composedShortcut
}

// processShortcutTags processes tag strings and creates Tag entities and bookmark_tag relationships
func (s *APIV1Service) processShortcutTags(ctx context.Context, shortcutID, userID int32, tagNames []string) error {
	for _, tagName := range tagNames {
		if tagName == "" {
			continue
		}

		// Generate abbreviation for the tag
		abbreviation := common.GenerateTagAbbreviation(tagName)

		// Try to find existing tag by abbreviation and user
		tag, err := s.Store.GetTag(ctx, &store.FindTag{
			CreatorID:    &userID,
			Abbreviation: &abbreviation,
		})
		if err != nil {
			return err
		}

		// If tag doesn't exist, create it
		if tag == nil {
			newTag := &storepb.Tag{
				Uuid:         uuid.New().String(),
				CreatorId:    userID,
				Name:         tagName,
				Abbreviation: abbreviation,
				Description:  "",
			}
			tag, err = s.Store.CreateTag(ctx, newTag)
			if err != nil {
				return err
			}
		}

		// Create bookmark_tag relationship
		_, err = s.Store.CreateBookmarkTag(ctx, &store.CreateBookmarkTag{
			ShortcutID: shortcutID,
			TagUUID:    tag.Uuid,
			UserID:     userID,
		})
		if err != nil {
			// Ignore duplicate errors (relationship already exists)
			continue
		}
	}

	return nil
}

// getTagInfoForShortcut retrieves tag UUIDs and names for a shortcut
func (s *APIV1Service) getTagInfoForShortcut(ctx context.Context, shortcutID int32) ([]*v1pb.Shortcut_TagInfo, error) {
	// Get bookmark_tag relationships
	bookmarkTags, err := s.Store.ListBookmarkTags(ctx, &store.FindBookmarkTag{
		ShortcutID: &shortcutID,
	})
	if err != nil {
		return nil, err
	}

	// Fetch tag details
	tagInfoList := make([]*v1pb.Shortcut_TagInfo, 0, len(bookmarkTags))
	for _, bt := range bookmarkTags {
		tag, err := s.Store.GetTag(ctx, &store.FindTag{
			UUID: &bt.TagUuid,
		})
		if err != nil || tag == nil {
			continue
		}

		tagInfoList = append(tagInfoList, &v1pb.Shortcut_TagInfo{
			Uuid: tag.Uuid,
			Name: tag.Name,
		})
	}

	return tagInfoList, nil
}
