package postgres

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
		VALUES ($1, $2, $3, $4, $5)
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
	argIndex := 0
	if find.ID != nil {
		argIndex++
		where, args = append(where, "id = "+placeholder(argIndex)), append(args, *find.ID)
	}

	order := ""
	if find.OrderByMeasuredTsDesc != nil && *find.OrderByMeasuredTsDesc {
		order = "ORDER BY measured_ts DESC"
	} else {
		order = "ORDER BY measured_ts ASC"
	}

	limit := ""
	if find.Limit != nil {
		argIndex++
		limit = "LIMIT " + placeholder(argIndex)
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
	argIndex := 0

	if update.ShortcutsCount != nil {
		argIndex++
		set, args = append(set, "shortcuts_count = "+placeholder(argIndex)), append(args, *update.ShortcutsCount)
	}
	if update.UsersCount != nil {
		argIndex++
		set, args = append(set, "users_count = "+placeholder(argIndex)), append(args, *update.UsersCount)
	}
	if update.CollectionsCount != nil {
		argIndex++
		set, args = append(set, "collections_count = "+placeholder(argIndex)), append(args, *update.CollectionsCount)
	}
	if update.HitsCount != nil {
		argIndex++
		set, args = append(set, "hits_count = "+placeholder(argIndex)), append(args, *update.HitsCount)
	}

	argIndex++
	args = append(args, update.ID)
	stmt := `UPDATE stats_measurement SET ` + strings.Join(set, ", ") + ` WHERE id = ` + placeholder(argIndex) + ` RETURNING id, measured_ts, shortcuts_count, users_count, collections_count, hits_count`
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
	stmt := `DELETE FROM stats_measurement WHERE id = $1`
	result, err := d.db.ExecContext(ctx, stmt, delete.ID)
	if err != nil {
		return err
	}
	if _, err = result.RowsAffected(); err != nil {
		return err
	}
	return nil
}