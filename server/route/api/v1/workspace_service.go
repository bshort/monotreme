package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1pb "github.com/bshort/monotreme/proto/gen/api/v1"
	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

func (s *APIV1Service) GetWorkspaceProfile(ctx context.Context, _ *v1pb.GetWorkspaceProfileRequest) (*v1pb.WorkspaceProfile, error) {
	workspaceProfile := &v1pb.WorkspaceProfile{
		Mode:    s.Profile.Mode,
		Version: s.Profile.Version,
	}

	// Load subscription plan from license service.
	subscription := s.LicenseService.GetSubscription()
	workspaceProfile.Subscription = subscription

	owner, err := s.GetInstanceOwner(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get instance owner: %v", err)
	}
	if owner != nil {
		workspaceProfile.Owner = fmt.Sprintf("%s%d", UserNamePrefix, owner.Id)
	}

	workspaceGeneralSetting, err := s.Store.GetWorkspaceGeneralSetting(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get workspace general setting")
	}
	workspaceProfile.Branding = workspaceGeneralSetting.GetBranding()

	return workspaceProfile, nil
}

func (s *APIV1Service) GetWorkspaceSetting(ctx context.Context, _ *v1pb.GetWorkspaceSettingRequest) (*v1pb.WorkspaceSetting, error) {
	currentUser, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	workspaceSettings, err := s.Store.ListWorkspaceSettings(ctx, &store.FindWorkspaceSetting{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list workspace settings: %v", err)
	}
	workspaceSetting := &v1pb.WorkspaceSetting{}
	for _, v := range workspaceSettings {
		if v.Key == storepb.WorkspaceSettingKey_WORKSPACE_SETTING_GENERAL {
			generalSetting := v.GetGeneral()
			workspaceSetting.Branding = generalSetting.GetBranding()
			workspaceSetting.CustomStyle = generalSetting.GetCustomStyle()
		} else if v.Key == storepb.WorkspaceSettingKey_WORKSPACE_SETTING_SECURITY {
			securitySetting := v.GetSecurity()
			workspaceSetting.DisallowUserRegistration = securitySetting.GetDisallowUserRegistration()
			workspaceSetting.DisallowPasswordAuth = securitySetting.GetDisallowPasswordAuth()
		} else if v.Key == storepb.WorkspaceSettingKey_WORKSPACE_SETTING_SHORTCUT_RELATED {
			shortcutRelatedSetting := v.GetShortcutRelated()
			workspaceSetting.DefaultVisibility = convertVisibilityFromStorepb(shortcutRelatedSetting.GetDefaultVisibility())
			workspaceSetting.ShortcutPrefix = shortcutRelatedSetting.GetShortcutPrefix()
			// Set default if empty
			if workspaceSetting.ShortcutPrefix == "" {
				workspaceSetting.ShortcutPrefix = "s"
			}
		} else if v.Key == storepb.WorkspaceSettingKey_WORKSPACE_SETTING_IDENTITY_PROVIDER {
			identityProviderSetting := v.GetIdentityProvider()
			workspaceSetting.IdentityProviders = []*v1pb.IdentityProvider{}
			for _, identityProvider := range identityProviderSetting.GetIdentityProviders() {
				identityProviderV1pb := convertIdentityProviderFromStore(identityProvider)
				if currentUser == nil || currentUser.Role != store.RoleAdmin {
					oauth2Config := identityProviderV1pb.Config.GetOauth2()
					if oauth2Config != nil {
						oauth2Config.ClientSecret = ""
					}
				}
				workspaceSetting.IdentityProviders = append(workspaceSetting.IdentityProviders, identityProviderV1pb)
			}
		}
	}
	return workspaceSetting, nil
}

func (s *APIV1Service) UpdateWorkspaceSetting(ctx context.Context, request *v1pb.UpdateWorkspaceSettingRequest) (*v1pb.WorkspaceSetting, error) {
	if request.UpdateMask == nil || len(request.UpdateMask.Paths) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "update mask is empty")
	}

	for _, path := range request.UpdateMask.Paths {
		if path == "branding" {
			generalSetting, err := s.Store.GetWorkspaceGeneralSetting(ctx)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to get workspace setting: %v", err)
			}
			generalSetting.Branding = request.Setting.Branding
			if _, err := s.Store.UpsertWorkspaceSetting(ctx, &storepb.WorkspaceSetting{
				Key: storepb.WorkspaceSettingKey_WORKSPACE_SETTING_GENERAL,
				Value: &storepb.WorkspaceSetting_General{
					General: generalSetting,
				},
			}); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update workspace setting: %v", err)
			}
		} else if path == "custom_style" {
			generalSetting, err := s.Store.GetWorkspaceGeneralSetting(ctx)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to get workspace setting: %v", err)
			}
			generalSetting.CustomStyle = request.Setting.CustomStyle
			if _, err := s.Store.UpsertWorkspaceSetting(ctx, &storepb.WorkspaceSetting{
				Key: storepb.WorkspaceSettingKey_WORKSPACE_SETTING_GENERAL,
				Value: &storepb.WorkspaceSetting_General{
					General: generalSetting,
				},
			}); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update workspace setting: %v", err)
			}
		} else if path == "default_visibility" {
			shortcutRelatedSetting, err := s.Store.GetWorkspaceSetting(ctx, &store.FindWorkspaceSetting{
				Key: storepb.WorkspaceSettingKey_WORKSPACE_SETTING_SHORTCUT_RELATED,
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to get workspace setting: %v", err)
			}
			if shortcutRelatedSetting == nil {
				shortcutRelatedSetting = &storepb.WorkspaceSetting{
					Key: storepb.WorkspaceSettingKey_WORKSPACE_SETTING_SHORTCUT_RELATED,
					Value: &storepb.WorkspaceSetting_ShortcutRelated{
						ShortcutRelated: &storepb.WorkspaceSetting_ShortcutRelatedSetting{},
					},
				}
			}
			shortcutRelatedSetting.GetShortcutRelated().DefaultVisibility = convertVisibilityToStorepb(request.Setting.DefaultVisibility)
			if _, err := s.Store.UpsertWorkspaceSetting(ctx, &storepb.WorkspaceSetting{
				Key: storepb.WorkspaceSettingKey_WORKSPACE_SETTING_SHORTCUT_RELATED,
				Value: &storepb.WorkspaceSetting_ShortcutRelated{
					ShortcutRelated: shortcutRelatedSetting.GetShortcutRelated(),
				},
			}); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update workspace setting: %v", err)
			}
		} else if path == "shortcut_prefix" {
			shortcutRelatedSetting, err := s.Store.GetWorkspaceSetting(ctx, &store.FindWorkspaceSetting{
				Key: storepb.WorkspaceSettingKey_WORKSPACE_SETTING_SHORTCUT_RELATED,
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to get workspace setting: %v", err)
			}
			if shortcutRelatedSetting == nil {
				shortcutRelatedSetting = &storepb.WorkspaceSetting{
					Key: storepb.WorkspaceSettingKey_WORKSPACE_SETTING_SHORTCUT_RELATED,
					Value: &storepb.WorkspaceSetting_ShortcutRelated{
						ShortcutRelated: &storepb.WorkspaceSetting_ShortcutRelatedSetting{},
					},
				}
			}
			shortcutRelatedSetting.GetShortcutRelated().ShortcutPrefix = request.Setting.ShortcutPrefix
			if _, err := s.Store.UpsertWorkspaceSetting(ctx, &storepb.WorkspaceSetting{
				Key: storepb.WorkspaceSettingKey_WORKSPACE_SETTING_SHORTCUT_RELATED,
				Value: &storepb.WorkspaceSetting_ShortcutRelated{
					ShortcutRelated: shortcutRelatedSetting.GetShortcutRelated(),
				},
			}); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update workspace setting: %v", err)
			}
		} else if path == "identity_providers" {
			identityProviderSetting := &storepb.WorkspaceSetting_IdentityProviderSetting{}
			for _, identityProvider := range request.Setting.IdentityProviders {
				identityProviderSetting.IdentityProviders = append(identityProviderSetting.IdentityProviders, convertIdentityProviderToStore(identityProvider))
			}
			if _, err := s.Store.UpsertWorkspaceSetting(ctx, &storepb.WorkspaceSetting{
				Key: storepb.WorkspaceSettingKey_WORKSPACE_SETTING_IDENTITY_PROVIDER,
				Value: &storepb.WorkspaceSetting_IdentityProvider{
					IdentityProvider: identityProviderSetting,
				},
			}); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update workspace setting: %v", err)
			}
		} else if path == "disallow_user_registration" {
			securitySetting, err := s.Store.GetWorkspaceSecuritySetting(ctx)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to get workspace setting: %v", err)
			}
			securitySetting.DisallowUserRegistration = request.Setting.DisallowUserRegistration
			if _, err := s.Store.UpsertWorkspaceSetting(ctx, &storepb.WorkspaceSetting{
				Key: storepb.WorkspaceSettingKey_WORKSPACE_SETTING_SECURITY,
				Value: &storepb.WorkspaceSetting_Security{
					Security: securitySetting,
				},
			}); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update workspace setting: %v", err)
			}
		} else if path == "disallow_password_auth" {
			securitySetting, err := s.Store.GetWorkspaceSecuritySetting(ctx)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to get workspace setting: %v", err)
			}
			securitySetting.DisallowPasswordAuth = request.Setting.DisallowPasswordAuth
			if _, err := s.Store.UpsertWorkspaceSetting(ctx, &storepb.WorkspaceSetting{
				Key: storepb.WorkspaceSettingKey_WORKSPACE_SETTING_SECURITY,
				Value: &storepb.WorkspaceSetting_Security{
					Security: securitySetting,
				},
			}); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update workspace setting: %v", err)
			}
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "invalid path: %s", path)
		}
	}

	workspaceSetting, err := s.GetWorkspaceSetting(ctx, &v1pb.GetWorkspaceSettingRequest{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get workspace setting: %v", err)
	}
	return workspaceSetting, nil
}

