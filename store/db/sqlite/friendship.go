package sqlite

import (
	"context"
	"database/sql"

	"github.com/bshort/monotreme/store"
)

func (d *DB) ListFriendships(ctx context.Context, userID int32, status string) ([]*store.Friendship, error) {
	query := `
		SELECT id, user_id, friend_id, status, created_ts, accepted_ts
		FROM friendship
		WHERE (user_id = ? OR friend_id = ?) AND status = ?
		ORDER BY created_ts DESC
	`

	rows, err := d.db.QueryContext(ctx, query, userID, userID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friendships := []*store.Friendship{}
	for rows.Next() {
		friendship := &store.Friendship{}
		var acceptedTs sql.NullInt64
		if err := rows.Scan(
			&friendship.ID,
			&friendship.UserID,
			&friendship.FriendID,
			&friendship.Status,
			&friendship.CreatedTs,
			&acceptedTs,
		); err != nil {
			return nil, err
		}
		if acceptedTs.Valid {
			ts := store.Timestamp(acceptedTs.Int64)
			friendship.AcceptedTs = &ts
		}
		friendships = append(friendships, friendship)
	}

	return friendships, rows.Err()
}

func (d *DB) ListIncomingFriendRequests(ctx context.Context, userID int32) ([]*store.Friendship, error) {
	query := `
		SELECT id, user_id, friend_id, status, created_ts, accepted_ts
		FROM friendship
		WHERE friend_id = ? AND status = 'PENDING'
		ORDER BY created_ts DESC
	`

	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friendships := []*store.Friendship{}
	for rows.Next() {
		friendship := &store.Friendship{}
		var acceptedTs sql.NullInt64
		if err := rows.Scan(
			&friendship.ID,
			&friendship.UserID,
			&friendship.FriendID,
			&friendship.Status,
			&friendship.CreatedTs,
			&acceptedTs,
		); err != nil {
			return nil, err
		}
		if acceptedTs.Valid {
			ts := store.Timestamp(acceptedTs.Int64)
			friendship.AcceptedTs = &ts
		}
		friendships = append(friendships, friendship)
	}

	return friendships, rows.Err()
}

func (d *DB) ListOutgoingFriendRequests(ctx context.Context, userID int32) ([]*store.Friendship, error) {
	query := `
		SELECT id, user_id, friend_id, status, created_ts, accepted_ts
		FROM friendship
		WHERE user_id = ? AND status = 'PENDING'
		ORDER BY created_ts DESC
	`

	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friendships := []*store.Friendship{}
	for rows.Next() {
		friendship := &store.Friendship{}
		var acceptedTs sql.NullInt64
		if err := rows.Scan(
			&friendship.ID,
			&friendship.UserID,
			&friendship.FriendID,
			&friendship.Status,
			&friendship.CreatedTs,
			&acceptedTs,
		); err != nil {
			return nil, err
		}
		if acceptedTs.Valid {
			ts := store.Timestamp(acceptedTs.Int64)
			friendship.AcceptedTs = &ts
		}
		friendships = append(friendships, friendship)
	}

	return friendships, rows.Err()
}

func (d *DB) GetFriendship(ctx context.Context, id int32) (*store.Friendship, error) {
	query := `
		SELECT id, user_id, friend_id, status, created_ts, accepted_ts
		FROM friendship
		WHERE id = ?
	`

	friendship := &store.Friendship{}
	var acceptedTs sql.NullInt64
	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&friendship.ID,
		&friendship.UserID,
		&friendship.FriendID,
		&friendship.Status,
		&friendship.CreatedTs,
		&acceptedTs,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if acceptedTs.Valid {
		ts := store.Timestamp(acceptedTs.Int64)
		friendship.AcceptedTs = &ts
	}

	return friendship, nil
}

func (d *DB) CreateFriendship(ctx context.Context, userID int32, friendID int32) (*store.Friendship, error) {
	query := `
		INSERT INTO friendship (user_id, friend_id, status)
		VALUES (?, ?, 'PENDING')
		RETURNING id, user_id, friend_id, status, created_ts, accepted_ts
	`

	friendship := &store.Friendship{}
	var acceptedTs sql.NullInt64
	err := d.db.QueryRowContext(ctx, query, userID, friendID).Scan(
		&friendship.ID,
		&friendship.UserID,
		&friendship.FriendID,
		&friendship.Status,
		&friendship.CreatedTs,
		&acceptedTs,
	)
	if err != nil {
		return nil, err
	}

	if acceptedTs.Valid {
		ts := store.Timestamp(acceptedTs.Int64)
		friendship.AcceptedTs = &ts
	}

	return friendship, nil
}

func (d *DB) AcceptFriendship(ctx context.Context, id int32) error {
	query := `
		UPDATE friendship
		SET status = 'ACCEPTED', accepted_ts = strftime('%s', 'now')
		WHERE id = ?
	`

	_, err := d.db.ExecContext(ctx, query, id)
	return err
}

func (d *DB) DeleteFriendship(ctx context.Context, id int32) error {
	query := `DELETE FROM friendship WHERE id = ?`
	_, err := d.db.ExecContext(ctx, query, id)
	return err
}
