package postgres

import (
	"context"
	"fmt"

	"github.com/bshort/monotreme/store"
	storepb "github.com/bshort/monotreme/proto/gen/store"
)

func (d *DB) ListFollowing(ctx context.Context, userID int32) ([]*store.Following, error) {
	return nil, fmt.Errorf("ListFollowing not implemented for postgres")
}

func (d *DB) ListFollowers(ctx context.Context, userID int32) ([]*store.Following, error) {
	return nil, fmt.Errorf("ListFollowers not implemented for postgres")
}

func (d *DB) GetFollowing(ctx context.Context, id int32) (*store.Following, error) {
	return nil, fmt.Errorf("GetFollowing not implemented for postgres")
}

func (d *DB) CreateFollowing(ctx context.Context, followerID int32, followingID int32) (*store.Following, error) {
	return nil, fmt.Errorf("CreateFollowing not implemented for postgres")
}

func (d *DB) DeleteFollowing(ctx context.Context, id int32) error {
	return fmt.Errorf("DeleteFollowing not implemented for postgres")
}

func (d *DB) GetFollowingUserShortcuts(ctx context.Context, followerID int32, followingID int32) ([]*storepb.Shortcut, error) {
	return nil, fmt.Errorf("GetFollowingUserShortcuts not implemented for postgres")
}
