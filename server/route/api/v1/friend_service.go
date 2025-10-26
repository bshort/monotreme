package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1pb "github.com/bshort/monotreme/proto/gen/api/v1"
	"github.com/bshort/monotreme/store"
)

func (s *APIV1Service) ListFriends(ctx context.Context, request *v1pb.ListFriendsRequest) (*v1pb.ListFriendsResponse, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// Get accepted friendships for the current user
	friendships, err := s.Store.ListFriendships(ctx, user.ID, "ACCEPTED")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list friendships: %v", err)
	}

	// Convert to proto format
	convertedFriendships := make([]*v1pb.Friendship, 0, len(friendships))
	for _, friendship := range friendships {
		convertedFriendship, err := s.convertFriendshipToProto(ctx, friendship)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert friendship: %v", err)
		}
		convertedFriendships = append(convertedFriendships, convertedFriendship)
	}

	return &v1pb.ListFriendsResponse{
		Friends: convertedFriendships,
	}, nil
}

func (s *APIV1Service) SendFriendRequest(ctx context.Context, request *v1pb.SendFriendRequestRequest) (*v1pb.Friendship, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// Find the friend user
	var friendUser *store.User
	if request.FriendEmail != nil && *request.FriendEmail != "" {
		friendUser, err = s.Store.GetUser(ctx, &store.FindUser{
			Email: request.FriendEmail,
		})
	} else if request.FriendId != nil {
		friendUser, err = s.Store.GetUser(ctx, &store.FindUser{
			ID: request.FriendId,
		})
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "either friend_email or friend_id must be provided")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find friend user: %v", err)
	}
	if friendUser == nil {
		return nil, status.Errorf(codes.NotFound, "friend user not found")
	}

	if friendUser.ID == user.ID {
		return nil, status.Errorf(codes.InvalidArgument, "cannot send friend request to yourself")
	}

	// Create friendship
	friendship, err := s.Store.CreateFriendship(ctx, user.ID, friendUser.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create friendship: %v", err)
	}

	convertedFriendship, err := s.convertFriendshipToProto(ctx, friendship)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert friendship: %v", err)
	}

	return convertedFriendship, nil
}

func (s *APIV1Service) ListFriendRequests(ctx context.Context, request *v1pb.ListFriendRequestsRequest) (*v1pb.ListFriendRequestsResponse, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	var friendships []*store.Friendship
	switch request.Type {
	case v1pb.FriendRequestType_REQUEST_TYPE_INCOMING:
		// Get pending requests where current user is the friend
		friendships, err = s.Store.ListIncomingFriendRequests(ctx, user.ID)
	case v1pb.FriendRequestType_REQUEST_TYPE_OUTGOING:
		// Get pending requests where current user is the sender
		friendships, err = s.Store.ListOutgoingFriendRequests(ctx, user.ID)
	default:
		// Get all pending requests
		friendships, err = s.Store.ListFriendships(ctx, user.ID, "PENDING")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list friend requests: %v", err)
	}

	convertedFriendships := make([]*v1pb.Friendship, 0, len(friendships))
	for _, friendship := range friendships {
		convertedFriendship, err := s.convertFriendshipToProto(ctx, friendship)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert friendship: %v", err)
		}
		convertedFriendships = append(convertedFriendships, convertedFriendship)
	}

	return &v1pb.ListFriendRequestsResponse{
		FriendRequests: convertedFriendships,
	}, nil
}

func (s *APIV1Service) AcceptFriendRequest(ctx context.Context, request *v1pb.AcceptFriendRequestRequest) (*v1pb.Friendship, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// Get the friendship
	friendship, err := s.Store.GetFriendship(ctx, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get friendship: %v", err)
	}
	if friendship == nil {
		return nil, status.Errorf(codes.NotFound, "friendship not found")
	}

	// Only the friend can accept the request
	if friendship.FriendID != user.ID {
		return nil, status.Errorf(codes.PermissionDenied, "only the recipient can accept a friend request")
	}

	// Accept the friendship
	err = s.Store.AcceptFriendship(ctx, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to accept friendship: %v", err)
	}

	// Get updated friendship
	friendship, err = s.Store.GetFriendship(ctx, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get updated friendship: %v", err)
	}

	convertedFriendship, err := s.convertFriendshipToProto(ctx, friendship)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert friendship: %v", err)
	}

	return convertedFriendship, nil
}

func (s *APIV1Service) RemoveFriend(ctx context.Context, request *v1pb.RemoveFriendRequest) (*emptypb.Empty, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// Get the friendship
	friendship, err := s.Store.GetFriendship(ctx, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get friendship: %v", err)
	}
	if friendship == nil {
		return nil, status.Errorf(codes.NotFound, "friendship not found")
	}

	// Only participants can remove a friendship
	if friendship.UserID != user.ID && friendship.FriendID != user.ID {
		return nil, status.Errorf(codes.PermissionDenied, "you can only remove your own friendships")
	}

	// Delete the friendship
	err = s.Store.DeleteFriendship(ctx, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete friendship: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *APIV1Service) convertFriendshipToProto(ctx context.Context, friendship *store.Friendship) (*v1pb.Friendship, error) {
	// Get user details
	friendUser, err := s.Store.GetUser(ctx, &store.FindUser{ID: &friendship.UserID})
	if err != nil {
		return nil, err
	}
	user, err := s.Store.GetUser(ctx, &store.FindUser{ID: &friendship.FriendID})
	if err != nil {
		return nil, err
	}

	result := &v1pb.Friendship{
		Id:       friendship.ID,
		UserId:   friendship.UserID,
		FriendId: friendship.FriendID,
		Status:   convertFriendshipStatus(friendship.Status),
	}

	if friendship.CreatedTs > 0 {
		result.CreatedTime = timestamppb.New(friendship.CreatedTs.AsTime())
	}
	if friendship.AcceptedTs != nil && *friendship.AcceptedTs > 0 {
		result.AcceptedTime = timestamppb.New(friendship.AcceptedTs.AsTime())
	}
	if friendUser != nil {
		result.User = convertUserFromStore(friendUser)
	}
	if user != nil {
		result.Friend = convertUserFromStore(user)
	}

	return result, nil
}

func convertFriendshipStatus(status string) v1pb.FriendshipStatus {
	switch status {
	case "PENDING":
		return v1pb.FriendshipStatus_PENDING
	case "ACCEPTED":
		return v1pb.FriendshipStatus_ACCEPTED
	case "DECLINED":
		return v1pb.FriendshipStatus_DECLINED
	case "BLOCKED":
		return v1pb.FriendshipStatus_BLOCKED
	default:
		return v1pb.FriendshipStatus_FRIENDSHIP_STATUS_UNSPECIFIED
	}
}