var ownerCache *v1pb.User

func (s *APIV1Service) GetWorkspaceStats(ctx context.Context, _ *v1pb.GetWorkspaceStatsRequest) (*v1pb.WorkspaceStats, error) {
	// Get total shortcuts count
	shortcuts, err := s.Store.ListShortcuts(ctx, &store.FindShortcut{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list shortcuts: %v", err)
	}
	totalShortcuts := int32(len(shortcuts))

	// Get total users count
	users, err := s.Store.ListUsers(ctx, &store.FindUser{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}
	totalUsers := int32(len(users))

	// Get total collections count
	collections, err := s.Store.ListCollections(ctx, &store.FindCollection{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list collections: %v", err)
	}
	totalCollections := int32(len(collections))

	// Get total hits count (shortcut views)
	activities, err := s.Store.ListActivities(ctx, &store.FindActivity{
		Type: store.ActivityShortcutView,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list activities: %v", err)
	}
	totalHits := int32(len(activities))

	// Get historical measurements (last 100 entries, ordered by timestamp desc)
	orderByMeasuredTsDesc := true
	limit := int32(100)
	measurements, err := s.Store.ListStatsMeasurements(ctx, &store.FindStatsMeasurement{
		OrderByMeasuredTsDesc: &orderByMeasuredTsDesc,
		Limit:                 &limit,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list stats measurements: %v", err)
	}

	// Convert measurements to API format
	historicalData := make([]*v1pb.StatsMeasurement, len(measurements))
	for i, measurement := range measurements {
		historicalData[i] = &v1pb.StatsMeasurement{
			MeasuredTs:       measurement.MeasuredTs,
			ShortcutsCount:   measurement.ShortcutsCount,
			UsersCount:       measurement.UsersCount,
			CollectionsCount: measurement.CollectionsCount,
			HitsCount:        measurement.HitsCount,
		}
	}

	return &v1pb.WorkspaceStats{
		TotalShortcuts:   totalShortcuts,
		TotalUsers:       totalUsers,
		TotalCollections: totalCollections,
		TotalHits:        totalHits,
		HistoricalData:   historicalData,
	}, nil
}

func (s *APIV1Service) GetDatabaseStats(ctx context.Context, _ *v1pb.GetDatabaseStatsRequest) (*v1pb.DatabaseStats, error) {
	// Get total users count
	users, err := s.Store.ListUsers(ctx, &store.FindUser{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	// Get total shortcuts count
	shortcuts, err := s.Store.ListShortcuts(ctx, &store.FindShortcut{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list shortcuts: %v", err)
	}

	// Get total collections count
	collections, err := s.Store.ListCollections(ctx, &store.FindCollection{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list collections: %v", err)
	}

	// Get total tags count
	tags, err := s.Store.ListTags(ctx, &store.FindTag{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tags: %v", err)
	}

	// Get total bookmark tags count
	bookmarkTags, err := s.Store.ListBookmarkTags(ctx, &store.FindBookmarkTag{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list bookmark tags: %v", err)
	}

	// Get total friendships count
	friendships, err := s.Store.ListFriendships(ctx, 0, "")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list friendships: %v", err)
	}

	// Get total followings count
	followings, err := s.Store.ListFollowing(ctx, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list followings: %v", err)
	}

	// Get total activities count
	activities, err := s.Store.ListActivities(ctx, &store.FindActivity{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list activities: %v", err)
	}

	// Get total invitations count
	invitations, err := s.Store.ListInvitations(ctx, &store.FindInvitation{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list invitations: %v", err)
	}

	return &v1pb.DatabaseStats{
		Users:        int32(len(users)),
		Shortcuts:    int32(len(shortcuts)),
		Collections:  int32(len(collections)),
		Tags:         int32(len(tags)),
		BookmarkTags: int32(len(bookmarkTags)),
		Friendships:  int32(len(friendships)),
		Followings:   int32(len(followings)),
		Activities:   int32(len(activities)),
		Invitations:  int32(len(invitations)),
	}, nil
}

func (s *APIV1Service) ExportDatabase(ctx context.Context, request *v1pb.ExportDatabaseRequest) (*v1pb.ExportDatabaseResponse, error) {
	// Verify user is admin
	currentUser, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if currentUser == nil || currentUser.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "only admin can export database")
	}

	exportData := make(map[string]interface{})

	for _, entity := range request.Entities {
		switch entity {
		case "users":
			users, err := s.Store.ListUsers(ctx, &store.FindUser{})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
			}
			exportData["users"] = users

		case "shortcuts":
			shortcuts, err := s.Store.ListShortcuts(ctx, &store.FindShortcut{})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to list shortcuts: %v", err)
			}
			exportData["shortcuts"] = shortcuts

		case "collections":
			collections, err := s.Store.ListCollections(ctx, &store.FindCollection{})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to list collections: %v", err)
			}
			exportData["collections"] = collections

		case "tags":
			tags, err := s.Store.ListTags(ctx, &store.FindTag{})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to list tags: %v", err)
			}
			exportData["tags"] = tags

		case "friendships":
			friendships, err := s.Store.ListFriendships(ctx, 0, "")
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to list friendships: %v", err)
			}
			exportData["friendships"] = friendships

		case "followings":
			followings, err := s.Store.ListFollowing(ctx, 0)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to list followings: %v", err)
			}
			exportData["followings"] = followings

		default:
			return nil, status.Errorf(codes.InvalidArgument, "unknown entity type: %s", entity)
		}
	}

	// Serialize to JSON
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to serialize data: %v", err)
	}

	// Generate filename with timestamp
	filename := fmt.Sprintf("monotreme_export_%s.json", time.Now().Format("2006_01_02"))

	return &v1pb.ExportDatabaseResponse{
		Data:     string(jsonData),
		Filename: filename,
	}, nil
}

func (s *APIV1Service) ImportDatabase(ctx context.Context, request *v1pb.ImportDatabaseRequest) (*v1pb.ImportDatabaseResponse, error) {
	// Verify user is admin
	currentUser, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if currentUser == nil || currentUser.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "only admin can import database")
	}

	// Parse JSON data
	var importData map[string]interface{}
	if err := json.Unmarshal([]byte(request.Data), &importData); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse JSON data: %v", err)
	}

	importedCounts := make(map[string]int32)
	messages := []string{}

	for _, entity := range request.Entities {
		switch entity {
		case "users":
			if usersData, ok := importData["users"]; ok {
				count, err := s.importUsers(ctx, usersData, request.Mode)
				if err != nil {
					messages = append(messages, fmt.Sprintf("Failed to import users: %v", err))
				} else {
					importedCounts["users"] = count
				}
			}

		case "shortcuts":
			if shortcutsData, ok := importData["shortcuts"]; ok {
				count, err := s.importShortcuts(ctx, shortcutsData, request.Mode)
				if err != nil {
					messages = append(messages, fmt.Sprintf("Failed to import shortcuts: %v", err))
				} else {
					importedCounts["shortcuts"] = count
				}
			}

		case "collections":
			if collectionsData, ok := importData["collections"]; ok {
				count, err := s.importCollections(ctx, collectionsData, request.Mode)
				if err != nil {
					messages = append(messages, fmt.Sprintf("Failed to import collections: %v", err))
				} else {
					importedCounts["collections"] = count
				}
			}

		case "tags":
			if tagsData, ok := importData["tags"]; ok {
				count, err := s.importTags(ctx, tagsData, request.Mode)
				if err != nil {
					messages = append(messages, fmt.Sprintf("Failed to import tags: %v", err))
				} else {
					importedCounts["tags"] = count
				}
			}

		case "friendships":
			if friendshipsData, ok := importData["friendships"]; ok {
				count, err := s.importFriendships(ctx, friendshipsData, request.Mode)
				if err != nil {
					messages = append(messages, fmt.Sprintf("Failed to import friendships: %v", err))
				} else {
					importedCounts["friendships"] = count
				}
			}

		case "followings":
			if followingsData, ok := importData["followings"]; ok {
				count, err := s.importFollowings(ctx, followingsData, request.Mode)
				if err != nil {
					messages = append(messages, fmt.Sprintf("Failed to import followings: %v", err))
				} else {
					importedCounts["followings"] = count
				}
			}
		}
	}

	return &v1pb.ImportDatabaseResponse{
		ImportedCounts: importedCounts,
		Messages:       messages,
	}, nil
}

