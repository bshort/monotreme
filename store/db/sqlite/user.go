package sqlite

import (
	"context"
	"errors"
	"strings"

	"github.com/bshort/monotreme/store"
)

func (d *DB) CreateUser(ctx context.Context, create *store.User) (*store.User, error) {
	stmt := `
		INSERT INTO user (
			email,
			nickname,
			password_hash,
			role,
			locale,
			color_theme,
			default_visibility,
			auto_generate_title,
			auto_generate_icon,
			auto_generate_name
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_ts, updated_ts, row_status
	`
	var rowStatus string
	if err := d.db.QueryRowContext(ctx, stmt,
		create.Email,
		create.Nickname,
		create.PasswordHash,
		create.Role,
		create.Locale,
		create.ColorTheme,
		create.DefaultVisibility,
		create.AutoGenerateTitle,
		create.AutoGenerateIcon,
		create.AutoGenerateName,
	).Scan(
		&create.ID,
		&create.CreatedTs,
		&create.UpdatedTs,
		&rowStatus,
	); err != nil {
		return nil, err
	}

	user := create
	user.RowStatus = store.ConvertRowStatusStringToStorepb(rowStatus)
	return user, nil
}

func (d *DB) UpdateUser(ctx context.Context, update *store.UpdateUser) (*store.User, error) {
	set, args := []string{}, []any{}
	if v := update.RowStatus; v != nil {
		set, args = append(set, "row_status = ?"), append(args, v.String())
	}
	if v := update.Email; v != nil {
		set, args = append(set, "email = ?"), append(args, *v)
	}
	if v := update.Nickname; v != nil {
		set, args = append(set, "nickname = ?"), append(args, *v)
	}
	if v := update.PasswordHash; v != nil {
		set, args = append(set, "password_hash = ?"), append(args, *v)
	}
	if v := update.Role; v != nil {
		set, args = append(set, "role = ?"), append(args, *v)
	}
	if v := update.Locale; v != nil {
		set, args = append(set, "locale = ?"), append(args, *v)
	}
	if v := update.ColorTheme; v != nil {
		set, args = append(set, "color_theme = ?"), append(args, *v)
	}
	if v := update.DefaultVisibility; v != nil {
		set, args = append(set, "default_visibility = ?"), append(args, *v)
	}
	if v := update.AutoGenerateTitle; v != nil {
		set, args = append(set, "auto_generate_title = ?"), append(args, *v)
	}
	if v := update.AutoGenerateIcon; v != nil {
		set, args = append(set, "auto_generate_icon = ?"), append(args, *v)
	}
	if v := update.AutoGenerateName; v != nil {
		set, args = append(set, "auto_generate_name = ?"), append(args, *v)
	}

	if len(set) == 0 {
		return nil, errors.New("no fields to update")
	}

	stmt := `
		UPDATE user
		SET ` + strings.Join(set, ", ") + `
		WHERE id = ?
		RETURNING id, created_ts, updated_ts, row_status, email, nickname, password_hash, role, locale, color_theme, default_visibility, auto_generate_title, auto_generate_icon, auto_generate_name
	`
	args = append(args, update.ID)
	user := &store.User{}
	var rowStatus string
	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&user.ID,
		&user.CreatedTs,
		&user.UpdatedTs,
		&rowStatus,
		&user.Email,
		&user.Nickname,
		&user.PasswordHash,
		&user.Role,
		&user.Locale,
		&user.ColorTheme,
		&user.DefaultVisibility,
		&user.AutoGenerateTitle,
		&user.AutoGenerateIcon,
		&user.AutoGenerateName,
	); err != nil {
		return nil, err
	}
	user.RowStatus = store.ConvertRowStatusStringToStorepb(rowStatus)
	return user, nil
}

func (d *DB) ListUsers(ctx context.Context, find *store.FindUser) ([]*store.User, error) {
	where, args := []string{"1 = 1"}, []any{}

	if v := find.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := find.RowStatus; v != nil {
		where, args = append(where, "row_status = ?"), append(args, v.String())
	}
	if v := find.Email; v != nil {
		where, args = append(where, "email = ?"), append(args, *v)
	}
	if v := find.Nickname; v != nil {
		where, args = append(where, "nickname = ?"), append(args, *v)
	}
	if v := find.Role; v != nil {
		where, args = append(where, "role = ?"), append(args, *v)
	}

	query := `
		SELECT
			id,
			created_ts,
			updated_ts,
			row_status,
			email,
			nickname,
			password_hash,
			role,
			locale,
			color_theme,
			default_visibility,
			auto_generate_title,
			auto_generate_icon,
			auto_generate_name
		FROM user
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY updated_ts DESC, created_ts DESC
	`
	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*store.User, 0)
	for rows.Next() {
		user := &store.User{}
		var rowStatus string
		if err := rows.Scan(
			&user.ID,
			&user.CreatedTs,
			&user.UpdatedTs,
			&rowStatus,
			&user.Email,
			&user.Nickname,
			&user.PasswordHash,
			&user.Role,
			&user.Locale,
			&user.ColorTheme,
			&user.DefaultVisibility,
			&user.AutoGenerateTitle,
			&user.AutoGenerateIcon,
			&user.AutoGenerateName,
		); err != nil {
			return nil, err
		}
		user.RowStatus = store.ConvertRowStatusStringToStorepb(rowStatus)
		list = append(list, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (d *DB) DeleteUser(ctx context.Context, delete *store.DeleteUser) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM user WHERE id = ?
	`, delete.ID); err != nil {
		return err
	}

	if err := vacuumUserSetting(ctx, tx); err != nil {
		return err
	}
	if err := vacuumShortcut(ctx, tx); err != nil {
		return err
	}
	if err := vacuumCollection(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}
