package postgres

import (
	"context"
	"fmt"

	"github.com/bshort/monotreme/store"
)

func (d *DB) ListFriendships(ctx context.Context, userID int32, status string) ([]*store.Friendship, error) {
	return nil, fmt.Errorf("ListFriendships not implemented for postgres")
}

func (d *DB) ListIncomingFriendRequests(ctx context.Context, userID int32) ([]*store.Friendship, error) {
	return nil, fmt.Errorf("ListIncomingFriendRequests not implemented for postgres")
}

func (d *DB) ListOutgoingFriendRequests(ctx context.Context, userID int32) ([]*store.Friendship, error) {
	return nil, fmt.Errorf("ListOutgoingFriendRequests not implemented for postgres")
}

func (d *DB) GetFriendship(ctx context.Context, id int32) (*store.Friendship, error) {
	return nil, fmt.Errorf("GetFriendship not implemented for postgres")
}

func (d *DB) CreateFriendship(ctx context.Context, userID int32, friendID int32) (*store.Friendship, error) {
	return nil, fmt.Errorf("CreateFriendship not implemented for postgres")
}

func (d *DB) AcceptFriendship(ctx context.Context, id int32) error {
	return fmt.Errorf("AcceptFriendship not implemented for postgres")
}

func (d *DB) DeleteFriendship(ctx context.Context, id int32) error {
	return fmt.Errorf("DeleteFriendship not implemented for postgres")
}
