package v1

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1pb "github.com/bshort/monotreme/proto/gen/api/v1"
	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

const (
	DefaultRecentActivityLimit = 5
	DefaultPageSize           = 20
	MaxPageSize               = 100
)

func (s *APIV1Service) GetRecentActivity(ctx context.Context, request *v1pb.GetRecentActivityRequest) (*v1pb.GetRecentActivityResponse, error) {
	limit := DefaultRecentActivityLimit
	if request.Limit != nil && *request.Limit > 0 {
		limit = int(*request.Limit)
	}

	response := &v1pb.GetRecentActivityResponse{}

	// Get recent users (most recently created)
	recentUsers, err := s.getRecentUsers(ctx, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get recent users: %v", err)
	}
	response.RecentUsers = recentUsers

	// Get recent shortcuts (most recently created)
	recentShortcuts, err := s.getRecentShortcuts(ctx, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get recent shortcuts: %v", err)
	}
	response.RecentShortcuts = recentShortcuts

	// Get recent collections (most recently created)
	recentCollections, err := s.getRecentCollections(ctx, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get recent collections: %v", err)
	}
	response.RecentCollections = recentCollections

	// Get recent clicks (most recent shortcut views)
	recentClicks, err := s.getRecentClicks(ctx, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get recent clicks: %v", err)
	}
	response.RecentClicks = recentClicks

	// Get most clicked shortcuts
	mostClickedShortcuts, err := s.getMostClickedShortcuts(ctx, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get most clicked shortcuts: %v", err)
	}
	response.MostClickedShortcuts = mostClickedShortcuts

	return response, nil
}

