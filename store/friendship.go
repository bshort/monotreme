package store

import (
	"context"
	"time"
)

// Timestamp is a custom type for Unix timestamps stored as int64.
type Timestamp int64

// AsTime converts Timestamp to time.Time.
func (t Timestamp) AsTime() time.Time {
	return time.Unix(int64(t), 0)
}

// Friendship represents a friendship relationship between two users.
type Friendship struct {
	ID         int32
	UserID     int32
	FriendID   int32
	Status     string
	CreatedTs  Timestamp
	AcceptedTs *Timestamp
}

func (s *Store) ListFriendships(ctx context.Context, userID int32, status string) ([]*Friendship, error) {
	return s.driver.ListFriendships(ctx, userID, status)
}

func (s *Store) ListIncomingFriendRequests(ctx context.Context, userID int32) ([]*Friendship, error) {
	return s.driver.ListIncomingFriendRequests(ctx, userID)
}

func (s *Store) ListOutgoingFriendRequests(ctx context.Context, userID int32) ([]*Friendship, error) {
	return s.driver.ListOutgoingFriendRequests(ctx, userID)
}

func (s *Store) GetFriendship(ctx context.Context, id int32) (*Friendship, error) {
	return s.driver.GetFriendship(ctx, id)
}

func (s *Store) CreateFriendship(ctx context.Context, userID int32, friendID int32) (*Friendship, error) {
	return s.driver.CreateFriendship(ctx, userID, friendID)
}

func (s *Store) AcceptFriendship(ctx context.Context, id int32) error {
	return s.driver.AcceptFriendship(ctx, id)
}

func (s *Store) DeleteFriendship(ctx context.Context, id int32) error {
	return s.driver.DeleteFriendship(ctx, id)
}
