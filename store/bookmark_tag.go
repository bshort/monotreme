package store

import (
	"context"

	storepb "github.com/bshort/monotreme/proto/gen/store"
)

type CreateBookmarkTag struct {
	ShortcutID int32
	TagUUID    string
	UserID     int32
}

type FindBookmarkTag struct {
	ShortcutID *int32
	TagUUID    *string
	UserID     *int32
}

type DeleteBookmarkTag struct {
	ShortcutID int32
	TagUUID    string
}

func (s *Store) CreateBookmarkTag(ctx context.Context, create *CreateBookmarkTag) (*storepb.BookmarkTag, error) {
	return s.driver.CreateBookmarkTag(ctx, create)
}

func (s *Store) ListBookmarkTags(ctx context.Context, find *FindBookmarkTag) ([]*storepb.BookmarkTag, error) {
	return s.driver.ListBookmarkTags(ctx, find)
}

func (s *Store) DeleteBookmarkTag(ctx context.Context, delete *DeleteBookmarkTag) error {
	return s.driver.DeleteBookmarkTag(ctx, delete)
}

// DeleteBookmarkTagsByShortcut deletes all bookmark-tag relationships for a given shortcut
func (s *Store) DeleteBookmarkTagsByShortcut(ctx context.Context, shortcutID int32) error {
	return s.driver.DeleteBookmarkTagsByShortcut(ctx, shortcutID)
}

// DeleteBookmarkTagsByTag deletes all bookmark-tag relationships for a given tag
func (s *Store) DeleteBookmarkTagsByTag(ctx context.Context, tagUUID string) error {
	return s.driver.DeleteBookmarkTagsByTag(ctx, tagUUID)
}
