package teststore

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

func TestInvitationStore(t *testing.T) {
	ctx := context.Background()
	ts := NewTestingStore(ctx, t)

	// Create test users
	user1, err := ts.CreateUser(ctx, &store.User{
		Role:         store.RoleAdmin,
		Email:        "test1@test.com",
		Nickname:     "test_user_1",
		PasswordHash: "test-hash",
		UUID:         "user-uuid-1",
	})
	require.NoError(t, err)

	user2, err := ts.CreateUser(ctx, &store.User{
		Role:         store.RoleUser,
		Email:        "test2@test.com",
		Nickname:     "test_user_2",
		PasswordHash: "test-hash",
		UUID:         "user-uuid-2",
	})
	require.NoError(t, err)

	t.Run("CreateInvitation", func(t *testing.T) {
		invitation, err := ts.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(user1.UUID),
			To:   store.GenerateUserPath(user2.UUID),
		})
		require.NoError(t, err)
		require.NotZero(t, invitation.Id)
		require.Equal(t, store.GenerateUserPath(user1.UUID), invitation.From)
		require.Equal(t, store.GenerateUserPath(user2.UUID), invitation.To)
		require.NotZero(t, invitation.CreatedTs)
		require.NotZero(t, invitation.UpdatedTs)
		require.Empty(t, invitation.AcceptedAt)
		require.Empty(t, invitation.DeletedAt)
	})

	t.Run("ListInvitations", func(t *testing.T) {
		// Create a fresh test store to avoid conflicts
		listTS := NewTestingStore(ctx, t)

		// Create fresh test users
		listUser1, err := listTS.CreateUser(ctx, &store.User{
			Role:         store.RoleAdmin,
			Email:        "list1@test.com",
			Nickname:     "list_user_1",
			PasswordHash: "test-hash",
			UUID:         "list-user-uuid-1",
		})
		require.NoError(t, err)

		listUser2, err := listTS.CreateUser(ctx, &store.User{
			Role:         store.RoleUser,
			Email:        "list2@test.com",
			Nickname:     "list_user_2",
			PasswordHash: "test-hash",
			UUID:         "list-user-uuid-2",
		})
		require.NoError(t, err)

		// Create a test invitation
		invitation, err := listTS.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(listUser1.UUID),
			To:   store.GenerateUserPath(listUser2.UUID),
		})
		require.NoError(t, err)

		// List all invitations
		invitations, err := listTS.ListInvitations(ctx, &store.FindInvitation{})
		require.NoError(t, err)
		require.Len(t, invitations, 1)

		// List invitations by from user
		invitations, err = listTS.ListInvitations(ctx, &store.FindInvitation{
			From: &invitation.From,
		})
		require.NoError(t, err)
		require.Len(t, invitations, 1)
		require.Equal(t, invitation.Id, invitations[0].Id)

		// List invitations by to user
		invitations, err = listTS.ListInvitations(ctx, &store.FindInvitation{
			To: &invitation.To,
		})
		require.NoError(t, err)
		require.Len(t, invitations, 1)
		require.Equal(t, invitation.Id, invitations[0].Id)
	})

	t.Run("GetInvitation", func(t *testing.T) {
		// Create a test invitation
		created, err := ts.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(user1.UUID),
			To:   store.GenerateUserPath(user2.UUID),
		})
		require.NoError(t, err)

		// Get invitation by ID
		invitation, err := ts.GetInvitation(ctx, &store.FindInvitation{
			ID: &created.Id,
		})
		require.NoError(t, err)
		require.Equal(t, created.Id, invitation.Id)
		require.Equal(t, created.From, invitation.From)
		require.Equal(t, created.To, invitation.To)
	})

	t.Run("UpdateInvitation", func(t *testing.T) {
		// Create a test invitation
		invitation, err := ts.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(user1.UUID),
			To:   store.GenerateUserPath(user2.UUID),
		})
		require.NoError(t, err)

		// Accept the invitation
		acceptedTime := time.Now().Format(time.RFC3339)
		updated, err := ts.UpdateInvitation(ctx, &store.UpdateInvitation{
			ID:         invitation.Id,
			AcceptedAt: &acceptedTime,
		})
		require.NoError(t, err)
		require.Equal(t, acceptedTime, updated.AcceptedAt)

		// Verify the update persisted
		retrieved, err := ts.GetInvitation(ctx, &store.FindInvitation{
			ID: &invitation.Id,
		})
		require.NoError(t, err)
		require.Equal(t, acceptedTime, retrieved.AcceptedAt)
	})

	t.Run("DeleteInvitation", func(t *testing.T) {
		// Create a test invitation
		invitation, err := ts.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(user1.UUID),
			To:   store.GenerateUserPath(user2.UUID),
		})
		require.NoError(t, err)

		// Soft delete the invitation
		deletedTime := time.Now().Format(time.RFC3339)
		updated, err := ts.UpdateInvitation(ctx, &store.UpdateInvitation{
			ID:        invitation.Id,
			DeletedAt: &deletedTime,
		})
		require.NoError(t, err)
		require.Equal(t, deletedTime, updated.DeletedAt)

		// Verify it's still retrievable but marked as deleted
		retrieved, err := ts.GetInvitation(ctx, &store.FindInvitation{
			ID: &invitation.Id,
		})
		require.NoError(t, err)
		require.Equal(t, deletedTime, retrieved.DeletedAt)
	})

	t.Run("FilterByAcceptedStatus", func(t *testing.T) {
		// Create a fresh test store to avoid conflicts
		filterTS := NewTestingStore(ctx, t)

		// Create fresh test users
		filterUser1, err := filterTS.CreateUser(ctx, &store.User{
			Role:         store.RoleAdmin,
			Email:        "filter1@test.com",
			Nickname:     "filter_user_1",
			PasswordHash: "test-hash",
			UUID:         "filter-user-uuid-1",
		})
		require.NoError(t, err)

		filterUser2, err := filterTS.CreateUser(ctx, &store.User{
			Role:         store.RoleUser,
			Email:        "filter2@test.com",
			Nickname:     "filter_user_2",
			PasswordHash: "test-hash",
			UUID:         "filter-user-uuid-2",
		})
		require.NoError(t, err)

		// Create accepted invitation
		acceptedInvite, err := filterTS.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(filterUser1.UUID),
			To:   store.GenerateUserPath(filterUser2.UUID),
		})
		require.NoError(t, err)

		acceptedTime := time.Now().Format(time.RFC3339)
		_, err = filterTS.UpdateInvitation(ctx, &store.UpdateInvitation{
			ID:         acceptedInvite.Id,
			AcceptedAt: &acceptedTime,
		})
		require.NoError(t, err)

		// Create pending invitation
		pendingInvite, err := filterTS.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(filterUser2.UUID),
			To:   store.GenerateUserPath(filterUser1.UUID),
		})
		require.NoError(t, err)

		// Filter for accepted invitations
		acceptedTrue := true
		accepted, err := filterTS.ListInvitations(ctx, &store.FindInvitation{
			Accepted: &acceptedTrue,
		})
		require.NoError(t, err)
		require.Len(t, accepted, 1)
		require.Equal(t, acceptedInvite.Id, accepted[0].Id)
		require.NotEmpty(t, accepted[0].AcceptedAt)

		// Filter for pending invitations
		acceptedFalse := false
		pending, err := filterTS.ListInvitations(ctx, &store.FindInvitation{
			Accepted: &acceptedFalse,
		})
		require.NoError(t, err)
		require.Len(t, pending, 1)
		require.Equal(t, pendingInvite.Id, pending[0].Id)
		require.Empty(t, pending[0].AcceptedAt)
	})

	t.Run("FilterByDeletedStatus", func(t *testing.T) {
		// Create a fresh test store to avoid conflicts
		deletedTS := NewTestingStore(ctx, t)

		// Create fresh test users
		deletedUser1, err := deletedTS.CreateUser(ctx, &store.User{
			Role:         store.RoleAdmin,
			Email:        "deleted1@test.com",
			Nickname:     "deleted_user_1",
			PasswordHash: "test-hash",
			UUID:         "deleted-user-uuid-1",
		})
		require.NoError(t, err)

		deletedUser2, err := deletedTS.CreateUser(ctx, &store.User{
			Role:         store.RoleUser,
			Email:        "deleted2@test.com",
			Nickname:     "deleted_user_2",
			PasswordHash: "test-hash",
			UUID:         "deleted-user-uuid-2",
		})
		require.NoError(t, err)

		// Create active invitation
		activeInvite, err := deletedTS.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(deletedUser1.UUID),
			To:   store.GenerateUserPath(deletedUser2.UUID),
		})
		require.NoError(t, err)

		// Create deleted invitation
		deletedInvite, err := deletedTS.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(deletedUser2.UUID),
			To:   store.GenerateUserPath(deletedUser1.UUID),
		})
		require.NoError(t, err)

		deletedTime := time.Now().Format(time.RFC3339)
		_, err = deletedTS.UpdateInvitation(ctx, &store.UpdateInvitation{
			ID:        deletedInvite.Id,
			DeletedAt: &deletedTime,
		})
		require.NoError(t, err)

		// Filter for active invitations
		deletedFalse := false
		active, err := deletedTS.ListInvitations(ctx, &store.FindInvitation{
			Deleted: &deletedFalse,
		})
		require.NoError(t, err)
		require.Len(t, active, 1)
		require.Equal(t, activeInvite.Id, active[0].Id)
		require.Empty(t, active[0].DeletedAt)

		// Filter for deleted invitations
		deletedTrue := true
		deleted, err := deletedTS.ListInvitations(ctx, &store.FindInvitation{
			Deleted: &deletedTrue,
		})
		require.NoError(t, err)
		require.Len(t, deleted, 1)
		require.Equal(t, deletedInvite.Id, deleted[0].Id)
		require.NotEmpty(t, deleted[0].DeletedAt)
	})
}

