package v1

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/google/uuid"

	v1pb "github.com/bshort/monotreme/proto/gen/api/v1"
	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/server/common"
	"github.com/bshort/monotreme/store"
)

func (s *APIV1Service) ListTags(ctx context.Context, _ *v1pb.ListTagsRequest) (*v1pb.ListTagsResponse, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	tags, err := s.Store.ListTags(ctx, &store.FindTag{
		CreatorID: &user.ID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tag list, err: %v", err)
	}

	convertedTags := []*v1pb.Tag{}
	for _, tag := range tags {
		convertedTags = append(convertedTags, convertTagFromStore(tag))
	}

	response := &v1pb.ListTagsResponse{
		Tags: convertedTags,
	}
	return response, nil
}

func (s *APIV1Service) GetTag(ctx context.Context, request *v1pb.GetTagRequest) (*v1pb.Tag, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	tag, err := s.Store.GetTag(ctx, &store.FindTag{
		UUID: &request.Uuid,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tag: %v", err)
	}
	if tag == nil {
		return nil, status.Errorf(codes.NotFound, "tag not found")
	}

	if tag.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	return convertTagFromStore(tag), nil
}

func (s *APIV1Service) GetTagByAbbreviation(ctx context.Context, request *v1pb.GetTagByAbbreviationRequest) (*v1pb.Tag, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	tag, err := s.Store.GetTag(ctx, &store.FindTag{
		Abbreviation: &request.Abbreviation,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tag by abbreviation: %v", err)
	}
	if tag == nil {
		return nil, status.Errorf(codes.NotFound, "tag not found")
	}

	if tag.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	return convertTagFromStore(tag), nil
}

func (s *APIV1Service) CreateTag(ctx context.Context, request *v1pb.CreateTagRequest) (*v1pb.Tag, error) {
	if request.Tag.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	// Generate abbreviation from name if not provided
	abbreviation := request.Tag.Abbreviation
	if abbreviation == "" {
		abbreviation = common.GenerateTagAbbreviation(request.Tag.Name)
	}

	// Check if abbreviation already exists for this user
	existing, err := s.Store.GetTag(ctx, &store.FindTag{
		CreatorID:    &user.ID,
		Abbreviation: &abbreviation,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check existing tag: %v", err)
	}
	if existing != nil {
		return nil, status.Errorf(codes.AlreadyExists, "a tag with abbreviation %s already exists", abbreviation)
	}

	tagCreate := &storepb.Tag{
		Uuid:         uuid.New().String(),
		CreatorId:    user.ID,
		Name:         request.Tag.Name,
		Abbreviation: abbreviation,
		Description:  request.Tag.Description,
	}
	tag, err := s.Store.CreateTag(ctx, tagCreate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create tag, err: %v", err)
	}

	return convertTagFromStore(tag), nil
}

func (s *APIV1Service) UpdateTag(ctx context.Context, request *v1pb.UpdateTagRequest) (*v1pb.Tag, error) {
	if len(request.UpdateMask) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "updateMask is required")
	}

	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	tag, err := s.Store.GetTag(ctx, &store.FindTag{
		UUID: &request.Tag.Uuid,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tag: %v", err)
	}
	if tag == nil {
		return nil, status.Errorf(codes.NotFound, "tag not found")
	}

	if tag.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	update := &store.UpdateTag{
		UUID: tag.Uuid,
	}
	for _, path := range request.UpdateMask {
		switch path {
		case "name":
			update.Name = &request.Tag.Name
			// Regenerate abbreviation if name changes
			abbr := common.GenerateTagAbbreviation(request.Tag.Name)
			update.Abbreviation = &abbr
		case "description":
			update.Description = &request.Tag.Description
		}
	}

	updatedTag, err := s.Store.UpdateTag(ctx, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update tag, err: %v", err)
	}

	return convertTagFromStore(updatedTag), nil
}

func (s *APIV1Service) DeleteTag(ctx context.Context, request *v1pb.DeleteTagRequest) (*emptypb.Empty, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	tag, err := s.Store.GetTag(ctx, &store.FindTag{
		UUID: &request.Uuid,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tag: %v", err)
	}
	if tag == nil {
		return nil, status.Errorf(codes.NotFound, "tag not found")
	}

	if tag.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	if err := s.Store.DeleteTag(ctx, &store.DeleteTag{UUID: request.Uuid}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete tag: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func convertTagFromStore(tag *storepb.Tag) *v1pb.Tag {
	return &v1pb.Tag{
		Uuid:         tag.Uuid,
		CreatorId:    tag.CreatorId,
		CreatedTime:  timestamppb.New(time.Unix(tag.CreatedTs, 0)),
		UpdatedTime:  timestamppb.New(time.Unix(tag.UpdatedTs, 0)),
		Name:         tag.Name,
		Abbreviation: tag.Abbreviation,
		Description:  tag.Description,
	}
}

func (s *APIV1Service) AddBookmarkToTag(ctx context.Context, request *v1pb.AddBookmarkToTagRequest) (*emptypb.Empty, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	// Verify tag exists and user has access
	tag, err := s.Store.GetTag(ctx, &store.FindTag{
		UUID: &request.TagUuid,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tag: %v", err)
	}
	if tag == nil {
		return nil, status.Errorf(codes.NotFound, "tag not found")
	}
	if tag.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	// Verify shortcut exists and user has access
	shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
		ID: &request.ShortcutId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get shortcut: %v", err)
	}
	if shortcut == nil {
		return nil, status.Errorf(codes.NotFound, "shortcut not found")
	}
	if shortcut.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	// Create the relationship
	_, err = s.Store.CreateBookmarkTag(ctx, &store.CreateBookmarkTag{
		ShortcutID: request.ShortcutId,
		TagUUID:    request.TagUuid,
		UserID:     user.ID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add bookmark to tag: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *APIV1Service) RemoveBookmarkFromTag(ctx context.Context, request *v1pb.RemoveBookmarkFromTagRequest) (*emptypb.Empty, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	// Verify tag exists and user has access
	tag, err := s.Store.GetTag(ctx, &store.FindTag{
		UUID: &request.TagUuid,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tag: %v", err)
	}
	if tag == nil {
		return nil, status.Errorf(codes.NotFound, "tag not found")
	}
	if tag.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	// Delete the relationship
	if err := s.Store.DeleteBookmarkTag(ctx, &store.DeleteBookmarkTag{
		ShortcutID: request.ShortcutId,
		TagUUID:    request.TagUuid,
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove bookmark from tag: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *APIV1Service) ListBookmarksForTag(ctx context.Context, request *v1pb.ListBookmarksForTagRequest) (*v1pb.ListBookmarksForTagResponse, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	// Verify tag exists and user has access
	tag, err := s.Store.GetTag(ctx, &store.FindTag{
		UUID: &request.TagUuid,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get tag: %v", err)
	}
	if tag == nil {
		return nil, status.Errorf(codes.NotFound, "tag not found")
	}
	if tag.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	// Get all bookmark-tag relationships for this tag
	bookmarkTags, err := s.Store.ListBookmarkTags(ctx, &store.FindBookmarkTag{
		TagUUID: &request.TagUuid,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list bookmark tags: %v", err)
	}

	// Get the shortcuts
	shortcuts := []*v1pb.Shortcut{}
	for _, bt := range bookmarkTags {
		shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
			ID: &bt.ShortcutId,
		})
		if err != nil {
			continue // Skip shortcuts that can't be found
		}
		if shortcut != nil {
			convertedShortcut, err := s.convertShortcutFromStorepb(ctx, shortcut)
			if err != nil {
				continue // Skip shortcuts that can't be converted
			}
			shortcuts = append(shortcuts, convertedShortcut)
		}
	}

	return &v1pb.ListBookmarksForTagResponse{
		Shortcuts: shortcuts,
	}, nil
}

func (s *APIV1Service) ListTagsForBookmark(ctx context.Context, request *v1pb.ListTagsForBookmarkRequest) (*v1pb.ListTagsForBookmarkResponse, error) {
	user, err := getCurrentUser(ctx, s.Store)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get current user: %v", err)
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	// Verify shortcut exists and user has access
	shortcut, err := s.Store.GetShortcut(ctx, &store.FindShortcut{
		ID: &request.ShortcutId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get shortcut: %v", err)
	}
	if shortcut == nil {
		return nil, status.Errorf(codes.NotFound, "shortcut not found")
	}
	if shortcut.CreatorId != user.ID && user.Role != store.RoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}

	// Get all bookmark-tag relationships for this shortcut
	bookmarkTags, err := s.Store.ListBookmarkTags(ctx, &store.FindBookmarkTag{
		ShortcutID: &request.ShortcutId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list bookmark tags: %v", err)
	}

	// Get the tags
	tags := []*v1pb.Tag{}
	for _, bt := range bookmarkTags {
		tag, err := s.Store.GetTag(ctx, &store.FindTag{
			UUID: &bt.TagUuid,
		})
		if err != nil {
			continue // Skip tags that can't be found
		}
		if tag != nil {
			tags = append(tags, convertTagFromStore(tag))
		}
	}

	return &v1pb.ListTagsForBookmarkResponse{
		Tags: tags,
	}, nil
}