// Helper functions for importing each entity type
func (s *APIV1Service) importUsers(ctx context.Context, data interface{}, mode string) (int32, error) {
	// Convert data to JSON and back to proper struct
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal user data: %v", err)
	}

	var users []*store.User
	if err := json.Unmarshal(jsonBytes, &users); err != nil {
		return 0, fmt.Errorf("failed to unmarshal user data: %v", err)
	}

	// Handle wipe-and-import mode
	if mode == "wipe-and-import" {
		// Note: This is dangerous and should only be used with caution
		// We don't actually implement full wipe for users as it could break the system
		return 0, fmt.Errorf("wipe-and-import mode not supported for users (too dangerous)")
	}

	count := int32(0)
	for _, user := range users {
		// Check if user exists
		existingUser, err := s.Store.GetUser(ctx, &store.FindUser{
			Email: &user.Email,
		})
		if err != nil {
			return count, fmt.Errorf("failed to check existing user: %v", err)
		}

		if existingUser != nil {
			// User exists
			if mode == "new-only" {
				continue // Skip existing users
			}
			// Overwrite mode: update the user
			if _, err := s.Store.UpdateUser(ctx, &store.UpdateUser{
				ID:       existingUser.ID,
				Nickname: &user.Nickname,
				Role:     &user.Role,
			}); err != nil {
				return count, fmt.Errorf("failed to update user: %v", err)
			}
		} else {
			// User doesn't exist, create it
			if _, err := s.Store.CreateUser(ctx, user); err != nil {
				return count, fmt.Errorf("failed to create user: %v", err)
			}
		}
		count++
	}

	return count, nil
}

