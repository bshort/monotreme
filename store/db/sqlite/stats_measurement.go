package sqlite

import (
	"context"
	"strings"

	"github.com/bshort/monotreme/store"
)

func (d *DB) CreateStatsMeasurement(ctx context.Context, create *store.StatsMeasurement) (*store.StatsMeasurement, error) {
	stmt := `
		INSERT INTO stats_measurement (
			measured_ts,
			shortcuts_count,
			users_count,
			collections_count,
			hits_count
		)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id
	`
	if err := d.db.QueryRowContext(ctx, stmt,
		create.MeasuredTs,
		create.ShortcutsCount,
		create.UsersCount,
		create.CollectionsCount,
		create.HitsCount,
	).Scan(
		&create.ID,
	); err != nil {
		return nil, err
	}

	return create, nil
}

func (d *DB) ListStatsMeasurements(ctx context.Context, find *store.FindStatsMeasurement) ([]*store.StatsMeasurement, error) {
	where, args := []string{"1 = 1"}, []any{}
	if find.ID != nil {
		where, args = append(where, "id = ?"), append(args, *find.ID)
	}

	order := ""
	if find.OrderByMeasuredTsDesc != nil && *find.OrderByMeasuredTsDesc {
		order = "ORDER BY measured_ts DESC"
	} else {
		order = "ORDER BY measured_ts ASC"
	}

	limit := ""
	if find.Limit != nil {
		limit = "LIMIT ?"
		args = append(args, *find.Limit)
	}

	stmt := `
		SELECT
			id,
			measured_ts,
			shortcuts_count,
			users_count,
			collections_count,
			hits_count
		FROM stats_measurement
		WHERE ` + strings.Join(where, " AND ") + `
		` + order + `
		` + limit + `
	`

	rows, err := d.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statsMeasurements := []*store.StatsMeasurement{}
	for rows.Next() {
		statsMeasurement := &store.StatsMeasurement{}
		if err := rows.Scan(
			&statsMeasurement.ID,
			&statsMeasurement.MeasuredTs,
			&statsMeasurement.ShortcutsCount,
			&statsMeasurement.UsersCount,
			&statsMeasurement.CollectionsCount,
			&statsMeasurement.HitsCount,
		); err != nil {
			return nil, err
		}
		statsMeasurements = append(statsMeasurements, statsMeasurement)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statsMeasurements, nil
}

func (d *DB) UpdateStatsMeasurement(ctx context.Context, update *store.UpdateStatsMeasurement) (*store.StatsMeasurement, error) {
	set, args := []string{}, []any{}
	if update.ShortcutsCount != nil {
		set, args = append(set, "shortcuts_count = ?"), append(args, *update.ShortcutsCount)
	}
	if update.UsersCount != nil {
		set, args = append(set, "users_count = ?"), append(args, *update.UsersCount)
	}
	if update.CollectionsCount != nil {
		set, args = append(set, "collections_count = ?"), append(args, *update.CollectionsCount)
	}
	if update.HitsCount != nil {
		set, args = append(set, "hits_count = ?"), append(args, *update.HitsCount)
	}

	args = append(args, update.ID)
	stmt := `UPDATE stats_measurement SET ` + strings.Join(set, ", ") + ` WHERE id = ? RETURNING id, measured_ts, shortcuts_count, users_count, collections_count, hits_count`
	statsMeasurement := &store.StatsMeasurement{}
	if err := d.db.QueryRowContext(ctx, stmt, args...).Scan(
		&statsMeasurement.ID,
		&statsMeasurement.MeasuredTs,
		&statsMeasurement.ShortcutsCount,
		&statsMeasurement.UsersCount,
		&statsMeasurement.CollectionsCount,
		&statsMeasurement.HitsCount,
	); err != nil {
		return nil, err
	}

	return statsMeasurement, nil
}

func (d *DB) DeleteStatsMeasurement(ctx context.Context, delete *store.DeleteStatsMeasurement) error {
	stmt := `DELETE FROM stats_measurement WHERE id = ?`
	result, err := d.db.ExecContext(ctx, stmt, delete.ID)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}