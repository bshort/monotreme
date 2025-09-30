package v1

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1pb "github.com/bshort/monotreme/proto/gen/api/v1"
	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

func (s *APIV1Service) ListInvitations(ctx context.Context, request *v1pb.ListInvitationsRequest) (*v1pb.ListInvitationsResponse, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	find := &store.FindInvitation{}

	// Apply filters if provided
	if request.From != "" {
		find.From = &request.From
	}
	if request.To != "" {
		find.To = &request.To
	}
	if request.Accepted != nil {
		accepted := *request.Accepted
		find.Accepted = &accepted
	}
	if request.Deleted != nil {
		deleted := *request.Deleted
		find.Deleted = &deleted
	}

	invitations, err := s.Store.ListInvitations(ctx, find)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get invitation list, err: %v", err)
	}

	convertedInvitations := []*v1pb.Invitation{}
	for _, invitation := range invitations {
		// Only show invitations where the current user is involved
		currentUserPath := store.GenerateUserPath(user.UUID)
		if invitation.From == currentUserPath || invitation.To == currentUserPath {
			convertedInvitations = append(convertedInvitations, convertInvitationFromStore(invitation))
		}
	}

	response := &v1pb.ListInvitationsResponse{
		Invitations: convertedInvitations,
	}
	return response, nil
}

func (s *APIV1Service) GetInvitation(ctx context.Context, request *v1pb.GetInvitationRequest) (*v1pb.Invitation, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	invitation, err := s.Store.GetInvitation(ctx, &store.FindInvitation{
		ID: &request.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get invitation: %v", err)
	}
	if invitation == nil {
		return nil, status.Errorf(codes.NotFound, "invitation not found")
	}

	// Check permission - user must be involved in the invitation
	currentUserPath := store.GenerateUserPath(user.UUID)
	if invitation.From != currentUserPath && invitation.To != currentUserPath {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	return convertInvitationFromStore(invitation), nil
}

func (s *APIV1Service) CreateInvitation(ctx context.Context, request *v1pb.CreateInvitationRequest) (*v1pb.Invitation, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	if request.Invitation.To == "" {
		return nil, status.Errorf(codes.InvalidArgument, "to field is required")
	}

	// Validate that the 'to' path is for a valid user
	toUUID := store.ExtractUUIDFromUserPath(request.Invitation.To)
	if toUUID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid 'to' user path format")
	}

	// Find the user with the matching UUID
	var foundTargetUser *store.User
	users, err := s.Store.ListUsers(ctx, &store.FindUser{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	for _, u := range users {
		if u.UUID == toUUID {
			foundTargetUser = u
			break
		}
	}

	if foundTargetUser == nil {
		return nil, status.Errorf(codes.NotFound, "target user not found")
	}

	// Set the from field to the current user
	currentUserPath := store.GenerateUserPath(user.UUID)

	// Prevent self-invitation
	if currentUserPath == request.Invitation.To {
		return nil, status.Errorf(codes.InvalidArgument, "cannot invite yourself")
	}

	invitationCreate := &storepb.Invitation{
		From:       currentUserPath,
		To:         request.Invitation.To,
		AcceptedAt: "",
		DeletedAt:  "",
	}

	invitation, err := s.Store.CreateInvitation(ctx, invitationCreate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create invitation, err: %v", err)
	}

	return convertInvitationFromStore(invitation), nil
}

func (s *APIV1Service) AcceptInvitation(ctx context.Context, request *v1pb.AcceptInvitationRequest) (*v1pb.Invitation, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	invitation, err := s.Store.GetInvitation(ctx, &store.FindInvitation{
		ID: &request.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get invitation: %v", err)
	}
	if invitation == nil {
		return nil, status.Errorf(codes.NotFound, "invitation not found")
	}

	// Check permission - only the invited user can accept
	currentUserPath := store.GenerateUserPath(user.UUID)
	if invitation.To != currentUserPath {
		return nil, status.Errorf(codes.PermissionDenied, "only the invited user can accept this invitation")
	}

	// Check if already accepted
	if invitation.AcceptedAt != "" {
		return nil, status.Errorf(codes.FailedPrecondition, "invitation already accepted")
	}

	// Check if deleted
	if invitation.DeletedAt != "" {
		return nil, status.Errorf(codes.FailedPrecondition, "invitation has been deleted")
	}

	updatedInvitation, err := s.Store.AcceptInvitation(ctx, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to accept invitation, err: %v", err)
	}

	return convertInvitationFromStore(updatedInvitation), nil
}

func (s *APIV1Service) DeleteInvitation(ctx context.Context, request *v1pb.DeleteInvitationRequest) (*emptypb.Empty, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	invitation, err := s.Store.GetInvitation(ctx, &store.FindInvitation{
		ID: &request.Id,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get invitation: %v", err)
	}
	if invitation == nil {
		return nil, status.Errorf(codes.NotFound, "invitation not found")
	}

	// Check permission - only the creator or recipient can delete
	currentUserPath := store.GenerateUserPath(user.UUID)
	if invitation.From != currentUserPath && invitation.To != currentUserPath {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	// Soft delete the invitation
	_, err = s.Store.SoftDeleteInvitation(ctx, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete invitation, err: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func convertInvitationFromStore(invitation *storepb.Invitation) *v1pb.Invitation {
	return &v1pb.Invitation{
		Id:          invitation.Id,
		CreatedTime: timestamppb.New(time.Unix(invitation.CreatedTs, 0)),
		UpdatedTime: timestamppb.New(time.Unix(invitation.UpdatedTs, 0)),
		From:        invitation.From,
		To:          invitation.To,
		AcceptedAt:  invitation.AcceptedAt,
		DeletedAt:   invitation.DeletedAt,
	}
}