func (s *APIV1Service) importShortcuts(ctx context.Context, data interface{}, mode string) (int32, error) {
	// Convert data to JSON and back to proper struct
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal shortcut data: %v", err)
	}

	var shortcuts []*storepb.Shortcut
	if err := json.Unmarshal(jsonBytes, &shortcuts); err != nil {
		return 0, fmt.Errorf("failed to unmarshal shortcut data: %v", err)
	}

	// Handle wipe-and-import mode
	if mode == "wipe-and-import" {
		// Get all shortcuts and delete them
		existingShortcuts, err := s.Store.ListShortcuts(ctx, &store.FindShortcut{})
		if err != nil {
			return 0, fmt.Errorf("failed to list shortcuts for wipe: %v", err)
		}
		for _, shortcut := range existingShortcuts {
			if err := s.Store.DeleteShortcut(ctx, &store.DeleteShortcut{ID: shortcut.Id}); err != nil {
				return 0, fmt.Errorf("failed to delete shortcut during wipe: %v", err)
			}
		}
	}

	count := int32(0)
	for _, shortcut := range shortcuts {
		// Check if shortcut exists by name
		existingShortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
			Name: &shortcut.Name,
		})
		if err != nil {
			return count, fmt.Errorf("failed to check existing shortcut: %v", err)
		}

		if existingShortcut != nil {
			// Shortcut exists
			if mode == "new-only" {
				continue // Skip existing shortcuts
			}
			// Overwrite mode: update the shortcut
			if _, err := s.Store.UpdateShortcut(ctx, &store.UpdateShortcut{
				ID:          existingShortcut.Id,
				Link:        &shortcut.Link,
				Title:       &shortcut.Title,
				Description: &shortcut.Description,
				Visibility:  &shortcut.Visibility,
			}); err != nil {
				return count, fmt.Errorf("failed to update shortcut: %v", err)
			}
		} else {
			// Shortcut doesn't exist, create it
			if _, err := s.Store.CreateShortcut(ctx, shortcut); err != nil {
				return count, fmt.Errorf("failed to create shortcut: %v", err)
			}
		}
		count++
	}

	return count, nil
}