func (s *APIV1Service) GetActivitySummary(ctx context.Context, request *v1pb.GetActivitySummaryRequest) (*v1pb.GetActivitySummaryResponse, error) {
	currentUser, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}

	response := &v1pb.GetActivitySummaryResponse{}

	// Get total counts
	users, err := s.Store.ListUsers(ctx, &store.FindUser{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}
	response.TotalUsers = int32(len(users))

	shortcuts, err := s.Store.ListShortcuts(ctx, &store.FindShortcut{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list shortcuts: %v", err)
	}
	response.TotalShortcuts = int32(len(shortcuts))

	collections, err := s.Store.ListCollections(ctx, &store.FindCollection{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list collections: %v", err)
	}
	response.TotalCollections = int32(len(collections))

	// Get total clicks
	activities, err := s.Store.ListActivities(ctx, &store.FindActivity{
		Type: store.ActivityShortcutView,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list activities: %v", err)
	}
	response.TotalClicks = int32(len(activities))

	// Get recent activity counts (last 24 hours)
	oneDayAgo := time.Now().AddDate(0, 0, -1).Unix()

	recentUsers, err := s.Store.ListUsers(ctx, &store.FindUser{
		CreatedTsAfter: &oneDayAgo,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list recent users: %v", err)
	}
	response.RecentUsersCount = int32(len(recentUsers))

	recentShortcuts, err := s.Store.ListShortcuts(ctx, &store.FindShortcut{
		CreatedTsAfter: &oneDayAgo,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list recent shortcuts: %v", err)
	}
	response.RecentShortcutsCount = int32(len(recentShortcuts))

	recentCollections, err := s.Store.ListCollections(ctx, &store.FindCollection{
		CreatedTsAfter: &oneDayAgo,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list recent collections: %v", err)
	}
	response.RecentCollectionsCount = int32(len(recentCollections))

	recentClickActivities, err := s.Store.ListActivities(ctx, &store.FindActivity{
		Type:           store.ActivityShortcutView,
		CreatedTsAfter: &oneDayAgo,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list recent click activities: %v", err)
	}
	response.RecentClicksCount = int32(len(recentClickActivities))

	// Get user-specific summary
	if currentUser != nil {
		userSummary, err := s.getUserSummary(ctx, currentUser.Id)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get user summary: %v", err)
		}
		response.UserSummary = userSummary
	}

	return response, nil
}

func (s *APIV1Service) ListActivities(ctx context.Context, request *v1pb.ListActivitiesRequest) (*v1pb.ListActivitiesResponse, error) {
	pageSize := DefaultPageSize
	if request.PageSize > 0 {
		pageSize = int(request.PageSize)
		if pageSize > MaxPageSize {
			pageSize = MaxPageSize
		}
	}

	// Parse page token to get offset
	offset := 0
	if request.PageToken != "" {
		if parsedOffset, err := strconv.Atoi(request.PageToken); err == nil {
			offset = parsedOffset
		}
	}

	// Build activity filter
	findActivity := &store.FindActivity{}

	// Add activity type filter
	if request.ActivityType != nil && *request.ActivityType != v1pb.ActivityType_ACTIVITY_TYPE_UNSPECIFIED {
		switch *request.ActivityType {
		case v1pb.ActivityType_SHORTCUT_CREATED:
			findActivity.Type = store.ActivityShortcutCreate
		case v1pb.ActivityType_SHORTCUT_VIEWED:
			findActivity.Type = store.ActivityShortcutView
		}
	}

	// Add time range filters
	if request.CreatedAfter != nil {
		createdAfter := request.CreatedAfter.AsTime().Unix()
		findActivity.CreatedTsAfter = &createdAfter
	}

	// Get activities with pagination
	activities, err := s.Store.ListActivities(ctx, findActivity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list activities: %v", err)
	}

	// Apply pagination
	totalCount := len(activities)
	start := offset
	end := offset + pageSize
	if start > totalCount {
		start = totalCount
	}
	if end > totalCount {
		end = totalCount
	}

	paginatedActivities := activities[start:end]

	// Convert activities to response format
	activityItems := make([]*v1pb.ActivityItem, len(paginatedActivities))
	for i, activity := range paginatedActivities {
		activityItem, err := s.convertActivityToActivityItem(ctx, activity)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert activity: %v", err)
		}
		activityItems[i] = activityItem
	}

	// Calculate next page token
	nextPageToken := ""
	if end < totalCount {
		nextPageToken = strconv.Itoa(end)
	}

	return &v1pb.ListActivitiesResponse{
		Activities:    activityItems,
		NextPageToken: nextPageToken,
		TotalCount:    int32(totalCount),
	}, nil
}

// Helper functions

func (s *APIV1Service) getRecentUsers(ctx context.Context, limit int) ([]*v1pb.RecentUser, error) {
	users, err := s.Store.ListUsers(ctx, &store.FindUser{})
	if err != nil {
		return nil, err
	}

	// Sort by created time (most recent first)
	sortedUsers := make([]*storepb.User, len(users))
	copy(sortedUsers, users)

	// Simple bubble sort by CreatedTs (newest first)
	for i := 0; i < len(sortedUsers); i++ {
		for j := i + 1; j < len(sortedUsers); j++ {
			if sortedUsers[i].CreatedTs < sortedUsers[j].CreatedTs {
				sortedUsers[i], sortedUsers[j] = sortedUsers[j], sortedUsers[i]
			}
		}
	}

	// Take the most recent ones
	if len(sortedUsers) > limit {
		sortedUsers = sortedUsers[:limit]
	}

	recentUsers := make([]*v1pb.RecentUser, len(sortedUsers))
	for i, user := range sortedUsers {
		recentUsers[i] = &v1pb.RecentUser{
			Id:          user.Id,
			Email:       user.Email,
			Nickname:    user.Nickname,
			CreatedTime: timestampFromUnix(user.CreatedTs),
			Role:        convertRoleFromStore(user.Role),
		}
	}

	return recentUsers, nil
}

func (s *APIV1Service) getRecentShortcuts(ctx context.Context, limit int) ([]*v1pb.RecentShortcut, error) {
	shortcuts, err := s.Store.ListShortcuts(ctx, &store.FindShortcut{})
	if err != nil {
		return nil, err
	}

	// Sort by created time (most recent first)
	sortedShortcuts := make([]*storepb.Shortcut, len(shortcuts))
	copy(sortedShortcuts, shortcuts)

	for i := 0; i < len(sortedShortcuts); i++ {
		for j := i + 1; j < len(sortedShortcuts); j++ {
			if sortedShortcuts[i].CreatedTs < sortedShortcuts[j].CreatedTs {
				sortedShortcuts[i], sortedShortcuts[j] = sortedShortcuts[j], sortedShortcuts[i]
			}
		}
	}

	if len(sortedShortcuts) > limit {
		sortedShortcuts = sortedShortcuts[:limit]
	}

	recentShortcuts := make([]*v1pb.RecentShortcut, len(sortedShortcuts))
	for i, shortcut := range sortedShortcuts {
		creator, _ := s.Store.GetUser(ctx, &store.FindUser{ID: &shortcut.CreatorId})
		creatorName := ""
		if creator != nil {
			creatorName = creator.Nickname
			if creatorName == "" {
				creatorName = creator.Email
			}
		}

		recentShortcuts[i] = &v1pb.RecentShortcut{
			Id:          shortcut.Id,
			Name:        shortcut.Name,
			Title:       shortcut.Title,
			Link:        shortcut.Link,
			CreatorId:   shortcut.CreatorId,
			CreatorName: creatorName,
			CreatedTime: timestampFromUnix(shortcut.CreatedTs),
			ViewCount:   shortcut.ViewCount,
			Tags:        shortcut.Tags,
		}
	}

	return recentShortcuts, nil
}

func (s *APIV1Service) getRecentCollections(ctx context.Context, limit int) ([]*v1pb.RecentCollection, error) {
	collections, err := s.Store.ListCollections(ctx, &store.FindCollection{})
	if err != nil {
		return nil, err
	}

	// Sort by created time (most recent first)
	sortedCollections := make([]*storepb.Collection, len(collections))
	copy(sortedCollections, collections)

	for i := 0; i < len(sortedCollections); i++ {
		for j := i + 1; j < len(sortedCollections); j++ {
			if sortedCollections[i].CreatedTs < sortedCollections[j].CreatedTs {
				sortedCollections[i], sortedCollections[j] = sortedCollections[j], sortedCollections[i]
			}
		}
	}

	if len(sortedCollections) > limit {
		sortedCollections = sortedCollections[:limit]
	}

	recentCollections := make([]*v1pb.RecentCollection, len(sortedCollections))
	for i, collection := range sortedCollections {
		creator, _ := s.Store.GetUser(ctx, &store.FindUser{ID: &collection.CreatorId})
		creatorName := ""
		if creator != nil {
			creatorName = creator.Nickname
			if creatorName == "" {
				creatorName = creator.Email
			}
		}

		// Count shortcuts in collection
		shortcuts, _ := s.Store.ListShortcuts(ctx, &store.FindShortcut{
			CollectionIds: []int32{collection.Id},
		})
		shortcutCount := int32(0)
		if shortcuts != nil {
			shortcutCount = int32(len(shortcuts))
		}

		recentCollections[i] = &v1pb.RecentCollection{
			Id:            collection.Id,
			Name:          collection.Name,
			Title:         collection.Title,
			Description:   collection.Description,
			CreatorId:     collection.CreatorId,
			CreatorName:   creatorName,
			CreatedTime:   timestampFromUnix(collection.CreatedTs),
			ShortcutCount: shortcutCount,
		}
	}

	return recentCollections, nil
}

func (s *APIV1Service) getRecentClicks(ctx context.Context, limit int) ([]*v1pb.RecentClick, error) {
	activities, err := s.Store.ListActivities(ctx, &store.FindActivity{
		Type: store.ActivityShortcutView,
	})
	if err != nil {
		return nil, err
	}

	// Sort by created time (most recent first)
	for i := 0; i < len(activities); i++ {
		for j := i + 1; j < len(activities); j++ {
			if activities[i].CreatedTs < activities[j].CreatedTs {
				activities[i], activities[j] = activities[j], activities[i]
			}
		}
	}

	if len(activities) > limit {
		activities = activities[:limit]
	}

	recentClicks := make([]*v1pb.RecentClick, 0, len(activities))
	for _, activity := range activities {
		payload := &storepb.ActivityShorcutViewPayload{}
		if err := protojson.Unmarshal([]byte(activity.Payload), payload); err != nil {
			continue // Skip invalid payloads
		}

		shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{ID: &payload.ShortcutId})
		if err != nil || shortcut == nil {
			continue // Skip if shortcut not found
		}

		creator, _ := s.Store.GetUser(ctx, &store.FindUser{ID: &shortcut.CreatorId})
		creatorName := ""
		if creator != nil {
			creatorName = creator.Nickname
			if creatorName == "" {
				creatorName = creator.Email
			}
		}

		recentClicks = append(recentClicks, &v1pb.RecentClick{
			ShortcutId:    shortcut.Id,
			ShortcutName:  shortcut.Name,
			ShortcutTitle: shortcut.Title,
			ShortcutLink:  shortcut.Link,
			CreatorId:     shortcut.CreatorId,
			CreatorName:   creatorName,
			ClickedTime:   timestampFromUnix(activity.CreatedTs),
			UserAgent:     payload.UserAgent,
			Referer:       payload.Referer,
		})
	}

	return recentClicks, nil
}

func (s *APIV1Service) getMostClickedShortcuts(ctx context.Context, limit int) ([]*v1pb.MostClickedShortcut, error) {
	shortcuts, err := s.Store.ListShortcuts(ctx, &store.FindShortcut{})
	if err != nil {
		return nil, err
	}

	// Filter shortcuts with clicks and sort by view count (descending)
	clickedShortcuts := make([]*storepb.Shortcut, 0)
	for _, shortcut := range shortcuts {
		if shortcut.ViewCount > 0 {
			clickedShortcuts = append(clickedShortcuts, shortcut)
		}
	}

	// Sort by view count (highest first)
	for i := 0; i < len(clickedShortcuts); i++ {
		for j := i + 1; j < len(clickedShortcuts); j++ {
			if clickedShortcuts[i].ViewCount < clickedShortcuts[j].ViewCount {
				clickedShortcuts[i], clickedShortcuts[j] = clickedShortcuts[j], clickedShortcuts[i]
			}
		}
	}

	if len(clickedShortcuts) > limit {
		clickedShortcuts = clickedShortcuts[:limit]
	}

	mostClickedShortcuts := make([]*v1pb.MostClickedShortcut, len(clickedShortcuts))
	for i, shortcut := range clickedShortcuts {
		creator, _ := s.Store.GetUser(ctx, &store.FindUser{ID: &shortcut.CreatorId})
		creatorName := ""
		if creator != nil {
			creatorName = creator.Nickname
			if creatorName == "" {
				creatorName = creator.Email
			}
		}

		mostClickedShortcuts[i] = &v1pb.MostClickedShortcut{
			Id:          shortcut.Id,
			Name:        shortcut.Name,
			Title:       shortcut.Title,
			Link:        shortcut.Link,
			CreatorId:   shortcut.CreatorId,
			CreatorName: creatorName,
			ViewCount:   shortcut.ViewCount,
			LastClicked: timestampFromUnix(shortcut.UpdatedTs), // Using updated timestamp as proxy
		}
	}

	return mostClickedShortcuts, nil
}

func (s *APIV1Service) getUserSummary(ctx context.Context, userID int32) (*v1pb.UserSummary, error) {
	// Get user's shortcuts
	userShortcuts, err := s.Store.ListShortcuts(ctx, &store.FindShortcut{
		CreatorID: &userID,
	})
	if err != nil {
		return nil, err
	}

	// Get user's collections
	userCollections, err := s.Store.ListCollections(ctx, &store.FindCollection{
		CreatorID: &userID,
	})
	if err != nil {
		return nil, err
	}

	// Calculate total clicks for user's shortcuts
	totalClicks := int32(0)
	tagMap := make(map[string]bool)

	for _, shortcut := range userShortcuts {
		totalClicks += shortcut.ViewCount
		for _, tag := range shortcut.Tags {
			tagMap[tag] = true
		}
	}

	// Convert tag map to slice
	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}

	return &v1pb.UserSummary{
		UserShortcutsCount:   int32(len(userShortcuts)),
		UserCollectionsCount: int32(len(userCollections)),
		UserTotalClicks:      totalClicks,
		UserTags:             tags,
	}, nil
}

func (s *APIV1Service) convertActivityToActivityItem(ctx context.Context, activity *store.Activity) (*v1pb.ActivityItem, error) {
	user, err := s.Store.GetUser(ctx, &store.FindUser{ID: &activity.CreatorID})
	if err != nil {
		return nil, err
	}

	userName := ""
	if user != nil {
		userName = user.Nickname
		if userName == "" {
			userName = user.Email
		}
	}

	activityItem := &v1pb.ActivityItem{
		Id:          activity.ID,
		UserId:      activity.CreatorID,
		UserName:    userName,
		CreatedTime: timestampFromUnix(activity.CreatedTs),
	}

	// Set activity type and data based on store activity type
	switch activity.Type {
	case store.ActivityShortcutCreate:
		activityItem.Type = v1pb.ActivityType_SHORTCUT_CREATED
		// Parse payload to get shortcut data
		payload := &storepb.ActivityShorcutCreatePayload{}
		if err := protojson.Unmarshal([]byte(activity.Payload), payload); err == nil {
			shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{ID: &payload.ShortcutId})
			if err == nil && shortcut != nil {
				activityItem.Data = &v1pb.ActivityItem_ShortcutCreated{
					ShortcutCreated: &v1pb.ShortcutCreatedData{
						ShortcutId: shortcut.Id,
						Name:       shortcut.Name,
						Title:      shortcut.Title,
						Link:       shortcut.Link,
					},
				}
			}
		}

	case store.ActivityShortcutView:
		activityItem.Type = v1pb.ActivityType_SHORTCUT_VIEWED
		// Parse payload to get view data
		payload := &storepb.ActivityShorcutViewPayload{}
		if err := protojson.Unmarshal([]byte(activity.Payload), payload); err == nil {
			shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{ID: &payload.ShortcutId})
			if err == nil && shortcut != nil {
				activityItem.Data = &v1pb.ActivityItem_ShortcutViewed{
					ShortcutViewed: &v1pb.ShortcutViewedData{
						ShortcutId: shortcut.Id,
						Name:       shortcut.Name,
						Title:      shortcut.Title,
						UserAgent:  payload.UserAgent,
						Referer:    payload.Referer,
					},
				}
			}
		}
	}

	return activityItem, nil
}

// Helper functions for type conversion
func convertRoleFromStore(role storepb.User_Role) v1pb.Role {
	switch role {
	case storepb.User_ADMIN:
		return v1pb.Role_ADMIN
	case storepb.User_USER:
		return v1pb.Role_USER
	default:
		return v1pb.Role_ROLE_UNSPECIFIED
	}
}

func timestampFromUnix(ts int64) *timestamppb.Timestamp {
	return timestamppb.New(time.Unix(ts, 0))
}