package store

import (
	"context"
	"database/sql"

	storepb "github.com/bshort/monotreme/proto/gen/store"
)

// Driver is an interface for store driver.
// It contains all methods that store database driver should implement.
type Driver interface {
	GetDB() *sql.DB
	Close() error

	// MigrationHistory model related methods.
	UpsertMigrationHistory(ctx context.Context, upsert *UpsertMigrationHistory) (*MigrationHistory, error)
	ListMigrationHistories(ctx context.Context, find *FindMigrationHistory) ([]*MigrationHistory, error)

	// Activity model related methods.
	CreateActivity(ctx context.Context, create *Activity) (*Activity, error)
	ListActivities(ctx context.Context, find *FindActivity) ([]*Activity, error)

	// Collection model related methods.
	CreateCollection(ctx context.Context, create *storepb.Collection) (*storepb.Collection, error)
	UpdateCollection(ctx context.Context, update *UpdateCollection) (*storepb.Collection, error)
	ListCollections(ctx context.Context, find *FindCollection) ([]*storepb.Collection, error)
	DeleteCollection(ctx context.Context, delete *DeleteCollection) error

	// Shortcut model related methods.
	CreateShortcut(ctx context.Context, create *storepb.Shortcut) (*storepb.Shortcut, error)
	UpdateShortcut(ctx context.Context, update *UpdateShortcut) (*storepb.Shortcut, error)
	ListShortcuts(ctx context.Context, find *FindShortcut) ([]*storepb.Shortcut, error)
	CountShortcuts(ctx context.Context, find *FindShortcut) (int32, error)
	DeleteShortcut(ctx context.Context, delete *DeleteShortcut) error

	// User model related methods.
	CreateUser(ctx context.Context, create *User) (*User, error)
	UpdateUser(ctx context.Context, update *UpdateUser) (*User, error)
	ListUsers(ctx context.Context, find *FindUser) ([]*User, error)
	DeleteUser(ctx context.Context, delete *DeleteUser) error

	// UserSetting model related methods.
	UpsertUserSetting(ctx context.Context, upsert *storepb.UserSetting) (*storepb.UserSetting, error)
	ListUserSettings(ctx context.Context, find *FindUserSetting) ([]*storepb.UserSetting, error)

	// WorkspaceSetting model related methods.
	UpsertWorkspaceSetting(ctx context.Context, upsert *storepb.WorkspaceSetting) (*storepb.WorkspaceSetting, error)
	ListWorkspaceSettings(ctx context.Context, find *FindWorkspaceSetting) ([]*storepb.WorkspaceSetting, error)
	DeleteWorkspaceSetting(ctx context.Context, key storepb.WorkspaceSettingKey) error

	// StatsMeasurement model related methods.
	CreateStatsMeasurement(ctx context.Context, create *StatsMeasurement) (*StatsMeasurement, error)
	ListStatsMeasurements(ctx context.Context, find *FindStatsMeasurement) ([]*StatsMeasurement, error)
	UpdateStatsMeasurement(ctx context.Context, update *UpdateStatsMeasurement) (*StatsMeasurement, error)
	DeleteStatsMeasurement(ctx context.Context, delete *DeleteStatsMeasurement) error

	// RssFeed model related methods.
	CreateRssFeed(ctx context.Context, create *RssFeed) (*RssFeed, error)
	UpdateRssFeed(ctx context.Context, update *UpdateRssFeed) (*RssFeed, error)
	ListRssFeeds(ctx context.Context, find *FindRssFeed) ([]*RssFeed, error)
	DeleteRssFeed(ctx context.Context, delete *DeleteRssFeed) error

	// RssFeedItem model related methods.
	CreateRssFeedItem(ctx context.Context, create *RssFeedItem) (*RssFeedItem, error)
	GetRssFeedItem(ctx context.Context, find *FindRssFeedItem) (*RssFeedItem, error)
	ListRssFeedItems(ctx context.Context, find *FindRssFeedItem) ([]*RssFeedItem, error)

	// Invitation model related methods.
	CreateInvitation(ctx context.Context, create *storepb.Invitation) (*storepb.Invitation, error)
	UpdateInvitation(ctx context.Context, update *UpdateInvitation) (*storepb.Invitation, error)
	ListInvitations(ctx context.Context, find *FindInvitation) ([]*storepb.Invitation, error)
	DeleteInvitation(ctx context.Context, delete *DeleteInvitation) error

	// Tag model related methods.
	CreateTag(ctx context.Context, create *storepb.Tag) (*storepb.Tag, error)
	UpdateTag(ctx context.Context, update *UpdateTag) (*storepb.Tag, error)
	ListTags(ctx context.Context, find *FindTag) ([]*storepb.Tag, error)
	DeleteTag(ctx context.Context, delete *DeleteTag) error

	// BookmarkTag model related methods.
	CreateBookmarkTag(ctx context.Context, create *CreateBookmarkTag) (*storepb.BookmarkTag, error)
	ListBookmarkTags(ctx context.Context, find *FindBookmarkTag) ([]*storepb.BookmarkTag, error)
	DeleteBookmarkTag(ctx context.Context, delete *DeleteBookmarkTag) error
	DeleteBookmarkTagsByShortcut(ctx context.Context, shortcutID int32) error
	DeleteBookmarkTagsByTag(ctx context.Context, tagUUID string) error

	// Friendship model related methods.
	ListFriendships(ctx context.Context, userID int32, status string) ([]*Friendship, error)
	ListIncomingFriendRequests(ctx context.Context, userID int32) ([]*Friendship, error)
	ListOutgoingFriendRequests(ctx context.Context, userID int32) ([]*Friendship, error)
	GetFriendship(ctx context.Context, id int32) (*Friendship, error)
	CreateFriendship(ctx context.Context, userID int32, friendID int32) (*Friendship, error)
	AcceptFriendship(ctx context.Context, id int32) error
	DeleteFriendship(ctx context.Context, id int32) error

	// Following model related methods.
	ListFollowing(ctx context.Context, userID int32) ([]*Following, error)
	ListFollowers(ctx context.Context, userID int32) ([]*Following, error)
	GetFollowing(ctx context.Context, id int32) (*Following, error)
	CreateFollowing(ctx context.Context, followerID int32, followingID int32) (*Following, error)
	DeleteFollowing(ctx context.Context, id int32) error
	GetFollowingUserShortcuts(ctx context.Context, followerID int32, followingID int32) ([]*storepb.Shortcut, error)
}