func (s *APIV1Service) importCollections(ctx context.Context, data interface{}, mode string) (int32, error) {
	// Convert data to JSON and back to proper struct
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal collection data: %v", err)
	}

	var collections []*storepb.Collection
	if err := json.Unmarshal(jsonBytes, &collections); err != nil {
		return 0, fmt.Errorf("failed to unmarshal collection data: %v", err)
	}

	// Handle wipe-and-import mode
	if mode == "wipe-and-import" {
		existingCollections, err := s.Store.ListCollections(ctx, &store.FindCollection{})
		if err != nil {
			return 0, fmt.Errorf("failed to list collections for wipe: %v", err)
		}
		for _, collection := range existingCollections {
			if err := s.Store.DeleteCollection(ctx, &store.DeleteCollection{ID: collection.Id}); err != nil {
				return 0, fmt.Errorf("failed to delete collection during wipe: %v", err)
			}
		}
	}

	count := int32(0)
	for _, collection := range collections {
		// Check if collection exists by name
		existingCollection, err := s.Store.GetCollection(ctx, &store.FindCollection{
			Name: &collection.Name,
		})
		if err != nil {
			return count, fmt.Errorf("failed to check existing collection: %v", err)
		}

		if existingCollection != nil {
			// Collection exists
			if mode == "new-only" {
				continue // Skip existing collections
			}
			// Overwrite mode: update the collection
			if _, err := s.Store.UpdateCollection(ctx, &store.UpdateCollection{
				ID:          existingCollection.Id,
				Title:       &collection.Title,
				Description: &collection.Description,
				Visibility:  &collection.Visibility,
			}); err != nil {
				return count, fmt.Errorf("failed to update collection: %v", err)
			}
		} else {
			// Collection doesn't exist, create it
			if _, err := s.Store.CreateCollection(ctx, collection); err != nil {
				return count, fmt.Errorf("failed to create collection: %v", err)
			}
		}
		count++
	}

	return count, nil
}

