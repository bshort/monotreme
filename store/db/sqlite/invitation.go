package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

func (d *DB) CreateInvitation(ctx context.Context, create *storepb.Invitation) (*storepb.Invitation, error) {
	stmt := `
		INSERT INTO invitation (
			from_path,
			to_path,
			accepted_at,
			deleted_at
		)
		VALUES (?, ?, ?, ?)
		RETURNING id, created_ts, updated_ts
	`

	// Handle NULL values for optional timestamp fields
	var acceptedAt, deletedAt interface{}
	if create.AcceptedAt != "" {
		acceptedAt = create.AcceptedAt
	}
	if create.DeletedAt != "" {
		deletedAt = create.DeletedAt
	}

	if err := d.db.QueryRowContext(ctx, stmt,
		create.From,
		create.To,
		acceptedAt,
		deletedAt,
	).Scan(
		&create.Id,
		&create.CreatedTs,
		&create.UpdatedTs,
	); err != nil {
		return nil, err
	}

	invitation := create
	return invitation, nil
}

func (d *DB) UpdateInvitation(ctx context.Context, update *store.UpdateInvitation) (*storepb.Invitation, error) {
	set, args := []string{}, []any{}
	if update.From != nil {
		set, args = append(set, "from_path = ?"), append(args, *update.From)
	}
	if update.To != nil {
		set, args = append(set, "to_path = ?"), append(args, *update.To)
	}
	if update.AcceptedAt != nil {
		var acceptedAt interface{}
		if *update.AcceptedAt != "" {
			acceptedAt = *update.AcceptedAt
		}
		set, args = append(set, "accepted_at = ?"), append(args, acceptedAt)
	}
	if update.DeletedAt != nil {
		var deletedAt interface{}
		if *update.DeletedAt != "" {
			deletedAt = *update.DeletedAt
		}
		set, args = append(set, "deleted_at = ?"), append(args, deletedAt)
	}
	if len(set) == 0 {
		return nil, errors.New("no fields to update")
	}
	args = append(args, update.ID)

	stmt := `
		UPDATE invitation
		SET
			` + strings.Join(set, ", ") + `
		WHERE
			id = ?
		RETURNING id, created_ts, updated_ts, from_path, to_path, accepted_at, deleted_at
	`
	invitation := &storepb.Invitation{}
	var acceptedAt, deletedAt sql.NullString
	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&invitation.Id,
		&invitation.CreatedTs,
		&invitation.UpdatedTs,
		&invitation.From,
		&invitation.To,
		&acceptedAt,
		&deletedAt,
	); err != nil {
		return nil, err
	}

	if acceptedAt.Valid {
		invitation.AcceptedAt = acceptedAt.String
	}
	if deletedAt.Valid {
		invitation.DeletedAt = deletedAt.String
	}

	return invitation, nil
}

func (d *DB) ListInvitations(ctx context.Context, find *store.FindInvitation) ([]*storepb.Invitation, error) {
	where, args := []string{"1 = 1"}, []any{}

	if v := find.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := find.From; v != nil {
		where, args = append(where, "from_path = ?"), append(args, *v)
	}
	if v := find.To; v != nil {
		where, args = append(where, "to_path = ?"), append(args, *v)
	}
	if v := find.Accepted; v != nil {
		if *v {
			where = append(where, "accepted_at IS NOT NULL")
		} else {
			where = append(where, "accepted_at IS NULL")
		}
	}
	if v := find.Deleted; v != nil {
		if *v {
			where = append(where, "deleted_at IS NOT NULL")
		} else {
			where = append(where, "deleted_at IS NULL")
		}
	}

	rows, err := d.db.QueryContext(ctx, `
		SELECT
			id,
			created_ts,
			updated_ts,
			from_path,
			to_path,
			accepted_at,
			deleted_at
		FROM invitation
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY created_ts DESC`,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*storepb.Invitation, 0)
	for rows.Next() {
		invitation := &storepb.Invitation{}
		var acceptedAt, deletedAt sql.NullString
		if err := rows.Scan(
			&invitation.Id,
			&invitation.CreatedTs,
			&invitation.UpdatedTs,
			&invitation.From,
			&invitation.To,
			&acceptedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}

		if acceptedAt.Valid {
			invitation.AcceptedAt = acceptedAt.String
		}
		if deletedAt.Valid {
			invitation.DeletedAt = deletedAt.String
		}

		list = append(list, invitation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (d *DB) DeleteInvitation(ctx context.Context, delete *store.DeleteInvitation) error {
	if _, err := d.db.ExecContext(ctx, `DELETE FROM invitation WHERE id = ?`, delete.ID); err != nil {
		return err
	}

	return nil
}