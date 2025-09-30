package store

import (
	"context"
	"time"

	storepb "github.com/bshort/monotreme/proto/gen/store"
)

type UpdateInvitation struct {
	ID         int32
	From       *string
	To         *string
	AcceptedAt *string
	DeletedAt  *string
}

type FindInvitation struct {
	ID       *int32
	From     *string
	To       *string
	Accepted *bool
	Deleted  *bool
}

type DeleteInvitation struct {
	ID int32
}

func (s *Store) CreateInvitation(ctx context.Context, create *storepb.Invitation) (*storepb.Invitation, error) {
	return s.driver.CreateInvitation(ctx, create)
}

func (s *Store) UpdateInvitation(ctx context.Context, update *UpdateInvitation) (*storepb.Invitation, error) {
	return s.driver.UpdateInvitation(ctx, update)
}

func (s *Store) ListInvitations(ctx context.Context, find *FindInvitation) ([]*storepb.Invitation, error) {
	return s.driver.ListInvitations(ctx, find)
}

func (s *Store) GetInvitation(ctx context.Context, find *FindInvitation) (*storepb.Invitation, error) {
	invitations, err := s.ListInvitations(ctx, find)
	if err != nil {
		return nil, err
	}

	if len(invitations) == 0 {
		return nil, nil
	}

	return invitations[0], nil
}

func (s *Store) DeleteInvitation(ctx context.Context, delete *DeleteInvitation) error {
	return s.driver.DeleteInvitation(ctx, delete)
}

func (s *Store) AcceptInvitation(ctx context.Context, invitationID int32) (*storepb.Invitation, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	return s.UpdateInvitation(ctx, &UpdateInvitation{
		ID:         invitationID,
		AcceptedAt: &now,
	})
}

func (s *Store) SoftDeleteInvitation(ctx context.Context, invitationID int32) (*storepb.Invitation, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	return s.UpdateInvitation(ctx, &UpdateInvitation{
		ID:        invitationID,
		DeletedAt: &now,
	})
}

// Helper function to generate user path from UUID
func GenerateUserPath(userUUID string) string {
	return "/users/" + userUUID
}

// Helper function to extract UUID from user path
func ExtractUUIDFromUserPath(path string) string {
	if len(path) > 7 && path[:7] == "/users/" {
		return path[7:]
	}
	return ""
}