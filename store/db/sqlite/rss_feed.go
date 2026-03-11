package sqlite

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

func (d *DB) CreateRssFeed(ctx context.Context, create *store.RssFeed) (*store.RssFeed, error) {
	if create.UUID == "" {
		create.UUID = uuid.New().String()
	}

	defaultTags := strings.Join(create.DefaultTags, " ")

	stmt := `
		INSERT INTO rss_feed (
			uuid, creator_id, title, url, description, auto_import,
			import_frequency_hours, default_tags, default_visibility,
			shortcut_prefix, is_active, last_error, total_imported
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_ts, updated_ts
	`

	var createdTs, updatedTs int64
	if err := d.db.QueryRowContext(ctx, stmt,
		create.UUID,
		create.CreatorID,
		create.Title,
		create.URL,
		create.Description,
		create.AutoImport,
		create.ImportFrequencyHours,
		defaultTags,
		create.DefaultVisibility.String(),
		create.ShortcutPrefix,
		create.IsActive,
		create.LastError,
		create.TotalImported,
	).Scan(&create.ID, &createdTs, &updatedTs); err != nil {
		return nil, err
	}

	create.CreatedTs = timestampFromSeconds(createdTs)
	create.UpdatedTs = timestampFromSeconds(updatedTs)
	create.RowStatus = storepb.RowStatus_NORMAL

	return create, nil
}

func (d *DB) UpdateRssFeed(ctx context.Context, update *store.UpdateRssFeed) (*store.RssFeed, error) {
	set, args := []string{}, []any{}

	if update.Title != nil {
		set, args = append(set, "title = ?"), append(args, *update.Title)
	}
	if update.URL != nil {
		set, args = append(set, "url = ?"), append(args, *update.URL)
	}
	if update.Description != nil {
		set, args = append(set, "description = ?"), append(args, *update.Description)
	}
	if update.AutoImport != nil {
		set, args = append(set, "auto_import = ?"), append(args, *update.AutoImport)
	}
	if update.ImportFrequencyHours != nil {
		set, args = append(set, "import_frequency_hours = ?"), append(args, *update.ImportFrequencyHours)
	}
	if update.LastImportTime != nil {
		set, args = append(set, "last_import_ts = ?"), append(args, update.LastImportTime.Unix())
	}
	if update.DefaultTags != nil {
		defaultTags := strings.Join(update.DefaultTags, " ")
		set, args = append(set, "default_tags = ?"), append(args, defaultTags)
	}
	if update.DefaultVisibility != nil {
		set, args = append(set, "default_visibility = ?"), append(args, update.DefaultVisibility.String())
	}
	if update.ShortcutPrefix != nil {
		set, args = append(set, "shortcut_prefix = ?"), append(args, *update.ShortcutPrefix)
	}
	if update.IsActive != nil {
		set, args = append(set, "is_active = ?"), append(args, *update.IsActive)
	}
	if update.LastError != nil {
		set, args = append(set, "last_error = ?"), append(args, *update.LastError)
	}
	if update.TotalImported != nil {
		set, args = append(set, "total_imported = ?"), append(args, *update.TotalImported)
	}

	if len(set) == 0 {
		return nil, errors.New("no update specified")
	}

	set = append(set, "updated_ts = ?")
	args = append(args, time.Now().Unix())
	args = append(args, update.ID)

	stmt := `
		UPDATE rss_feed
		SET ` + strings.Join(set, ", ") + `
		WHERE id = ?
		RETURNING id, uuid, creator_id, created_ts, updated_ts, title, url, description,
		          auto_import, import_frequency_hours, last_import_ts, default_tags,
		          default_visibility, shortcut_prefix, is_active, last_error, total_imported
	`

	rssFeed := &store.RssFeed{}
	var createdTs, updatedTs int64
	var lastImportTs sql.NullInt64
	var defaultTags, defaultVisibility string

	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&rssFeed.ID,
		&rssFeed.UUID,
		&rssFeed.CreatorID,
		&createdTs,
		&updatedTs,
		&rssFeed.Title,
		&rssFeed.URL,
		&rssFeed.Description,
		&rssFeed.AutoImport,
		&rssFeed.ImportFrequencyHours,
		&lastImportTs,
		&defaultTags,
		&defaultVisibility,
		&rssFeed.ShortcutPrefix,
		&rssFeed.IsActive,
		&rssFeed.LastError,
		&rssFeed.TotalImported,
	); err != nil {
		return nil, err
	}

	rssFeed.CreatedTs = timestampFromSeconds(createdTs)
	rssFeed.UpdatedTs = timestampFromSeconds(updatedTs)
	if lastImportTs.Valid {
		rssFeed.LastImportTs = timestampFromSeconds(lastImportTs.Int64)
	}
	rssFeed.DefaultTags = filterTags(strings.Split(defaultTags, " "))
	rssFeed.DefaultVisibility = store.ConvertVisibilityStringToStorepb(defaultVisibility)
	rssFeed.RowStatus = storepb.RowStatus_NORMAL

	return rssFeed, nil
}

