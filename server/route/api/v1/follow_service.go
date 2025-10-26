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

func (s *APIV1Service) ListFollowing(ctx context.Context, request *v1pb.ListFollowingRequest) (*v1pb.ListFollowingResponse, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// Get list of users the current user is following
	followings, err := s.Store.ListFollowing(ctx, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list following: %v", err)
	}

	// Convert to proto format
	convertedFollowings := make([]*v1pb.Following, 0, len(followings))
	for _, following := range followings {
		convertedFollowing, err := s.convertFollowingToProto(ctx, following)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert following: %v", err)
		}
		convertedFollowings = append(convertedFollowings, convertedFollowing)
	}

	return &v1pb.ListFollowingResponse{
		Following: convertedFollowings,
	}, nil
}

func (s *APIV1Service) FollowUser(ctx context.Context, request *v1pb.FollowUserRequest) (*v1pb.Following, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// Find the user to follow
	var followingUser *store.User
	if request.UserEmail != nil && *request.UserEmail != "" {
		followingUser, err = s.Store.GetUser(ctx, &store.FindUser{
			Email: request.UserEmail,
		})
	} else if request.UserId != nil {
		followingUser, err = s.Store.GetUser(ctx, &store.FindUser{
			ID: request.UserId,
		})
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "either user_email or user_id must be provided")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}
	if followingUser == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	if followingUser.ID == user.ID {
		return nil, status.Errorf(codes.InvalidArgument, "cannot follow yourself")
	}

	// Create following relationship
	following, err := s.Store.CreateFollowing(ctx, user.ID, followingUser.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create following: %v", err)
	}

	convertedFollowing, err := s.convertFollowingToProto(ctx, following)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert following: %v", err)
	}

	return convertedFollowing, nil
}

func (s *APIV1Service) UnfollowUser(ctx context.Context, request *v1pb.UnfollowUserRequest) (*emptypb.Empty, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// Get the following relationship
	following, err := s.Store.GetFollowing(ctx, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get following: %v", err)
	}
	if following == nil {
		return nil, status.Errorf(codes.NotFound, "following relationship not found")
	}

	// Only the follower can unfollow
	if following.FollowerID != user.ID {
		return nil, status.Errorf(codes.PermissionDenied, "you can only unfollow users you are following")
	}

	// Delete the following relationship
	err = s.Store.DeleteFollowing(ctx, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete following: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *APIV1Service) ListFollowers(ctx context.Context, request *v1pb.ListFollowersRequest) (*v1pb.ListFollowersResponse, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// Get list of users following the current user
	followers, err := s.Store.ListFollowers(ctx, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list followers: %v", err)
	}

	// Convert to proto format
	convertedFollowers := make([]*v1pb.Following, 0, len(followers))
	for _, follower := range followers {
		convertedFollower, err := s.convertFollowingToProto(ctx, follower)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert follower: %v", err)
		}
		convertedFollowers = append(convertedFollowers, convertedFollower)
	}

	return &v1pb.ListFollowersResponse{
		Followers: convertedFollowers,
	}, nil
}

func (s *APIV1Service) ListFollowingUserShortcuts(ctx context.Context, request *v1pb.ListFollowingUserShortcutsRequest) (*v1pb.ListFollowingUserShortcutsResponse, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// Get public shortcuts from the user being followed
	shortcuts, err := s.Store.GetFollowingUserShortcuts(ctx, user.ID, request.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get shortcuts: %v", err)
	}

	// Convert to proto format
	convertedShortcuts := make([]*v1pb.Shortcut, 0, len(shortcuts))
	for _, shortcut := range shortcuts {
		convertedShortcut, err := s.convertShortcutFromStorepb(ctx, shortcut)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert shortcut: %v", err)
		}
		convertedShortcuts = append(convertedShortcuts, convertedShortcut)
	}

	return &v1pb.ListFollowingUserShortcutsResponse{
		Shortcuts: convertedShortcuts,
	}, nil
}

func (s *APIV1Service) convertFollowingToProto(ctx context.Context, following *store.Following) (*v1pb.Following, error) {
	// Get user details
	follower, err := s.Store.GetUser(ctx, &store.FindUser{ID: &following.FollowerID})
	if err != nil {
		return nil, err
	}
	followingUser, err := s.Store.GetUser(ctx, &store.FindUser{ID: &following.FollowingID})
	if err != nil {
		return nil, err
	}

	result := &v1pb.Following{
		Id:          following.ID,
		FollowerId:  following.FollowerID,
		FollowingId: following.FollowingID,
	}

	if following.CreatedTs > 0 {
		result.CreatedTime = timestamppb.New(following.CreatedTs.AsTime())
	}
	if follower != nil {
		result.Follower = convertUserFromStore(follower)
	}
	if followingUser != nil {
		result.Following = convertUserFromStore(followingUser)
	}

	return result, nil
}
