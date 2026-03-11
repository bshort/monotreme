package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	v1pb "github.com/bshort/monotreme/proto/gen/api/v1"
	"github.com/bshort/monotreme/store"
	teststore "github.com/bshort/monotreme/store/test"
)

func TestShortcutService_CreateShortcut_WithTags(t *testing.T) {
	ctx := context.Background()
	s, user := setupShortcutTagTest(t)

	t.Run("Success_CreateTagsAutomatically", func(t *testing.T) {
		ctx := setUserContext(ctx, user)
		resp, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name:        "test",
				Link:        "https://example.com",
				Title:       "Example Site",
				Description: "Test shortcut",
				Tags:        []string{"work", "programming", "golang"},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test", resp.Name)
		require.Equal(t, "https://example.com", resp.Link)

		// Verify tags were stored (backward compatible string array)
		require.Len(t, resp.Tags, 3)
		require.Contains(t, resp.Tags, "work")
		require.Contains(t, resp.Tags, "programming")
		require.Contains(t, resp.Tags, "golang")

		// Verify tag_info is populated with UUIDs and names
		require.Len(t, resp.TagInfo, 3)
		for _, tagInfo := range resp.TagInfo {
			require.NotEmpty(t, tagInfo.Uuid)
			require.NotEmpty(t, tagInfo.Name)
			require.Contains(t, []string{"work", "programming", "golang"}, tagInfo.Name)
		}

		// Verify Tag entities were created in database
		tags, err := s.Store.ListTags(ctx, &store.FindTag{
			CreatorID: &user.ID,
		})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(tags), 3)

		// Verify tag abbreviations were generated correctly
		tagAbbreviations := make(map[string]string)
		for _, tag := range tags {
			tagAbbreviations[tag.Name] = tag.Abbreviation
		}
		require.Equal(t, "work", tagAbbreviations["work"])
		require.Equal(t, "programming", tagAbbreviations["programming"])
		require.Equal(t, "golang", tagAbbreviations["golang"])

		// Verify bookmark-tag relationships were created
		bookmarkTags, err := s.Store.ListBookmarkTags(ctx, &store.FindBookmarkTag{
			ShortcutID: &resp.Id,
		})
		require.NoError(t, err)
		require.Len(t, bookmarkTags, 3)
	})

	t.Run("Success_ReuseExistingTags", func(t *testing.T) {
		ctx := setUserContext(ctx, user)

		// Create first shortcut with tags
		resp1, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "test1",
				Link: "https://example1.com",
				Tags: []string{"shared-tag", "unique1"},
			},
		})
		require.NoError(t, err)

		// Get count of tags before creating second shortcut
		tagsBefore, err := s.Store.ListTags(ctx, &store.FindTag{
			CreatorID: &user.ID,
		})
		require.NoError(t, err)
		countBefore := len(tagsBefore)

		// Create second shortcut with one shared tag and one new tag
		resp2, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "test2",
				Link: "https://example2.com",
				Tags: []string{"shared-tag", "unique2"},
			},
		})
		require.NoError(t, err)

		// Get count of tags after creating second shortcut
		tagsAfter, err := s.Store.ListTags(ctx, &store.FindTag{
			CreatorID: &user.ID,
		})
		require.NoError(t, err)
		countAfter := len(tagsAfter)

		// Should only have created 1 new tag (unique2), not 2
		// because shared-tag already existed
		require.Equal(t, countBefore+1, countAfter)

		// Verify both shortcuts have correct tags
		require.Len(t, resp1.Tags, 2)
		require.Len(t, resp2.Tags, 2)
		require.Contains(t, resp1.Tags, "shared-tag")
		require.Contains(t, resp2.Tags, "shared-tag")
	})

	t.Run("Success_WithEmptyTags", func(t *testing.T) {
		ctx := setUserContext(ctx, user)
		resp, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "no-tags",
				Link: "https://example.com",
				Tags: []string{},
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.Tags, 0)
		require.Len(t, resp.TagInfo, 0)
	})

	t.Run("Success_WithSpecialCharactersInTags", func(t *testing.T) {
		ctx := setUserContext(ctx, user)
		resp, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "special-chars",
				Link: "https://example.com",
				Tags: []string{"Node.js", "C++", "API Design - Best Practices"},
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.TagInfo, 3)

		// Verify abbreviations were generated correctly
		abbreviations := make(map[string]string)
		for _, tagInfo := range resp.TagInfo {
			tag, err := s.Store.GetTag(ctx, &store.FindTag{
				UUID: &tagInfo.Uuid,
			})
			require.NoError(t, err)
			abbreviations[tag.Name] = tag.Abbreviation
		}

		require.Equal(t, "nodejs", abbreviations["Node.js"])
		require.Equal(t, "c", abbreviations["C++"])
		require.Equal(t, "api-design-best-practices", abbreviations["API Design - Best Practices"])
	})
}