func (d *DB) ListRssFeeds(ctx context.Context, find *store.FindRssFeed) ([]*store.RssFeed, error) {
	where, args := []string{"row_status = 'NORMAL'"}, []any{}

	if v := find.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := find.CreatorID; v != nil {
		where, args = append(where, "creator_id = ?"), append(args, *v)
	}
	if v := find.IsActive; v != nil {
		where, args = append(where, "is_active = ?"), append(args, *v)
	}

	stmt := `
		SELECT id, uuid, creator_id, created_ts, updated_ts, title, url, description,
		       auto_import, import_frequency_hours, last_import_ts, default_tags,
		       default_visibility, shortcut_prefix, is_active, last_error, total_imported
		FROM rss_feed
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY created_ts DESC
	`

	rows, err := d.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.RssFeed, 0)
	for rows.Next() {
		rssFeed := &store.RssFeed{}
		var createdTs, updatedTs int64
		var lastImportTs sql.NullInt64
		var defaultTags, defaultVisibility string

		if err := rows.Scan(
			&rssFeed.ID,
			&rssFeed.UUID,
			&rssFeed.CreatorID,
			&createdTs,
			&updatedTs,
			&rssFeed.Title,
			&rssFeed.URL,
			&rssFeed.Description,
			&rssFeed.AutoImport,
			&rssFeed.ImportFrequencyHours,
			&lastImportTs,
			&defaultTags,
			&defaultVisibility,
			&rssFeed.ShortcutPrefix,
			&rssFeed.IsActive,
			&rssFeed.LastError,
			&rssFeed.TotalImported,
		); err != nil {
			return nil, err
		}

		rssFeed.CreatedTs = timestampFromSeconds(createdTs)
		rssFeed.UpdatedTs = timestampFromSeconds(updatedTs)
		if lastImportTs.Valid {
			rssFeed.LastImportTs = timestampFromSeconds(lastImportTs.Int64)
		}
		rssFeed.DefaultTags = filterTags(strings.Split(defaultTags, " "))
		rssFeed.DefaultVisibility = store.ConvertVisibilityStringToStorepb(defaultVisibility)
		rssFeed.RowStatus = storepb.RowStatus_NORMAL

		list = append(list, rssFeed)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *DB) DeleteRssFeed(ctx context.Context, delete *store.DeleteRssFeed) error {
	stmt := `UPDATE rss_feed SET row_status = 'ARCHIVED' WHERE id = ?`
	if _, err := d.db.ExecContext(ctx, stmt, delete.ID); err != nil {
		return err
	}
	return nil
}

func (d *DB) CreateRssFeedItem(ctx context.Context, create *store.RssFeedItem) (*store.RssFeedItem, error) {
	if create.UUID == "" {
		create.UUID = uuid.New().String()
	}

	stmt := `
		INSERT INTO rss_feed_item (
			uuid, rss_feed_id, item_guid, item_link, item_title,
			item_description, item_published_ts, shortcut_id,
			import_success, import_error
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_ts
	`

	var createdTs int64
	var itemPublishedTs sql.NullInt64
	if create.ItemPublishedTs != nil {
		itemPublishedTs.Valid = true
		itemPublishedTs.Int64 = create.ItemPublishedTs.Unix()
	}

	if err := d.db.QueryRowContext(ctx, stmt,
		create.UUID,
		create.RssFeedID,
		create.ItemGUID,
		create.ItemLink,
		create.ItemTitle,
		create.ItemDescription,
		itemPublishedTs,
		create.ShortcutID,
		create.ImportSuccess,
		create.ImportError,
	).Scan(&create.ID, &createdTs); err != nil {
		return nil, err
	}

	create.CreatedTs = timestampFromSeconds(createdTs)
	create.RowStatus = storepb.RowStatus_NORMAL

	return create, nil
}

