package store

import (
	"context"

	storepb "github.com/bshort/monotreme/proto/gen/store"
)

type UpdateTag struct {
	UUID         string
	Name         *string
	Abbreviation *string
	Description  *string
}

type FindTag struct {
	UUID         *string
	CreatorID    *int32
	Abbreviation *string
}

type DeleteTag struct {
	UUID string
}

func (s *Store) CreateTag(ctx context.Context, create *storepb.Tag) (*storepb.Tag, error) {
	return s.driver.CreateTag(ctx, create)
}

func (s *Store) UpdateTag(ctx context.Context, update *UpdateTag) (*storepb.Tag, error) {
	return s.driver.UpdateTag(ctx, update)
}

func (s *Store) ListTags(ctx context.Context, find *FindTag) ([]*storepb.Tag, error) {
	return s.driver.ListTags(ctx, find)
}

func (s *Store) GetTag(ctx context.Context, find *FindTag) (*storepb.Tag, error) {
	tags, err := s.ListTags(ctx, find)
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return nil, nil
	}

	tag := tags[0]
	return tag, nil
}

func (s *Store) DeleteTag(ctx context.Context, delete *DeleteTag) error {
	return s.driver.DeleteTag(ctx, delete)
}