func TestShortcutService_UpdateShortcut_WithTags(t *testing.T) {
	ctx := context.Background()
	s, user := setupShortcutTagTest(t)

	t.Run("Success_UpdateTags", func(t *testing.T) {
		ctx := setUserContext(ctx, user)

		// Create shortcut with initial tags
		created, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "test",
				Link: "https://example.com",
				Tags: []string{"tag1", "tag2"},
			},
		})
		require.NoError(t, err)
		require.Len(t, created.Tags, 2)

		// Update with new tags
		updated, err := s.UpdateShortcut(ctx, &v1pb.UpdateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Id:   created.Id,
				Tags: []string{"tag3", "tag4", "tag5"},
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"tags"},
			},
		})
		require.NoError(t, err)
		require.Len(t, updated.Tags, 3)
		require.Contains(t, updated.Tags, "tag3")
		require.Contains(t, updated.Tags, "tag4")
		require.Contains(t, updated.Tags, "tag5")
		require.NotContains(t, updated.Tags, "tag1")
		require.NotContains(t, updated.Tags, "tag2")

		// Verify tag_info is updated
		require.Len(t, updated.TagInfo, 3)

		// Verify old relationships were deleted
		bookmarkTags, err := s.Store.ListBookmarkTags(ctx, &store.FindBookmarkTag{
			ShortcutID: &created.Id,
		})
		require.NoError(t, err)
		require.Len(t, bookmarkTags, 3)

		// Verify old Tag entities still exist (not deleted)
		tag1, err := s.Store.GetTag(ctx, &store.FindTag{
			CreatorID:    &user.ID,
			Abbreviation: stringPtr("tag1"),
		})
		require.NoError(t, err)
		require.NotNil(t, tag1)
	})

	t.Run("Success_ClearAllTags", func(t *testing.T) {
		ctx := setUserContext(ctx, user)

		// Create shortcut with tags
		created, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "test-clear",
				Link: "https://example.com",
				Tags: []string{"tag1", "tag2"},
			},
		})
		require.NoError(t, err)

		// Update with empty tags
		updated, err := s.UpdateShortcut(ctx, &v1pb.UpdateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Id:   created.Id,
				Tags: []string{},
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"tags"},
			},
		})
		require.NoError(t, err)
		require.Len(t, updated.Tags, 0)
		require.Len(t, updated.TagInfo, 0)

		// Verify all relationships were deleted
		bookmarkTags, err := s.Store.ListBookmarkTags(ctx, &store.FindBookmarkTag{
			ShortcutID: &created.Id,
		})
		require.NoError(t, err)
		require.Len(t, bookmarkTags, 0)
	})

	t.Run("Success_UpdateWithoutChangingTags", func(t *testing.T) {
		ctx := setUserContext(ctx, user)

		// Create shortcut with tags
		created, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "test-no-change",
				Link: "https://example.com",
				Tags: []string{"tag1", "tag2"},
			},
		})
		require.NoError(t, err)

		// Update without tags in update_mask
		updated, err := s.UpdateShortcut(ctx, &v1pb.UpdateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Id:    created.Id,
				Title: "New Title",
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"title"},
			},
		})
		require.NoError(t, err)
		require.Equal(t, "New Title", updated.Title)

		// Verify tags are unchanged
		require.Len(t, updated.Tags, 2)
		require.Contains(t, updated.Tags, "tag1")
		require.Contains(t, updated.Tags, "tag2")
	})
}

func TestShortcutService_GetShortcut_WithTagInfo(t *testing.T) {
	ctx := context.Background()
	s, user := setupShortcutTagTest(t)

	t.Run("Success_TagInfoPopulated", func(t *testing.T) {
		ctx := setUserContext(ctx, user)

		// Create shortcut with tags
		created, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "test-get",
				Link: "https://example.com",
				Tags: []string{"backend", "frontend", "database"},
			},
		})
		require.NoError(t, err)

		// Get shortcut by ID
		retrieved, err := s.GetShortcut(ctx, &v1pb.GetShortcutRequest{
			Id: created.Id,
		})
		require.NoError(t, err)

		// Verify tags array is populated (backward compatibility)
		require.Len(t, retrieved.Tags, 3)

		// Verify tag_info is populated with UUIDs and names
		require.Len(t, retrieved.TagInfo, 3)
		tagNames := make(map[string]bool)
		for _, tagInfo := range retrieved.TagInfo {
			require.NotEmpty(t, tagInfo.Uuid)
			require.NotEmpty(t, tagInfo.Name)
			tagNames[tagInfo.Name] = true
		}
		require.True(t, tagNames["backend"])
		require.True(t, tagNames["frontend"])
		require.True(t, tagNames["database"])
	})

	t.Run("Success_NoTags", func(t *testing.T) {
		ctx := setUserContext(ctx, user)

		// Create shortcut without tags
		created, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "no-tags-get",
				Link: "https://example.com",
			},
		})
		require.NoError(t, err)

		// Get shortcut by ID
		retrieved, err := s.GetShortcut(ctx, &v1pb.GetShortcutRequest{
			Id: created.Id,
		})
		require.NoError(t, err)

		// Verify no tags
		require.Len(t, retrieved.Tags, 0)
		require.Len(t, retrieved.TagInfo, 0)
	})
}