func (s *APIV1Service) importTags(ctx context.Context, data interface{}, mode string) (int32, error) {
	// Convert data to JSON and back to proper struct
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal tag data: %v", err)
	}

	var tags []*storepb.Tag
	if err := json.Unmarshal(jsonBytes, &tags); err != nil {
		return 0, fmt.Errorf("failed to unmarshal tag data: %v", err)
	}

	// Handle wipe-and-import mode
	if mode == "wipe-and-import" {
		existingTags, err := s.Store.ListTags(ctx, &store.FindTag{})
		if err != nil {
			return 0, fmt.Errorf("failed to list tags for wipe: %v", err)
		}
		for _, tag := range existingTags {
			if err := s.Store.DeleteTag(ctx, &store.DeleteTag{UUID: tag.Uuid}); err != nil {
				return 0, fmt.Errorf("failed to delete tag during wipe: %v", err)
			}
		}
	}

	count := int32(0)
	for _, tag := range tags {
		// Check if tag exists by UUID
		existingTag, err := s.Store.GetTag(ctx, &store.FindTag{
			UUID: &tag.Uuid,
		})
		if err != nil {
			return count, fmt.Errorf("failed to check existing tag: %v", err)
		}

		if existingTag != nil {
			// Tag exists
			if mode == "new-only" {
				continue // Skip existing tags
			}
			// Overwrite mode: update the tag
			if _, err := s.Store.UpdateTag(ctx, &store.UpdateTag{
				UUID: existingTag.Uuid,
				Name: &tag.Name,
			}); err != nil {
				return count, fmt.Errorf("failed to update tag: %v", err)
			}
		} else {
			// Tag doesn't exist, create it
			if _, err := s.Store.CreateTag(ctx, tag); err != nil {
				return count, fmt.Errorf("failed to create tag: %v", err)
			}
		}
		count++
	}

	return count, nil
}

