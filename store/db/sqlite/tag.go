package sqlite

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

func (d *DB) CreateTag(ctx context.Context, create *storepb.Tag) (*storepb.Tag, error) {
	set := []string{"uuid", "creator_id", "name", "abbreviation", "description"}
	args := []any{create.Uuid, create.CreatorId, create.Name, create.Abbreviation, create.Description}
	placeholder := []string{"?", "?", "?", "?", "?"}

	stmt := `
		INSERT INTO tag (
			` + strings.Join(set, ", ") + `
		)
		VALUES (` + strings.Join(placeholder, ",") + `)
		RETURNING created_ts, updated_ts
	`
	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&create.CreatedTs,
		&create.UpdatedTs,
	); err != nil {
		return nil, err
	}
	tag := create
	return tag, nil
}

func (d *DB) UpdateTag(ctx context.Context, update *store.UpdateTag) (*storepb.Tag, error) {
	set, args := []string{}, []any{}
	if update.Name != nil {
		set, args = append(set, "name = ?"), append(args, *update.Name)
	}
	if update.Abbreviation != nil {
		set, args = append(set, "abbreviation = ?"), append(args, *update.Abbreviation)
	}
	if update.Description != nil {
		set, args = append(set, "description = ?"), append(args, *update.Description)
	}
	if len(set) == 0 {
		return nil, errors.New("no update specified")
	}
	args = append(args, update.UUID)

	stmt := `
		UPDATE tag
		SET
			` + strings.Join(set, ", ") + `
		WHERE
			uuid = ?
		RETURNING uuid, creator_id, created_ts, updated_ts, name, abbreviation, description
	`
	tag := &storepb.Tag{}
	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&tag.Uuid,
		&tag.CreatorId,
		&tag.CreatedTs,
		&tag.UpdatedTs,
		&tag.Name,
		&tag.Abbreviation,
		&tag.Description,
	); err != nil {
		return nil, err
	}

	return tag, nil
}

func (d *DB) ListTags(ctx context.Context, find *store.FindTag) ([]*storepb.Tag, error) {
	where, args := []string{"1 = 1"}, []any{}
	if v := find.UUID; v != nil {
		where, args = append(where, "uuid = ?"), append(args, *v)
	}
	if v := find.CreatorID; v != nil {
		where, args = append(where, "creator_id = ?"), append(args, *v)
	}
	if v := find.Abbreviation; v != nil {
		where, args = append(where, "abbreviation = ?"), append(args, *v)
	}

	rows, err := d.db.QueryContext(ctx, `
		SELECT
			uuid,
			creator_id,
			created_ts,
			updated_ts,
			name,
			abbreviation,
			description
		FROM tag
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY created_ts DESC`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*storepb.Tag, 0)
	for rows.Next() {
		tag := &storepb.Tag{}
		if err := rows.Scan(
			&tag.Uuid,
			&tag.CreatorId,
			&tag.CreatedTs,
			&tag.UpdatedTs,
			&tag.Name,
			&tag.Abbreviation,
			&tag.Description,
		); err != nil {
			return nil, err
		}

		list = append(list, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (d *DB) DeleteTag(ctx context.Context, delete *store.DeleteTag) error {
	if _, err := d.db.ExecContext(ctx, `DELETE FROM tag WHERE uuid = ?`, delete.UUID); err != nil {
		return err
	}

	return nil
}