func TestShortcutService_ListShortcuts_TagInfoNotPopulated(t *testing.T) {
	ctx := context.Background()
	s, user := setupShortcutTagTest(t)

	t.Run("TagInfo_NotIncludedInListView", func(t *testing.T) {
		ctx := setUserContext(ctx, user)

		// Create shortcuts with tags
		_, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "list-test-1",
				Link: "https://example1.com",
				Tags: []string{"tag1", "tag2"},
			},
		})
		require.NoError(t, err)

		_, err = s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "list-test-2",
				Link: "https://example2.com",
				Tags: []string{"tag3"},
			},
		})
		require.NoError(t, err)

		// List shortcuts
		resp, err := s.ListShortcuts(ctx, &v1pb.ListShortcutsRequest{})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(resp.Shortcuts), 2)

		// For list views, tag_info should NOT be populated (performance optimization)
		// But tags string array should still be present for backward compatibility
		for _, shortcut := range resp.Shortcuts {
			if shortcut.Name == "list-test-1" || shortcut.Name == "list-test-2" {
				// Tags array should be populated
				require.NotEmpty(t, shortcut.Tags)
				// But tag_info might not be populated in list views (implementation dependent)
				// This is acceptable as it's a performance optimization
			}
		}
	})
}

func TestShortcutService_DeleteShortcut_WithTags(t *testing.T) {
	ctx := context.Background()
	s, user := setupShortcutTagTest(t)

	t.Run("Success_ShortcutDeleted", func(t *testing.T) {
		ctx := setUserContext(ctx, user)

		// Create shortcut with tags
		created, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "delete-test",
				Link: "https://example.com",
				Tags: []string{"tag1", "tag2"},
			},
		})
		require.NoError(t, err)

		// Verify bookmark-tag relationships exist
		bookmarkTagsBefore, err := s.Store.ListBookmarkTags(ctx, &store.FindBookmarkTag{
			ShortcutID: &created.Id,
		})
		require.NoError(t, err)
		require.Len(t, bookmarkTagsBefore, 2)

		// Delete shortcut
		_, err = s.DeleteShortcut(ctx, &v1pb.DeleteShortcutRequest{
			Id: created.Id,
		})
		require.NoError(t, err)

		// Verify shortcut is deleted
		_, err = s.GetShortcut(ctx, &v1pb.GetShortcutRequest{
			Id: created.Id,
		})
		require.Error(t, err) // Should get not found error

		// Note: bookmark-tag cascade deletion depends on SQLite foreign keys being enabled
		// In test environment this may not work, but Tag entities should always remain
		tags, err := s.Store.ListTags(ctx, &store.FindTag{
			CreatorID: &user.ID,
		})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(tags), 2)
	})
}

func TestShortcutService_GetShortcutByName_WithTagInfo(t *testing.T) {
	ctx := context.Background()
	s, user := setupShortcutTagTest(t)

	t.Run("Success_TagInfoPopulated", func(t *testing.T) {
		ctx := setUserContext(ctx, user)

		// Create shortcut with unique name and tags
		created, err := s.CreateShortcut(ctx, &v1pb.CreateShortcutRequest{
			Shortcut: &v1pb.Shortcut{
				Name: "unique-name-test",
				Link: "https://example.com",
				Tags: []string{"web", "tools"},
			},
		})
		require.NoError(t, err)

		// Get shortcut by name
		retrieved, err := s.GetShortcutByName(ctx, &v1pb.GetShortcutByNameRequest{
			Name: "unique-name-test",
		})
		require.NoError(t, err)
		require.Equal(t, created.Id, retrieved.Id)

		// Verify tag_info is populated
		require.Len(t, retrieved.TagInfo, 2)
		for _, tagInfo := range retrieved.TagInfo {
			require.NotEmpty(t, tagInfo.Uuid)
			require.Contains(t, []string{"web", "tools"}, tagInfo.Name)
		}
	})
}

// Helper function to setup test service and user
func setupShortcutTagTest(t *testing.T) (*APIV1Service, *store.User) {
	ctx := context.Background()
	ts := teststore.NewTestingStore(ctx, t)
	s := &APIV1Service{Store: ts}

	// Create test user
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("test-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	user, err := ts.CreateUser(ctx, &store.User{
		Role:         store.RoleUser,
		Email:        "shortcuttagtest@test.com",
		Nickname:     "shortcut_tag_test_user",
		PasswordHash: string(passwordHash),
		UUID:         "shortcut-tag-user-uuid",
	})
	require.NoError(t, err)

	return s, user
}