func (d *DB) GetRssFeedItem(ctx context.Context, find *store.FindRssFeedItem) (*store.RssFeedItem, error) {
	where, args := []string{"row_status = 'NORMAL'"}, []any{}

	if v := find.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := find.RssFeedID; v != nil {
		where, args = append(where, "rss_feed_id = ?"), append(args, *v)
	}
	if v := find.ItemGUID; v != nil {
		where, args = append(where, "item_guid = ?"), append(args, *v)
	}
	if v := find.ItemLink; v != nil {
		where, args = append(where, "item_link = ?"), append(args, *v)
	}

	stmt := `
		SELECT id, uuid, rss_feed_id, created_ts, item_guid, item_link,
		       item_title, item_description, item_published_ts, shortcut_id,
		       import_success, import_error
		FROM rss_feed_item
		WHERE ` + strings.Join(where, " AND ") + `
		LIMIT 1
	`

	item := &store.RssFeedItem{}
	var createdTs int64
	var itemPublishedTs, shortcutID sql.NullInt64

	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&item.ID,
		&item.UUID,
		&item.RssFeedID,
		&createdTs,
		&item.ItemGUID,
		&item.ItemLink,
		&item.ItemTitle,
		&item.ItemDescription,
		&itemPublishedTs,
		&shortcutID,
		&item.ImportSuccess,
		&item.ImportError,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	item.CreatedTs = timestampFromSeconds(createdTs)
	if itemPublishedTs.Valid {
		item.ItemPublishedTs = timestampFromSeconds(itemPublishedTs.Int64)
	}
	if shortcutID.Valid {
		id := int32(shortcutID.Int64)
		item.ShortcutID = &id
	}
	item.RowStatus = storepb.RowStatus_NORMAL

	return item, nil
}

func (d *DB) ListRssFeedItems(ctx context.Context, find *store.FindRssFeedItem) ([]*store.RssFeedItem, error) {
	where, args := []string{"row_status = 'NORMAL'"}, []any{}

	if v := find.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := find.RssFeedID; v != nil {
		where, args = append(where, "rss_feed_id = ?"), append(args, *v)
	}
	if v := find.ItemGUID; v != nil {
		where, args = append(where, "item_guid = ?"), append(args, *v)
	}
	if v := find.ItemLink; v != nil {
		where, args = append(where, "item_link = ?"), append(args, *v)
	}

	stmt := `
		SELECT id, uuid, rss_feed_id, created_ts, item_guid, item_link,
		       item_title, item_description, item_published_ts, shortcut_id,
		       import_success, import_error
		FROM rss_feed_item
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY created_ts DESC
	`

	rows, err := d.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.RssFeedItem, 0)
	for rows.Next() {
		item := &store.RssFeedItem{}
		var createdTs int64
		var itemPublishedTs, shortcutID sql.NullInt64

		if err := rows.Scan(
			&item.ID,
			&item.UUID,
			&item.RssFeedID,
			&createdTs,
			&item.ItemGUID,
			&item.ItemLink,
			&item.ItemTitle,
			&item.ItemDescription,
			&itemPublishedTs,
			&shortcutID,
			&item.ImportSuccess,
			&item.ImportError,
		); err != nil {
			return nil, err
		}

		item.CreatedTs = timestampFromSeconds(createdTs)
		if itemPublishedTs.Valid {
			item.ItemPublishedTs = timestampFromSeconds(itemPublishedTs.Int64)
		}
		if shortcutID.Valid {
			id := int32(shortcutID.Int64)
			item.ShortcutID = &id
		}
		item.RowStatus = storepb.RowStatus_NORMAL

		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func timestampFromSeconds(seconds int64) *time.Time {
	t := time.Unix(seconds, 0)
	return &t
}