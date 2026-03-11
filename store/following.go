package store

import (
	"context"

	storepb "github.com/bshort/monotreme/proto/gen/store"
)

// Following represents a following relationship between two users.
type Following struct {
	ID          int32
	FollowerID  int32
	FollowingID int32
	CreatedTs   Timestamp
}

func (s *Store) ListFollowing(ctx context.Context, userID int32) ([]*Following, error) {
	return s.driver.ListFollowing(ctx, userID)
}

func (s *Store) ListFollowers(ctx context.Context, userID int32) ([]*Following, error) {
	return s.driver.ListFollowers(ctx, userID)
}

func (s *Store) GetFollowing(ctx context.Context, id int32) (*Following, error) {
	return s.driver.GetFollowing(ctx, id)
}

func (s *Store) CreateFollowing(ctx context.Context, followerID int32, followingID int32) (*Following, error) {
	return s.driver.CreateFollowing(ctx, followerID, followingID)
}

func (s *Store) DeleteFollowing(ctx context.Context, id int32) error {
	return s.driver.DeleteFollowing(ctx, id)
}

func (s *Store) GetFollowingUserShortcuts(ctx context.Context, followerID int32, followingID int32) ([]*storepb.Shortcut, error) {
	return s.driver.GetFollowingUserShortcuts(ctx, followerID, followingID)
}
