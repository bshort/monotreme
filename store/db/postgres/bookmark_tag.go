package postgres

import (
	"context"

	"github.com/pkg/errors"

	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

func (d *DB) CreateBookmarkTag(ctx context.Context, create *store.CreateBookmarkTag) (*storepb.BookmarkTag, error) {
	return nil, errors.New("not implemented for postgres")
}

func (d *DB) ListBookmarkTags(ctx context.Context, find *store.FindBookmarkTag) ([]*storepb.BookmarkTag, error) {
	return nil, errors.New("not implemented for postgres")
}

func (d *DB) DeleteBookmarkTag(ctx context.Context, delete *store.DeleteBookmarkTag) error {
	return errors.New("not implemented for postgres")
}

func (d *DB) DeleteBookmarkTagsByShortcut(ctx context.Context, shortcutID int32) error {
	return errors.New("not implemented for postgres")
}

func (d *DB) DeleteBookmarkTagsByTag(ctx context.Context, tagUUID string) error {
	return errors.New("not implemented for postgres")
}
