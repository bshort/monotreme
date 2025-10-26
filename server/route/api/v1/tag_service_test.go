package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	storepb "github.com/bshort/monotreme/proto/gen/store"
	v1pb "github.com/bshort/monotreme/proto/gen/api/v1"
	"github.com/bshort/monotreme/store"
	teststore "github.com/bshort/monotreme/store/test"
)

func TestTagService_ListTags(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupTagTest(t)

	// Create test tags for user1
	tag1, err := s.Store.CreateTag(ctx, &storepb.Tag{
		Uuid:         "tag-uuid-1",
		CreatorId:    user1.ID,
		Name:         "Work",
		Abbreviation: "work",
		Description:  "Work related bookmarks",
	})
	require.NoError(t, err)

	tag2, err := s.Store.CreateTag(ctx, &storepb.Tag{
		Uuid:         "tag-uuid-2",
		CreatorId:    user1.ID,
		Name:         "Personal",
		Abbreviation: "personal",
		Description:  "Personal bookmarks",
	})
	require.NoError(t, err)

	// Create tag for user2 (should not be visible to user1)
	_, err = s.Store.CreateTag(ctx, &storepb.Tag{
		Uuid:         "tag-uuid-3",
		CreatorId:    user2.ID,
		Name:         "User2 Tag",
		Abbreviation: "user2-tag",
		Description:  "User 2's tag",
	})
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.ListTags(ctx, &v1pb.ListTagsRequest{})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(resp.Tags), 2)

		// Verify user1 can only see their own tags
		foundTag1 := false
		foundTag2 := false
		for _, tag := range resp.Tags {
			require.Equal(t, user1.ID, tag.CreatorId)
			if tag.Uuid == tag1.Uuid {
				foundTag1 = true
				require.Equal(t, "Work", tag.Name)
				require.Equal(t, "work", tag.Abbreviation)
			}
			if tag.Uuid == tag2.Uuid {
				foundTag2 = true
				require.Equal(t, "Personal", tag.Name)
				require.Equal(t, "personal", tag.Abbreviation)
			}
		}
		require.True(t, foundTag1)
		require.True(t, foundTag2)
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.ListTags(ctx, &v1pb.ListTagsRequest{})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestTagService_GetTag(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupTagTest(t)

	// Create test tag for user1
	tag, err := s.Store.CreateTag(ctx, &storepb.Tag{
		Uuid:         "tag-uuid-1",
		CreatorId:    user1.ID,
		Name:         "Test Tag",
		Abbreviation: "test-tag",
		Description:  "Test description",
	})
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.GetTag(ctx, &v1pb.GetTagRequest{
			Uuid: tag.Uuid,
		})
		require.NoError(t, err)
		require.Equal(t, tag.Uuid, resp.Uuid)
		require.Equal(t, tag.Name, resp.Name)
		require.Equal(t, tag.Abbreviation, resp.Abbreviation)
		require.Equal(t, tag.Description, resp.Description)
	})

	t.Run("NotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.GetTag(ctx, &v1pb.GetTagRequest{
			Uuid: "nonexistent-uuid",
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		ctx := setUserContext(ctx, user2)
		_, err := s.GetTag(ctx, &v1pb.GetTagRequest{
			Uuid: tag.Uuid,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.GetTag(ctx, &v1pb.GetTagRequest{
			Uuid: tag.Uuid,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestTagService_GetTagByAbbreviation(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupTagTest(t)

	// Create test tag for user1
	tag, err := s.Store.CreateTag(ctx, &storepb.Tag{
		Uuid:         "tag-uuid-1",
		CreatorId:    user1.ID,
		Name:         "Programming",
		Abbreviation: "programming",
		Description:  "Programming resources",
	})
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.GetTagByAbbreviation(ctx, &v1pb.GetTagByAbbreviationRequest{
			Abbreviation: "programming",
		})
		require.NoError(t, err)
		require.Equal(t, tag.Uuid, resp.Uuid)
		require.Equal(t, tag.Name, resp.Name)
		require.Equal(t, tag.Abbreviation, resp.Abbreviation)
	})

	t.Run("NotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.GetTagByAbbreviation(ctx, &v1pb.GetTagByAbbreviationRequest{
			Abbreviation: "nonexistent",
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		ctx := setUserContext(ctx, user2)
		_, err := s.GetTagByAbbreviation(ctx, &v1pb.GetTagByAbbreviationRequest{
			Abbreviation: "programming",
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.GetTagByAbbreviation(ctx, &v1pb.GetTagByAbbreviationRequest{
			Abbreviation: "programming",
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestTagService_CreateTag(t *testing.T) {
	ctx := context.Background()
	s, user1, _ := setupTagTest(t)

	t.Run("Success_WithAutoGeneratedAbbreviation", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.CreateTag(ctx, &v1pb.CreateTagRequest{
			Tag: &v1pb.Tag{
				Name:        "My Test Tag",
				Description: "A test tag",
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, resp.Uuid)
		require.Equal(t, "My Test Tag", resp.Name)
		require.Equal(t, "my-test-tag", resp.Abbreviation) // Auto-generated
		require.Equal(t, "A test tag", resp.Description)
		require.Equal(t, user1.ID, resp.CreatorId)
		require.NotNil(t, resp.CreatedTime)
		require.NotNil(t, resp.UpdatedTime)
	})

	t.Run("Success_WithCustomAbbreviation", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.CreateTag(ctx, &v1pb.CreateTagRequest{
			Tag: &v1pb.Tag{
				Name:         "Custom Tag",
				Abbreviation: "custom",
				Description:  "Custom abbreviation",
			},
		})
		require.NoError(t, err)
		require.Equal(t, "Custom Tag", resp.Name)
		require.Equal(t, "custom", resp.Abbreviation)
	})

	t.Run("Success_WithSpecialCharacters", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.CreateTag(ctx, &v1pb.CreateTagRequest{
			Tag: &v1pb.Tag{
				Name:        "Node.js & JavaScript!",
				Description: "JS resources",
			},
		})
		require.NoError(t, err)
		require.Equal(t, "Node.js & JavaScript!", resp.Name)
		require.Equal(t, "nodejs-javascript", resp.Abbreviation) // Special chars removed
	})

	t.Run("MissingName", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.CreateTag(ctx, &v1pb.CreateTagRequest{
			Tag: &v1pb.Tag{
				Description: "No name provided",
			},
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, st.Code())
		require.Contains(t, st.Message(), "name is required")
	})

	t.Run("DuplicateAbbreviation", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)

		// Create first tag
		_, err := s.CreateTag(ctx, &v1pb.CreateTagRequest{
			Tag: &v1pb.Tag{
				Name:        "Duplicate Test",
				Description: "First tag",
			},
		})
		require.NoError(t, err)

		// Try to create second tag with same abbreviation
		_, err = s.CreateTag(ctx, &v1pb.CreateTagRequest{
			Tag: &v1pb.Tag{
				Name:        "duplicate test", // Will generate same abbreviation
				Description: "Second tag",
			},
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.AlreadyExists, st.Code())
		require.Contains(t, st.Message(), "already exists")
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.CreateTag(ctx, &v1pb.CreateTagRequest{
			Tag: &v1pb.Tag{
				Name: "Test",
			},
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestTagService_UpdateTag(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupTagTest(t)

	// Create test tag
	tag, err := s.Store.CreateTag(ctx, &storepb.Tag{
		Uuid:         "tag-uuid-1",
		CreatorId:    user1.ID,
		Name:         "Original Name",
		Abbreviation: "original-name",
		Description:  "Original description",
	})
	require.NoError(t, err)

	t.Run("Success_UpdateName", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.UpdateTag(ctx, &v1pb.UpdateTagRequest{
			Tag: &v1pb.Tag{
				Uuid: tag.Uuid,
				Name: "Updated Name",
			},
			UpdateMask: []string{"name"},
		})
		require.NoError(t, err)
		require.Equal(t, "Updated Name", resp.Name)
		require.Equal(t, "updated-name", resp.Abbreviation) // Auto-regenerated
	})

	t.Run("Success_UpdateDescription", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.UpdateTag(ctx, &v1pb.UpdateTagRequest{
			Tag: &v1pb.Tag{
				Uuid:        tag.Uuid,
				Description: "New description",
			},
			UpdateMask: []string{"description"},
		})
		require.NoError(t, err)
		require.Equal(t, "New description", resp.Description)
	})

	t.Run("MissingUpdateMask", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.UpdateTag(ctx, &v1pb.UpdateTagRequest{
			Tag: &v1pb.Tag{
				Uuid: tag.Uuid,
				Name: "Test",
			},
			UpdateMask: []string{},
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("NotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.UpdateTag(ctx, &v1pb.UpdateTagRequest{
			Tag: &v1pb.Tag{
				Uuid: "nonexistent",
				Name: "Test",
			},
			UpdateMask: []string{"name"},
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		ctx := setUserContext(ctx, user2)
		_, err := s.UpdateTag(ctx, &v1pb.UpdateTagRequest{
			Tag: &v1pb.Tag{
				Uuid: tag.Uuid,
				Name: "Hacked",
			},
			UpdateMask: []string{"name"},
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.UpdateTag(ctx, &v1pb.UpdateTagRequest{
			Tag: &v1pb.Tag{
				Uuid: tag.Uuid,
				Name: "Test",
			},
			UpdateMask: []string{"name"},
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestTagService_DeleteTag(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupTagTest(t)

	t.Run("Success", func(t *testing.T) {
		// Create test tag
		tag, err := s.Store.CreateTag(ctx, &storepb.Tag{
			Uuid:         "tag-to-delete",
			CreatorId:    user1.ID,
			Name:         "Delete Me",
			Abbreviation: "delete-me",
			Description:  "Test",
		})
		require.NoError(t, err)

		ctx := setUserContext(ctx, user1)
		_, err = s.DeleteTag(ctx, &v1pb.DeleteTagRequest{
			Uuid: tag.Uuid,
		})
		require.NoError(t, err)

		// Verify tag is deleted
		deletedTag, err := s.Store.GetTag(ctx, &store.FindTag{
			UUID: &tag.Uuid,
		})
		// Tag should be nil after deletion
		require.Nil(t, deletedTag)
	})

	t.Run("NotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.DeleteTag(ctx, &v1pb.DeleteTagRequest{
			Uuid: "nonexistent",
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		// Create tag for user1
		tag, err := s.Store.CreateTag(ctx, &storepb.Tag{
			Uuid:         "tag-permission-test",
			CreatorId:    user1.ID,
			Name:         "User1 Tag",
			Abbreviation: "user1-tag",
			Description:  "Test",
		})
		require.NoError(t, err)

		// Try to delete as user2
		ctx := setUserContext(ctx, user2)
		_, err = s.DeleteTag(ctx, &v1pb.DeleteTagRequest{
			Uuid: tag.Uuid,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.DeleteTag(ctx, &v1pb.DeleteTagRequest{
			Uuid: "any-uuid",
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestTagService_AddBookmarkToTag(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupTagTest(t)

	// Create test tag
	tag, err := s.Store.CreateTag(ctx, &storepb.Tag{
		Uuid:         "tag-uuid-1",
		CreatorId:    user1.ID,
		Name:         "Test Tag",
		Abbreviation: "test-tag",
		Description:  "Test",
	})
	require.NoError(t, err)

	// Create test shortcut
	shortcut, err := s.Store.CreateShortcut(ctx, &storepb.Shortcut{
		Uuid:       "shortcut-uuid-1",
		CreatorId:  user1.ID,
		Name:       "test",
		Link:       "https://example.com",
		Visibility: storepb.Visibility_WORKSPACE,
		OgMetadata: &storepb.OpenGraphMetadata{},
	})
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.AddBookmarkToTag(ctx, &v1pb.AddBookmarkToTagRequest{
			TagUuid:    tag.Uuid,
			ShortcutId: shortcut.Id,
		})
		require.NoError(t, err)

		// Verify relationship exists
		bookmarkTags, err := s.Store.ListBookmarkTags(ctx, &store.FindBookmarkTag{
			ShortcutID: &shortcut.Id,
			TagUUID:    &tag.Uuid,
		})
		require.NoError(t, err)
		require.Len(t, bookmarkTags, 1)
	})

	t.Run("TagNotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.AddBookmarkToTag(ctx, &v1pb.AddBookmarkToTagRequest{
			TagUuid:    "nonexistent",
			ShortcutId: shortcut.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("ShortcutNotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.AddBookmarkToTag(ctx, &v1pb.AddBookmarkToTagRequest{
			TagUuid:    tag.Uuid,
			ShortcutId: 99999,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("PermissionDenied_Tag", func(t *testing.T) {
		ctx := setUserContext(ctx, user2)
		_, err := s.AddBookmarkToTag(ctx, &v1pb.AddBookmarkToTagRequest{
			TagUuid:    tag.Uuid,
			ShortcutId: shortcut.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.AddBookmarkToTag(ctx, &v1pb.AddBookmarkToTagRequest{
			TagUuid:    tag.Uuid,
			ShortcutId: shortcut.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestTagService_RemoveBookmarkFromTag(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupTagTest(t)

	// Create test tag
	tag, err := s.Store.CreateTag(ctx, &storepb.Tag{
		Uuid:         "tag-uuid-1",
		CreatorId:    user1.ID,
		Name:         "Test Tag",
		Abbreviation: "test-tag",
		Description:  "Test",
	})
	require.NoError(t, err)

	// Create test shortcut
	shortcut, err := s.Store.CreateShortcut(ctx, &storepb.Shortcut{
		Uuid:       "shortcut-uuid-1",
		CreatorId:  user1.ID,
		Name:       "test",
		Link:       "https://example.com",
		Visibility: storepb.Visibility_WORKSPACE,
		OgMetadata: &storepb.OpenGraphMetadata{},
	})
	require.NoError(t, err)

	// Create relationship
	_, err = s.Store.CreateBookmarkTag(ctx, &store.CreateBookmarkTag{
		ShortcutID: shortcut.Id,
		TagUUID:    tag.Uuid,
		UserID:     user1.ID,
	})
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.RemoveBookmarkFromTag(ctx, &v1pb.RemoveBookmarkFromTagRequest{
			TagUuid:    tag.Uuid,
			ShortcutId: shortcut.Id,
		})
		require.NoError(t, err)

		// Verify relationship is removed
		bookmarkTags, err := s.Store.ListBookmarkTags(ctx, &store.FindBookmarkTag{
			ShortcutID: &shortcut.Id,
			TagUUID:    &tag.Uuid,
		})
		require.NoError(t, err)
		require.Len(t, bookmarkTags, 0)
	})

	t.Run("TagNotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.RemoveBookmarkFromTag(ctx, &v1pb.RemoveBookmarkFromTagRequest{
			TagUuid:    "nonexistent",
			ShortcutId: shortcut.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		ctx := setUserContext(ctx, user2)
		_, err := s.RemoveBookmarkFromTag(ctx, &v1pb.RemoveBookmarkFromTagRequest{
			TagUuid:    tag.Uuid,
			ShortcutId: shortcut.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.RemoveBookmarkFromTag(ctx, &v1pb.RemoveBookmarkFromTagRequest{
			TagUuid:    tag.Uuid,
			ShortcutId: shortcut.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestTagService_ListBookmarksForTag(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupTagTest(t)

	// Create test tag
	tag, err := s.Store.CreateTag(ctx, &storepb.Tag{
		Uuid:         "tag-uuid-1",
		CreatorId:    user1.ID,
		Name:         "Test Tag",
		Abbreviation: "test-tag",
		Description:  "Test",
	})
	require.NoError(t, err)

	// Create test shortcuts
	shortcut1, err := s.Store.CreateShortcut(ctx, &storepb.Shortcut{
		Uuid:       "shortcut-uuid-1",
		CreatorId:  user1.ID,
		Name:       "test1",
		Link:       "https://example1.com",
		Visibility: storepb.Visibility_WORKSPACE,
		OgMetadata: &storepb.OpenGraphMetadata{},
	})
	require.NoError(t, err)

	shortcut2, err := s.Store.CreateShortcut(ctx, &storepb.Shortcut{
		Uuid:       "shortcut-uuid-2",
		CreatorId:  user1.ID,
		Name:       "test2",
		Link:       "https://example2.com",
		Visibility: storepb.Visibility_WORKSPACE,
		OgMetadata: &storepb.OpenGraphMetadata{},
	})
	require.NoError(t, err)

	// Create relationships
	_, err = s.Store.CreateBookmarkTag(ctx, &store.CreateBookmarkTag{
		ShortcutID: shortcut1.Id,
		TagUUID:    tag.Uuid,
		UserID:     user1.ID,
	})
	require.NoError(t, err)

	_, err = s.Store.CreateBookmarkTag(ctx, &store.CreateBookmarkTag{
		ShortcutID: shortcut2.Id,
		TagUUID:    tag.Uuid,
		UserID:     user1.ID,
	})
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.ListBookmarksForTag(ctx, &v1pb.ListBookmarksForTagRequest{
			TagUuid: tag.Uuid,
		})
		require.NoError(t, err)
		require.Len(t, resp.Shortcuts, 2)

		// Verify shortcuts are returned
		foundShortcut1 := false
		foundShortcut2 := false
		for _, sc := range resp.Shortcuts {
			if sc.Id == shortcut1.Id {
				foundShortcut1 = true
				require.Equal(t, "test1", sc.Name)
			}
			if sc.Id == shortcut2.Id {
				foundShortcut2 = true
				require.Equal(t, "test2", sc.Name)
			}
		}
		require.True(t, foundShortcut1)
		require.True(t, foundShortcut2)
	})

	t.Run("TagNotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.ListBookmarksForTag(ctx, &v1pb.ListBookmarksForTagRequest{
			TagUuid: "nonexistent",
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		ctx := setUserContext(ctx, user2)
		_, err := s.ListBookmarksForTag(ctx, &v1pb.ListBookmarksForTagRequest{
			TagUuid: tag.Uuid,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.ListBookmarksForTag(ctx, &v1pb.ListBookmarksForTagRequest{
			TagUuid: tag.Uuid,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestTagService_ListTagsForBookmark(t *testing.T) {
	ctx := context.Background()
	s, user1, user2 := setupTagTest(t)

	// Create test tags
	tag1, err := s.Store.CreateTag(ctx, &storepb.Tag{
		Uuid:         "tag-uuid-1",
		CreatorId:    user1.ID,
		Name:         "Tag 1",
		Abbreviation: "tag-1",
		Description:  "Test",
	})
	require.NoError(t, err)

	tag2, err := s.Store.CreateTag(ctx, &storepb.Tag{
		Uuid:         "tag-uuid-2",
		CreatorId:    user1.ID,
		Name:         "Tag 2",
		Abbreviation: "tag-2",
		Description:  "Test",
	})
	require.NoError(t, err)

	// Create test shortcut
	shortcut, err := s.Store.CreateShortcut(ctx, &storepb.Shortcut{
		Uuid:       "shortcut-uuid-1",
		CreatorId:  user1.ID,
		Name:       "test",
		Link:       "https://example.com",
		Visibility: storepb.Visibility_WORKSPACE,
		OgMetadata: &storepb.OpenGraphMetadata{},
	})
	require.NoError(t, err)

	// Create relationships
	_, err = s.Store.CreateBookmarkTag(ctx, &store.CreateBookmarkTag{
		ShortcutID: shortcut.Id,
		TagUUID:    tag1.Uuid,
		UserID:     user1.ID,
	})
	require.NoError(t, err)

	_, err = s.Store.CreateBookmarkTag(ctx, &store.CreateBookmarkTag{
		ShortcutID: shortcut.Id,
		TagUUID:    tag2.Uuid,
		UserID:     user1.ID,
	})
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		resp, err := s.ListTagsForBookmark(ctx, &v1pb.ListTagsForBookmarkRequest{
			ShortcutId: shortcut.Id,
		})
		require.NoError(t, err)
		require.Len(t, resp.Tags, 2)

		// Verify tags are returned
		foundTag1 := false
		foundTag2 := false
		for _, tag := range resp.Tags {
			if tag.Uuid == tag1.Uuid {
				foundTag1 = true
				require.Equal(t, "Tag 1", tag.Name)
			}
			if tag.Uuid == tag2.Uuid {
				foundTag2 = true
				require.Equal(t, "Tag 2", tag.Name)
			}
		}
		require.True(t, foundTag1)
		require.True(t, foundTag2)
	})

	t.Run("ShortcutNotFound", func(t *testing.T) {
		ctx := setUserContext(ctx, user1)
		_, err := s.ListTagsForBookmark(ctx, &v1pb.ListTagsForBookmarkRequest{
			ShortcutId: 99999,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		ctx := setUserContext(ctx, user2)
		_, err := s.ListTagsForBookmark(ctx, &v1pb.ListTagsForBookmarkRequest{
			ShortcutId: shortcut.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("Unauthenticated", func(t *testing.T) {
		_, err := s.ListTagsForBookmark(ctx, &v1pb.ListTagsForBookmarkRequest{
			ShortcutId: shortcut.Id,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.Unauthenticated, st.Code())
	})
}

// Helper function to setup test service and users
func setupTagTest(t *testing.T) (*APIV1Service, *store.User, *store.User) {
	ctx := context.Background()
	ts := teststore.NewTestingStore(ctx, t)
	s := &APIV1Service{Store: ts}

	// Create test users
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("test-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	user1, err := ts.CreateUser(ctx, &store.User{
		Role:         store.RoleUser,
		Email:        "tagtest1@test.com",
		Nickname:     "tag_test_user_1",
		PasswordHash: string(passwordHash),
		UUID:         "tag-user-uuid-1",
	})
	require.NoError(t, err)

	user2, err := ts.CreateUser(ctx, &store.User{
		Role:         store.RoleUser,
		Email:        "tagtest2@test.com",
		Nickname:     "tag_test_user_2",
		PasswordHash: string(passwordHash),
		UUID:         "tag-user-uuid-2",
	})
	require.NoError(t, err)

	return s, user1, user2
}
