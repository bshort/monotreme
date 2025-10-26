package sqlite

import (
	"context"
	"strings"

	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

func (d *DB) CreateBookmarkTag(ctx context.Context, create *store.CreateBookmarkTag) (*storepb.BookmarkTag, error) {
	stmt := `
		INSERT INTO bookmark_tag (
			shortcut_id,
			tag_uuid,
			user_id
		)
		VALUES (?, ?, ?)
		ON CONFLICT(shortcut_id, tag_uuid) DO NOTHING
		RETURNING shortcut_id, tag_uuid, user_id, created_ts
	`

	bookmarkTag := &storepb.BookmarkTag{}
	if err := d.db.QueryRowContext(ctx, stmt,
		create.ShortcutID,
		create.TagUUID,
		create.UserID,
	).Scan(
		&bookmarkTag.ShortcutId,
		&bookmarkTag.TagUuid,
		&bookmarkTag.UserId,
		&bookmarkTag.CreatedTs,
	); err != nil {
		return nil, err
	}

	return bookmarkTag, nil
}

func (d *DB) ListBookmarkTags(ctx context.Context, find *store.FindBookmarkTag) ([]*storepb.BookmarkTag, error) {
	where, args := []string{"1 = 1"}, []any{}

	if v := find.ShortcutID; v != nil {
		where, args = append(where, "shortcut_id = ?"), append(args, *v)
	}
	if v := find.TagUUID; v != nil {
		where, args = append(where, "tag_uuid = ?"), append(args, *v)
	}
	if v := find.UserID; v != nil {
		where, args = append(where, "user_id = ?"), append(args, *v)
	}

	rows, err := d.db.QueryContext(ctx, `
		SELECT
			shortcut_id,
			tag_uuid,
			user_id,
			created_ts
		FROM bookmark_tag
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY created_ts DESC`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*storepb.BookmarkTag, 0)
	for rows.Next() {
		bookmarkTag := &storepb.BookmarkTag{}
		if err := rows.Scan(
			&bookmarkTag.ShortcutId,
			&bookmarkTag.TagUuid,
			&bookmarkTag.UserId,
			&bookmarkTag.CreatedTs,
		); err != nil {
			return nil, err
		}

		list = append(list, bookmarkTag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (d *DB) DeleteBookmarkTag(ctx context.Context, delete *store.DeleteBookmarkTag) error {
	stmt := `DELETE FROM bookmark_tag WHERE shortcut_id = ? AND tag_uuid = ?`
	if _, err := d.db.ExecContext(ctx, stmt, delete.ShortcutID, delete.TagUUID); err != nil {
		return err
	}
	return nil
}

func (d *DB) DeleteBookmarkTagsByShortcut(ctx context.Context, shortcutID int32) error {
	stmt := `DELETE FROM bookmark_tag WHERE shortcut_id = ?`
	if _, err := d.db.ExecContext(ctx, stmt, shortcutID); err != nil {
		return err
	}
	return nil
}

func (d *DB) DeleteBookmarkTagsByTag(ctx context.Context, tagUUID string) error {
	stmt := `DELETE FROM bookmark_tag WHERE tag_uuid = ?`
	if _, err := d.db.ExecContext(ctx, stmt, tagUUID); err != nil {
		return err
	}
	return nil
}
