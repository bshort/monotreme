package v1

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	storepb "github.com/bshort/monotreme/proto/gen/store"
	v1pb "github.com/bshort/monotreme/proto/gen/api/v1"
	"github.com/bshort/monotreme/store"
	"github.com/bshort/monotreme/store/test"
)

func TestInvitationService_ListInvitations(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupInvitationTest(t)

	// Create test invitations
	invitation1, err := s.Store.CreateInvitation(ctx, &storepb.Invitation{
		From: store.GenerateUserPath(user1.UUID),
		To:   store.GenerateUserPath(user2.UUID),
	})
	require.NoError(t, err)

	invitation2, err := s.Store.CreateInvitation(ctx, &storepb.Invitation{
		From: store.GenerateUserPath(user2.UUID),
		To:   store.GenerateUserPath(user1.UUID),
	})
	require.NoError(t, err)

	t.Run("ListAll", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.ListInvitations(ctx, &v1pb.ListInvitationsRequest{})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(resp.Invitations), 2)

		// Check that user1 can see invitations they're involved in
		foundInvitation1 := false
		foundInvitation2 := false
		for _, inv := range resp.Invitations {
			if inv.Id == invitation1.Id {
				foundInvitation1 = true
				require.Equal(t, store.GenerateUserPath(user1.UUID), inv.From)
				require.Equal(t, store.GenerateUserPath(user2.UUID), inv.To)
			}
			if inv.Id == invitation2.Id {
				foundInvitation2 = true
				require.Equal(t, store.GenerateUserPath(user2.UUID), inv.From)
				require.Equal(t, store.GenerateUserPath(user1.UUID), inv.To)
			}
		}
		require.True(t, foundInvitation1)
		require.True(t, foundInvitation2)
	})

	t.Run("FilterByFromUser", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.ListInvitations(ctx, &v1pb.ListInvitationsRequest{
			From: store.GenerateUserPath(user1.UUID),
		})
		require.NoError(t, err)

		// Should find at least invitation1
		found := false
		for _, inv := range resp.Invitations {
			require.Equal(t, store.GenerateUserPath(user1.UUID), inv.From)
			if inv.Id == invitation1.Id {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("FilterByToUser", func(t *testing.T) {
		ctx := setUserContext(ctx, user2)
		resp, err := s.ListInvitations(ctx, &v1pb.ListInvitationsRequest{
			To: store.GenerateUserPath(user2.UUID),
		})
		require.NoError(t, err)

		// Should find at least invitation1
		found := false
		for _, inv := range resp.Invitations {
			require.Equal(t, store.GenerateUserPath(user2.UUID), inv.To)
			if inv.Id == invitation1.Id {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.ListInvitations(ctx, &v1pb.ListInvitationsRequest{})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestInvitationService_GetInvitation(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupInvitationTest(t)

	// Create test invitation
	invitation, err := s.Store.CreateInvitation(ctx, &storepb.Invitation{
		From: store.GenerateUserPath(user1.UUID),
		To:   store.GenerateUserPath(user2.UUID),
	})
	require.NoError(t, err)

	t.Run("Success_FromUser", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.GetInvitation(ctx, &v1pb.GetInvitationRequest{
			Id: invitation.Id,
		})
		require.NoError(t, err)
		require.Equal(t, invitation.Id, resp.Id)
		require.Equal(t, store.GenerateUserPath(user1.UUID), resp.From)
		require.Equal(t, store.GenerateUserPath(user2.UUID), resp.To)
	})

	t.Run("Success_ToUser", func(t *testing.T) {
		ctx := setUserContext(ctx, user2)
		resp, err := s.GetInvitation(ctx, &v1pb.GetInvitationRequest{
			Id: invitation.Id,
		})
		require.NoError(t, err)
		require.Equal(t, invitation.Id, resp.Id)
		require.Equal(t, store.GenerateUserPath(user1.UUID), resp.From)
		require.Equal(t, store.GenerateUserPath(user2.UUID), resp.To)
	})

	t.Run("NotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.GetInvitation(ctx, &v1pb.GetInvitationRequest{
			Id: 99999,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		// Create a third user who is not involved in the invitation
		thirdUser, err := s.Store.CreateUser(ctx, &store.User{
			Role:         store.RoleUser,
			Email:        "test3@test.com",
			Nickname:     "test_user_3",
			PasswordHash: "test-hash",
			UUID:         "user-uuid-3",
		})
		require.NoError(t, err)

		ctx := setUserContext(ctx, thirdUser)
		_, err = s.GetInvitation(ctx, &v1pb.GetInvitationRequest{
			Id: invitation.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.GetInvitation(ctx, &v1pb.GetInvitationRequest{
			Id: invitation.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestInvitationService_CreateInvitation(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupInvitationTest(t)

	t.Run("Success", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.CreateInvitation(ctx, &v1pb.CreateInvitationRequest{
			Invitation: &v1pb.Invitation{
				To: store.GenerateUserPath(user2.UUID),
			},
		})
		require.NoError(t, err)
		require.NotZero(t, resp.Id)
		require.Equal(t, store.GenerateUserPath(user1.UUID), resp.From)
		require.Equal(t, store.GenerateUserPath(user2.UUID), resp.To)
		require.Empty(t, resp.AcceptedAt)
		require.Empty(t, resp.DeletedAt)
		require.NotNil(t, resp.CreatedTime)
		require.NotNil(t, resp.UpdatedTime)
	})

	t.Run("MissingToField", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.CreateInvitation(ctx, &v1pb.CreateInvitationRequest{
			Invitation: &v1pb.Invitation{},
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, st.Code())
		require.Contains(t, st.Message(), "to field is required")
	})

	t.Run("InvalidUserPathFormat", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.CreateInvitation(ctx, &v1pb.CreateInvitationRequest{
			Invitation: &v1pb.Invitation{
				To: "invalid-path",
			},
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, st.Code())
		require.Contains(t, st.Message(), "invalid 'to' user path format")
	})

	t.Run("TargetUserNotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.CreateInvitation(ctx, &v1pb.CreateInvitationRequest{
			Invitation: &v1pb.Invitation{
				To: "/users/nonexistent-uuid",
			},
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
		require.Contains(t, st.Message(), "target user not found")
	})

	t.Run("SelfInvitation", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.CreateInvitation(ctx, &v1pb.CreateInvitationRequest{
			Invitation: &v1pb.Invitation{
				To: store.GenerateUserPath(user1.UUID),
			},
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, st.Code())
		require.Contains(t, st.Message(), "cannot invite yourself")
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.CreateInvitation(ctx, &v1pb.CreateInvitationRequest{
			Invitation: &v1pb.Invitation{
				To: store.GenerateUserPath(user2.UUID),
			},
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestInvitationService_AcceptInvitation(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupInvitationTest(t)

	// Create test invitation
	invitation, err := s.Store.CreateInvitation(ctx, &storepb.Invitation{
		From: store.GenerateUserPath(user1.UUID),
		To:   store.GenerateUserPath(user2.UUID),
	})
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		ctx := setUserContext(ctx, user2)
		resp, err := s.AcceptInvitation(ctx, &v1pb.AcceptInvitationRequest{
			Id: invitation.Id,
		})
		require.NoError(t, err)
		require.Equal(t, invitation.Id, resp.Id)
		require.NotEmpty(t, resp.AcceptedAt)

		// Verify it's a valid ISO 8601 timestamp
		_, parseErr := time.Parse(time.RFC3339, resp.AcceptedAt)
		require.NoError(t, parseErr)
	})

	t.Run("InviterCannotAccept", func(t *testing.T) {
		// Create new invitation for this test
		newInvitation, err := s.Store.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(user1.UUID),
			To:   store.GenerateUserPath(user2.UUID),
		})
		require.NoError(t, err)

		ctx := setUserContext(ctx, user1)
		_, err = s.AcceptInvitation(ctx, &v1pb.AcceptInvitationRequest{
			Id: newInvitation.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.PermissionDenied, st.Code())
		require.Contains(t, st.Message(), "only the invited user can accept this invitation")
	})

	t.Run("NotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user2)
		_, err := s.AcceptInvitation(ctx, &v1pb.AcceptInvitationRequest{
			Id: 99999,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.AcceptInvitation(ctx, &v1pb.AcceptInvitationRequest{
			Id: invitation.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestInvitationService_DeleteInvitation(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupInvitationTest(t)

	t.Run("Success_Inviter", func(t *testing.T) {
		// Create test invitation
		invitation, err := s.Store.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(user1.UUID),
			To:   store.GenerateUserPath(user2.UUID),
		})
		require.NoError(t, err)

		ctx := setUserContext(ctx, user1)
		_, err = s.DeleteInvitation(ctx, &v1pb.DeleteInvitationRequest{
			Id: invitation.Id,
		})
		require.NoError(t, err)

		// Verify the invitation is soft deleted
		updated, err := s.Store.GetInvitation(ctx, &store.FindInvitation{
			ID: &invitation.Id,
		})
		require.NoError(t, err)
		require.NotEmpty(t, updated.DeletedAt)

		// Verify it's a valid ISO 8601 timestamp
		_, parseErr := time.Parse(time.RFC3339, updated.DeletedAt)
		require.NoError(t, parseErr)
	})

	t.Run("Success_Invitee", func(t *testing.T) {
		// Create test invitation
		invitation, err := s.Store.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(user1.UUID),
			To:   store.GenerateUserPath(user2.UUID),
		})
		require.NoError(t, err)

		ctx := setUserContext(ctx, user2)
		_, err = s.DeleteInvitation(ctx, &v1pb.DeleteInvitationRequest{
			Id: invitation.Id,
		})
		require.NoError(t, err)

		// Verify the invitation is soft deleted
		updated, err := s.Store.GetInvitation(ctx, &store.FindInvitation{
			ID: &invitation.Id,
		})
		require.NoError(t, err)
		require.NotEmpty(t, updated.DeletedAt)
	})

	t.Run("NotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.DeleteInvitation(ctx, &v1pb.DeleteInvitationRequest{
			Id: 99999,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		// Create test invitation
		invitation, err := s.Store.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(user1.UUID),
			To:   store.GenerateUserPath(user2.UUID),
		})
		require.NoError(t, err)

		// Create a third user who is not involved in the invitation
		thirdUser, err := s.Store.CreateUser(ctx, &store.User{
			Role:         store.RoleUser,
			Email:        "test4@test.com",
			Nickname:     "test_user_4",
			PasswordHash: "test-hash",
			UUID:         "user-uuid-4",
		})
		require.NoError(t, err)

		ctx := setUserContext(ctx, thirdUser)
		_, err = s.DeleteInvitation(ctx, &v1pb.DeleteInvitationRequest{
			Id: invitation.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		// Create test invitation
		invitation, err := s.Store.CreateInvitation(ctx, &storepb.Invitation{
			From: store.GenerateUserPath(user1.UUID),
			To:   store.GenerateUserPath(user2.UUID),
		})
		require.NoError(t, err)

		_, err = s.DeleteInvitation(ctx, &v1pb.DeleteInvitationRequest{
			Id: invitation.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

// Helper function to setup test service and users
func setupInvitationTest(t *testing.T) (*APIV1Service, *store.User, *store.User) {
	ctx := context.Background()
	ts := teststore.NewTestingStore(ctx, t)
	s := &APIV1Service{Store: ts}

	// Create test users
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("test-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	user1, err := ts.CreateUser(ctx, &store.User{
		Role:         store.RoleAdmin,
		Email:        "test1@test.com",
		Nickname:     "test_user_1",
		PasswordHash: string(passwordHash),
		UUID:         "user-uuid-1",
	})
	require.NoError(t, err)

	user2, err := ts.CreateUser(ctx, &store.User{
		Role:         store.RoleUser,
		Email:        "test2@test.com",
		Nickname:     "test_user_2",
		PasswordHash: string(passwordHash),
		UUID:         "user-uuid-2",
	})
	require.NoError(t, err)

	return s, user1, user2
}

// Helper function to set user context for authentication
func setUserContext(ctx context.Context, user *store.User) context.Context {
	// This simulates the authenticated user context that would be set by the auth middleware
	return context.WithValue(ctx, userIDContextKey, user.ID)
}