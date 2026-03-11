package store

import (
	"context"

	storepb "github.com/bshort/monotreme/proto/gen/store"
)

// StatsMeasurement represents a single stats measurement record.
type StatsMeasurement struct {
	ID               int32
	MeasuredTs       int64
	ShortcutsCount   int32
	UsersCount       int32
	CollectionsCount int32
	HitsCount        int32
}

type FindStatsMeasurement struct {
	ID    *int32
	Limit *int32
	// If true, order by measured_ts DESC, else ASC.
	OrderByMeasuredTsDesc *bool
}

type UpdateStatsMeasurement struct {
	ID               int32
	ShortcutsCount   *int32
	UsersCount       *int32
	CollectionsCount *int32
	HitsCount        *int32
}

type DeleteStatsMeasurement struct {
	ID *int32
}

func (s *Store) CreateStatsMeasurement(ctx context.Context, create *StatsMeasurement) (*StatsMeasurement, error) {
	return s.driver.CreateStatsMeasurement(ctx, create)
}

func (s *Store) ListStatsMeasurements(ctx context.Context, find *FindStatsMeasurement) ([]*StatsMeasurement, error) {
	return s.driver.ListStatsMeasurements(ctx, find)
}

func (s *Store) GetStatsMeasurement(ctx context.Context, find *FindStatsMeasurement) (*StatsMeasurement, error) {
	list, err := s.ListStatsMeasurements(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (s *Store) UpdateStatsMeasurement(ctx context.Context, update *UpdateStatsMeasurement) (*StatsMeasurement, error) {
	return s.driver.UpdateStatsMeasurement(ctx, update)
}

func (s *Store) DeleteStatsMeasurement(ctx context.Context, delete *DeleteStatsMeasurement) error {
	return s.driver.DeleteStatsMeasurement(ctx, delete)
}

// Utility function to convert store StatsMeasurement to protobuf
func (s *StatsMeasurement) ToStorePb() *storepb.StatsMeasurement {
	return &storepb.StatsMeasurement{
		Id:               s.ID,
		MeasuredTs:       s.MeasuredTs,
		ShortcutsCount:   s.ShortcutsCount,
		UsersCount:       s.UsersCount,
		CollectionsCount: s.CollectionsCount,
		HitsCount:        s.HitsCount,
	}
}

// Utility function to convert protobuf StatsMeasurement to store
func NewStatsMeasurementFromStorePb(pb *storepb.StatsMeasurement) *StatsMeasurement {
	return &StatsMeasurement{
		ID:               pb.Id,
		MeasuredTs:       pb.MeasuredTs,
		ShortcutsCount:   pb.ShortcutsCount,
		UsersCount:       pb.UsersCount,
		CollectionsCount: pb.CollectionsCount,
		HitsCount:        pb.HitsCount,
	}
}