func (s *APIV1Service) importFriendships(ctx context.Context, data interface{}, mode string) (int32, error) {
	// Convert data to JSON and back to proper struct
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal friendship data: %v", err)
	}

	var friendships []*store.Friendship
	if err := json.Unmarshal(jsonBytes, &friendships); err != nil {
		return 0, fmt.Errorf("failed to unmarshal friendship data: %v", err)
	}

	// Handle wipe-and-import mode
	if mode == "wipe-and-import" {
		existingFriendships, err := s.Store.ListFriendships(ctx, 0, "")
		if err != nil {
			return 0, fmt.Errorf("failed to list friendships for wipe: %v", err)
		}
		for _, friendship := range existingFriendships {
			if err := s.Store.DeleteFriendship(ctx, friendship.ID); err != nil {
				return 0, fmt.Errorf("failed to delete friendship during wipe: %v", err)
			}
		}
	}

	count := int32(0)
	for _, friendship := range friendships {
		// For friendships, we check if a relationship already exists between these users
		existingFriendship, err := s.Store.GetFriendship(ctx, friendship.ID)
		if err != nil {
			return count, fmt.Errorf("failed to check existing friendship: %v", err)
		}

		if existingFriendship != nil {
			// Friendship exists
			if mode == "new-only" {
				continue // Skip existing friendships
			}
			// For overwrite mode with friendships, we could update status
			// But since we don't have an update method, we'll skip this
			continue
		} else {
			// Friendship doesn't exist, create it
			if _, err := s.Store.CreateFriendship(ctx, friendship.UserID, friendship.FriendID); err != nil {
				return count, fmt.Errorf("failed to create friendship: %v", err)
			}
		}
		count++
	}

	return count, nil
}

func (s *APIV1Service) importFollowings(ctx context.Context, data interface{}, mode string) (int32, error) {
	// Convert data to JSON and back to proper struct
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal following data: %v", err)
	}

	var followings []*store.Following
	if err := json.Unmarshal(jsonBytes, &followings); err != nil {
		return 0, fmt.Errorf("failed to unmarshal following data: %v", err)
	}

	// Handle wipe-and-import mode
	if mode == "wipe-and-import" {
		existingFollowings, err := s.Store.ListFollowing(ctx, 0)
		if err != nil {
			return 0, fmt.Errorf("failed to list followings for wipe: %v", err)
		}
		for _, following := range existingFollowings {
			if err := s.Store.DeleteFollowing(ctx, following.ID); err != nil {
				return 0, fmt.Errorf("failed to delete following during wipe: %v", err)
			}
		}
	}

	count := int32(0)
	for _, following := range followings {
		// Check if following relationship exists
		existingFollowing, err := s.Store.GetFollowing(ctx, following.ID)
		if err != nil {
			return count, fmt.Errorf("failed to check existing following: %v", err)
		}

		if existingFollowing != nil {
			// Following exists
			if mode == "new-only" {
				continue // Skip existing followings
			}
			// For overwrite mode with followings, there's nothing to update
			continue
		} else {
			// Following doesn't exist, create it
			if _, err := s.Store.CreateFollowing(ctx, following.FollowerID, following.FollowingID); err != nil {
				return count, fmt.Errorf("failed to create following: %v", err)
			}
		}
		count++
	}

	return count, nil
}

