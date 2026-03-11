package postgres

import (
	"context"

	"github.com/pkg/errors"

	storepb "github.com/bshort/monotreme/proto/gen/store"
	"github.com/bshort/monotreme/store"
)

func (d *DB) CreateTag(ctx context.Context, create *storepb.Tag) (*storepb.Tag, error) {
	return nil, errors.New("not implemented for postgres")
}

func (d *DB) UpdateTag(ctx context.Context, update *store.UpdateTag) (*storepb.Tag, error) {
	return nil, errors.New("not implemented for postgres")
}

func (d *DB) ListTags(ctx context.Context, find *store.FindTag) ([]*storepb.Tag, error) {
	return nil, errors.New("not implemented for postgres")
}

func (d *DB) DeleteTag(ctx context.Context, delete *store.DeleteTag) error {
	return errors.New("not implemented for postgres")
}
