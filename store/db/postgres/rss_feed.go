package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

func (d *DB) CreateRssFeed(ctx context.Context, create *store.RssFeed) (*store.RssFeed, error) {
	if create.UUID == "" {
		create.UUID = uuid.New().String()
	}

	stmt := `
		INSERT INTO rss_feed (
			uuid, creator_id, title, url, description, auto_import,
			import_frequency_hours, default_tags, default_visibility,
			shortcut_prefix, is_active, last_error, total_imported
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
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
		pq.Array(create.DefaultTags),
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
	argIndex := 1

	if update.Title != nil {
		set = append(set, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *update.Title)
		argIndex++
	}
	if update.URL != nil {
		set = append(set, fmt.Sprintf("url = $%d", argIndex))
		args = append(args, *update.URL)
		argIndex++
	}
	if update.Description != nil {
		set = append(set, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *update.Description)
		argIndex++
	}
	if update.AutoImport != nil {
		set = append(set, fmt.Sprintf("auto_import = $%d", argIndex))
		args = append(args, *update.AutoImport)
		argIndex++
	}
	if update.ImportFrequencyHours != nil {
		set = append(set, fmt.Sprintf("import_frequency_hours = $%d", argIndex))
		args = append(args, *update.ImportFrequencyHours)
		argIndex++
	}
	if update.LastImportTime != nil {
		set = append(set, fmt.Sprintf("last_import_ts = $%d", argIndex))
		args = append(args, update.LastImportTime.Unix())
		argIndex++
	}
	if update.DefaultTags != nil {
		set = append(set, fmt.Sprintf("default_tags = $%d", argIndex))
		args = append(args, pq.Array(update.DefaultTags))
		argIndex++
	}
	if update.DefaultVisibility != nil {
		set = append(set, fmt.Sprintf("default_visibility = $%d", argIndex))
		args = append(args, update.DefaultVisibility.String())
		argIndex++
	}
	if update.ShortcutPrefix != nil {
		set = append(set, fmt.Sprintf("shortcut_prefix = $%d", argIndex))
		args = append(args, *update.ShortcutPrefix)
		argIndex++
	}
	if update.IsActive != nil {
		set = append(set, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *update.IsActive)
		argIndex++
	}
	if update.LastError != nil {
		set = append(set, fmt.Sprintf("last_error = $%d", argIndex))
		args = append(args, *update.LastError)
		argIndex++
	}
	if update.TotalImported != nil {
		set = append(set, fmt.Sprintf("total_imported = $%d", argIndex))
		args = append(args, *update.TotalImported)
		argIndex++
	}

	if len(set) == 0 {
		return nil, errors.New("no update specified")
	}

	set = append(set, fmt.Sprintf("updated_ts = $%d", argIndex))
	args = append(args, time.Now().Unix())
	argIndex++

	args = append(args, update.ID)

	stmt := `
		UPDATE rss_feed
		SET ` + strings.Join(set, ", ") + `
		WHERE id = $` + fmt.Sprintf("%d", argIndex) + `
		RETURNING id, uuid, creator_id, created_ts, updated_ts, title, url, description,
		          auto_import, import_frequency_hours, last_import_ts, default_tags,
		          default_visibility, shortcut_prefix, is_active, last_error, total_imported
	`

	rssFeed := &store.RssFeed{}
	var createdTs, updatedTs int64
	var lastImportTs sql.NullInt64
	var defaultTags pq.StringArray
	var defaultVisibility string

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
	rssFeed.DefaultTags = []string(defaultTags)
	rssFeed.DefaultVisibility = store.ConvertVisibilityStringToStorepb(defaultVisibility)
	rssFeed.RowStatus = storepb.RowStatus_NORMAL

	return rssFeed, nil
}

func (d *DB) ListRssFeeds(ctx context.Context, find *store.FindRssFeed) ([]*store.RssFeed, error) {
	where, args := []string{"row_status = 'NORMAL'"}, []any{}
	argIndex := 1

	if v := find.ID; v != nil {
		where = append(where, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := find.CreatorID; v != nil {
		where = append(where, fmt.Sprintf("creator_id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := find.IsActive; v != nil {
		where = append(where, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *v)
		argIndex++
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
		var defaultTags pq.StringArray
		var defaultVisibility string

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
		rssFeed.DefaultTags = []string(defaultTags)
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
	stmt := `UPDATE rss_feed SET row_status = 'ARCHIVED' WHERE id = $1`
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
	argIndex := 1

	if v := find.ID; v != nil {
		where = append(where, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := find.RssFeedID; v != nil {
		where = append(where, fmt.Sprintf("rss_feed_id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := find.ItemGUID; v != nil {
		where = append(where, fmt.Sprintf("item_guid = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := find.ItemLink; v != nil {
		where = append(where, fmt.Sprintf("item_link = $%d", argIndex))
		args = append(args, *v)
		argIndex++
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
	argIndex := 1

	if v := find.ID; v != nil {
		where = append(where, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := find.RssFeedID; v != nil {
		where = append(where, fmt.Sprintf("rss_feed_id = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := find.ItemGUID; v != nil {
		where = append(where, fmt.Sprintf("item_guid = $%d", argIndex))
		args = append(args, *v)
		argIndex++
	}
	if v := find.ItemLink; v != nil {
		where = append(where, fmt.Sprintf("item_link = $%d", argIndex))
		args = append(args, *v)
		argIndex++
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