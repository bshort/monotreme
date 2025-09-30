package store

import (
	"context"
	"time"

	storepb "github.com/bshort/monotreme/proto/gen/store"
)

type RssFeed struct {
	ID          int32
	UUID        string
	CreatorID   int32
	CreatedTs   *time.Time
	UpdatedTs   *time.Time

	Title                string
	URL                  string
	Description          string

	// Import settings
	AutoImport           bool
	ImportFrequencyHours int32
	LastImportTs         *time.Time

	// Shortcut generation settings
	DefaultTags          []string
	DefaultVisibility    storepb.Visibility
	ShortcutPrefix       string

	// Status
	IsActive             bool
	LastError            string
	TotalImported        int32

	RowStatus            storepb.RowStatus
}

type RssFeedItem struct {
	ID              int32
	UUID            string
	RssFeedID       int32
	CreatedTs       *time.Time

	ItemGUID        string
	ItemLink        string
	ItemTitle       string
	ItemDescription string
	ItemPublishedTs *time.Time

	ShortcutID      *int32
	ImportSuccess   bool
	ImportError     string

	RowStatus       storepb.RowStatus
}

type UpdateRssFeed struct {
	ID                   int32

	Title                *string
	URL                  *string
	Description          *string
	AutoImport           *bool
	ImportFrequencyHours *int32
	LastImportTime       *time.Time
	DefaultTags          []string
	DefaultVisibility    *storepb.Visibility
	ShortcutPrefix       *string
	IsActive             *bool
	LastError            *string
	TotalImported        *int32
}

type FindRssFeed struct {
	ID        *int32
	CreatorID *int32
	IsActive  *bool
}

type DeleteRssFeed struct {
	ID int32
}

type FindRssFeedItem struct {
	ID         *int32
	RssFeedID  *int32
	ItemGUID   *string
	ItemLink   *string
}

func (s *Store) CreateRssFeed(ctx context.Context, create *RssFeed) (*RssFeed, error) {
	rssFeed, err := s.driver.CreateRssFeed(ctx, create)
	if err != nil {
		return nil, err
	}
	return rssFeed, nil
}

func (s *Store) UpdateRssFeed(ctx context.Context, update *UpdateRssFeed) (*RssFeed, error) {
	rssFeed, err := s.driver.UpdateRssFeed(ctx, update)
	if err != nil {
		return nil, err
	}
	return rssFeed, nil
}

func (s *Store) ListRssFeeds(ctx context.Context, find *FindRssFeed) ([]*RssFeed, error) {
	return s.driver.ListRssFeeds(ctx, find)
}

func (s *Store) GetRssFeed(ctx context.Context, find *FindRssFeed) (*RssFeed, error) {
	list, err := s.driver.ListRssFeeds(ctx, find)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}

func (s *Store) DeleteRssFeed(ctx context.Context, delete *DeleteRssFeed) error {
	return s.driver.DeleteRssFeed(ctx, delete)
}

func (s *Store) CreateRssFeedItem(ctx context.Context, create *RssFeedItem) (*RssFeedItem, error) {
	return s.driver.CreateRssFeedItem(ctx, create)
}

func (s *Store) CheckRssFeedItemExists(ctx context.Context, rssFeedID int32, itemGUID string) (bool, error) {
	item, err := s.driver.GetRssFeedItem(ctx, &FindRssFeedItem{
		RssFeedID: &rssFeedID,
		ItemGUID:  &itemGUID,
	})
	if err != nil {
		return false, err
	}
	return item != nil, nil
}

func (s *Store) ListRssFeedItems(ctx context.Context, find *FindRssFeedItem) ([]*RssFeedItem, error) {
	return s.driver.ListRssFeedItems(ctx, find)
}