func (s *APIV1Service) GetInstanceOwner(ctx context.Context) (*v1pb.User, error) {
	if ownerCache != nil {
		return ownerCache, nil
	}

	adminRole := store.RoleAdmin
	user, err := s.Store.GetUser(ctx, &store.FindUser{
		Role: &adminRole,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find admin")
	}
	if user == nil {
		return nil, nil
	}

	ownerCache = convertUserFromStore(user)
	return ownerCache, nil
}

func convertIdentityProviderFromStore(identityProvider *storepb.IdentityProvider) *v1pb.IdentityProvider {
	if identityProvider == nil {
		return nil
	}
	return &v1pb.IdentityProvider{
		Id:     identityProvider.Id,
		Title:  identityProvider.Title,
		Type:   v1pb.IdentityProvider_Type(identityProvider.Type),
		Config: convertIdentityProviderConfigFromStore(identityProvider.Config),
	}
}

func convertIdentityProviderConfigFromStore(identityProviderConfig *storepb.IdentityProviderConfig) *v1pb.IdentityProviderConfig {
	oauth2Config := identityProviderConfig.GetOauth2()
	if oauth2Config != nil {
		return &v1pb.IdentityProviderConfig{
			Config: &v1pb.IdentityProviderConfig_Oauth2{
				Oauth2: &v1pb.IdentityProviderConfig_OAuth2Config{
					ClientId:     oauth2Config.ClientId,
					ClientSecret: oauth2Config.ClientSecret,
					AuthUrl:      oauth2Config.AuthUrl,
					TokenUrl:     oauth2Config.TokenUrl,
					UserInfoUrl:  oauth2Config.UserInfoUrl,
					Scopes:       oauth2Config.Scopes,
					FieldMapping: &v1pb.IdentityProviderConfig_FieldMapping{
						Identifier:  oauth2Config.FieldMapping.Identifier,
						DisplayName: oauth2Config.FieldMapping.DisplayName,
					},
				},
			},
		}
	}
	return nil
}

func convertIdentityProviderToStore(identityProvider *v1pb.IdentityProvider) *storepb.IdentityProvider {
	if identityProvider == nil {
		return nil
	}
	return &storepb.IdentityProvider{
		Id:     identityProvider.Id,
		Title:  identityProvider.Title,
		Type:   storepb.IdentityProvider_Type(identityProvider.Type),
		Config: convertIdentityProviderConfigToStore(identityProvider.Config),
	}
}

func convertIdentityProviderConfigToStore(identityProviderConfig *v1pb.IdentityProviderConfig) *storepb.IdentityProviderConfig {
	oauth2Config := identityProviderConfig.GetOauth2()
	if oauth2Config != nil {
		return &storepb.IdentityProviderConfig{
			Config: &storepb.IdentityProviderConfig_Oauth2{
				Oauth2: &storepb.IdentityProviderConfig_OAuth2Config{
					ClientId:     oauth2Config.ClientId,
					ClientSecret: oauth2Config.ClientSecret,
					AuthUrl:      oauth2Config.AuthUrl,
					TokenUrl:     oauth2Config.TokenUrl,
					UserInfoUrl:  oauth2Config.UserInfoUrl,
					Scopes:       oauth2Config.Scopes,
					FieldMapping: &storepb.IdentityProviderConfig_FieldMapping{
						Identifier:  oauth2Config.FieldMapping.Identifier,
						DisplayName: oauth2Config.FieldMapping.DisplayName,
					},
				},
			},
		}
	}
	return nil
}