func TestInvitationHelperFunctions(t *testing.T) {
	t.Run("GenerateUserPath", func(t *testing.T) {
		uuid := "123e4567-e89b-12d3-a456-426614174000"
		expected := "/users/123e4567-e89b-12d3-a456-426614174000"
		result := store.GenerateUserPath(uuid)
		require.Equal(t, expected, result)
	})

	t.Run("ExtractUUIDFromUserPath", func(t *testing.T) {
		path := "/users/123e4567-e89b-12d3-a456-426614174000"
		expected := "123e4567-e89b-12d3-a456-426614174000"
		result := store.ExtractUUIDFromUserPath(path)
		require.Equal(t, expected, result)
	})

	t.Run("ExtractUUIDFromUserPath_InvalidPath", func(t *testing.T) {
		invalidPaths := []string{
			"",
			"/user/123e4567-e89b-12d3-a456-426614174000",  // missing 's'
			"/users/",                                       // no UUID
			"users/123e4567-e89b-12d3-a456-426614174000",   // no leading slash
			"/groups/123e4567-e89b-12d3-a456-426614174000", // wrong prefix
		}

		for _, path := range invalidPaths {
			result := store.ExtractUUIDFromUserPath(path)
			require.Empty(t, result, "Expected empty result for path: %s", path)
		}
	})
}