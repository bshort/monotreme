package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/bshort/monotreme/store"
	storepb "github.com/bshort/monotreme/proto/gen/store"
)

func (d *DB) ListFollowing(ctx context.Context, userID int32) ([]*store.Following, error) {
	var query string
	var args []interface{}

	// If userID is 0, get all following relationships
	if userID == 0 {
		query = `
			SELECT id, follower_id, following_id, created_ts
			FROM following
			ORDER BY created_ts DESC
		`
	} else {
		// Get following relationships for a specific user
		query = `
			SELECT id, follower_id, following_id, created_ts
			FROM following
			WHERE follower_id = ?
			ORDER BY created_ts DESC
		`
		args = []interface{}{userID}
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	followings := []*store.Following{}
	for rows.Next() {
		following := &store.Following{}
		if err := rows.Scan(
			&following.ID,
			&following.FollowerID,
			&following.FollowingID,
			&following.CreatedTs,
		); err != nil {
			return nil, err
		}
		followings = append(followings, following)
	}

	return followings, rows.Err()
}

func (d *DB) ListFollowers(ctx context.Context, userID int32) ([]*store.Following, error) {
	var query string
	var args []interface{}

	// If userID is 0, get all follower relationships
	if userID == 0 {
		query = `
			SELECT id, follower_id, following_id, created_ts
			FROM following
			ORDER BY created_ts DESC
		`
	} else {
		// Get followers for a specific user
		query = `
			SELECT id, follower_id, following_id, created_ts
			FROM following
			WHERE following_id = ?
			ORDER BY created_ts DESC
		`
		args = []interface{}{userID}
	}

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	followings := []*store.Following{}
	for rows.Next() {
		following := &store.Following{}
		if err := rows.Scan(
			&following.ID,
			&following.FollowerID,
			&following.FollowingID,
			&following.CreatedTs,
		); err != nil {
			return nil, err
		}
		followings = append(followings, following)
	}

	return followings, rows.Err()
}

func (d *DB) GetFollowing(ctx context.Context, id int32) (*store.Following, error) {
	query := `
		SELECT id, follower_id, following_id, created_ts
		FROM following
		WHERE id = ?
	`

	following := &store.Following{}
	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&following.ID,
		&following.FollowerID,
		&following.FollowingID,
		&following.CreatedTs,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return following, nil
}

func (d *DB) CreateFollowing(ctx context.Context, followerID int32, followingID int32) (*store.Following, error) {
	query := `
		INSERT INTO following (follower_id, following_id)
		VALUES (?, ?)
		RETURNING id, follower_id, following_id, created_ts
	`

	following := &store.Following{}
	err := d.db.QueryRowContext(ctx, query, followerID, followingID).Scan(
		&following.ID,
		&following.FollowerID,
		&following.FollowingID,
		&following.CreatedTs,
	)
	if err != nil {
		return nil, err
	}

	return following, nil
}

func (d *DB) DeleteFollowing(ctx context.Context, id int32) error {
	query := `DELETE FROM following WHERE id = ?`
	_, err := d.db.ExecContext(ctx, query, id)
	return err
}

func (d *DB) GetFollowingUserShortcuts(ctx context.Context, followerID int32, followingID int32) ([]*storepb.Shortcut, error) {
	// First verify the following relationship exists
	query := `
		SELECT COUNT(*) FROM following
		WHERE follower_id = ? AND following_id = ?
	`
	var count int
	err := d.db.QueryRowContext(ctx, query, followerID, followingID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		// Not following this user, return empty list
		return []*storepb.Shortcut{}, nil
	}

	// Get public shortcuts from the user being followed
	shortcutQuery := `
		SELECT
			id,
			creator_id,
			created_ts,
			updated_ts,
			name,
			link,
			title,
			description,
			visibility,
			tag,
			og_metadata,
			uuid,
			custom_icon,
			user_order
		FROM shortcut
		WHERE creator_id = ? AND visibility = 'PUBLIC'
		ORDER BY created_ts DESC
	`

	rows, err := d.db.QueryContext(ctx, shortcutQuery, followingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shortcuts := []*storepb.Shortcut{}
	for rows.Next() {
		shortcut := &storepb.Shortcut{}
		var visibility, tags, openGraphMetadataString string
		var customIcon sql.NullString
		var userOrder sql.NullInt32

		if err := rows.Scan(
			&shortcut.Id,
			&shortcut.CreatorId,
			&shortcut.CreatedTs,
			&shortcut.UpdatedTs,
			&shortcut.Name,
			&shortcut.Link,
			&shortcut.Title,
			&shortcut.Description,
			&visibility,
			&tags,
			&openGraphMetadataString,
			&shortcut.Uuid,
			&customIcon,
			&userOrder,
		); err != nil {
			return nil, err
		}

		shortcut.Visibility = store.ConvertVisibilityStringToStorepb(visibility)
		shortcut.Tags = filterTags(strings.Split(tags, " "))

		var ogMetadata storepb.OpenGraphMetadata
		if openGraphMetadataString != "" {
			if err := protojson.Unmarshal([]byte(openGraphMetadataString), &ogMetadata); err == nil {
				shortcut.OgMetadata = &ogMetadata
			}
		}

		if customIcon.Valid {
			shortcut.CustomIcon = customIcon.String
		}
		if userOrder.Valid {
			shortcut.UserOrder = userOrder.Int32
		}

		shortcuts = append(shortcuts, shortcut)
	}

	return shortcuts, rows.Err()
}
