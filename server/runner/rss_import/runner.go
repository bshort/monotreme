// Package rss_import provides a runner to automatically import RSS feeds based on their schedule.
package rss_import

import (
	"context"
	"log/slog"
	"time"

	"github.com/bshort/monotreme/store"
)

type Runner struct {
	Store       *store.Store
	apiV1Service RSSImportService
}

// RSSImportService interface to avoid circular dependency
type RSSImportService interface {
	ImportRssFeed(ctx context.Context, rssFeed *store.RssFeed) error
}

func NewRunner(store *store.Store, apiV1Service RSSImportService) *Runner {
	return &Runner{
		Store:       store,
		apiV1Service: apiV1Service,
	}
}

// Check for RSS feeds to import every 5 minutes
const runnerInterval = 5 * time.Minute

func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(runnerInterval)
	defer ticker.Stop()

	slog.Info("RSS import runner started", "interval", runnerInterval.String())

	for {
		select {
		case <-ticker.C:
			r.RunOnce(ctx)
		case <-ctx.Done():
			slog.Info("RSS import runner shutting down")
			return
		}
	}
}

func (r *Runner) RunOnce(ctx context.Context) {
	if err := r.processScheduledImports(ctx); err != nil {
		slog.Error("failed to process scheduled RSS imports", "error", err)
	}
}

func (r *Runner) processScheduledImports(ctx context.Context) error {
	// Get all active RSS feeds with auto-import enabled
	activeTrue := true
	rssFeeds, err := r.Store.ListRssFeeds(ctx, &store.FindRssFeed{
		IsActive: &activeTrue,
	})
	if err != nil {
		return err
	}

	now := time.Now()
	processed := 0

	for _, feed := range rssFeeds {
		// Skip if auto-import is disabled
		if !feed.AutoImport {
			continue
		}

		// Check if it's time to import this feed
		if r.shouldImportFeed(feed, now) {
			slog.Info("processing scheduled RSS import",
				"feed_id", feed.ID,
				"title", feed.Title,
				"url", feed.URL,
				"frequency_hours", feed.ImportFrequencyHours)

			if err := r.apiV1Service.ImportRssFeed(ctx, feed); err != nil {
				slog.Error("failed to import RSS feed",
					"feed_id", feed.ID,
					"title", feed.Title,
					"url", feed.URL,
					"error", err)

				// Update feed with error information
				errorMsg := err.Error()
				if _, updateErr := r.Store.UpdateRssFeed(ctx, &store.UpdateRssFeed{
					ID:        feed.ID,
					LastError: &errorMsg,
				}); updateErr != nil {
					slog.Error("failed to update RSS feed error", "feed_id", feed.ID, "error", updateErr)
				}
			} else {
				slog.Info("successfully imported RSS feed",
					"feed_id", feed.ID,
					"title", feed.Title)

				// Clear any previous error
				emptyError := ""
				if _, updateErr := r.Store.UpdateRssFeed(ctx, &store.UpdateRssFeed{
					ID:        feed.ID,
					LastError: &emptyError,
				}); updateErr != nil {
					slog.Error("failed to clear RSS feed error", "feed_id", feed.ID, "error", updateErr)
				}
			}
			processed++
		}
	}

	if processed > 0 {
		slog.Info("completed RSS import batch", "feeds_processed", processed)
	}

	return nil
}

func (r *Runner) shouldImportFeed(feed *store.RssFeed, now time.Time) bool {
	// If this feed has never been imported, import it now
	if feed.LastImportTs == nil {
		return true
	}

	// Calculate the next import time based on frequency
	nextImportTime := feed.LastImportTs.Add(time.Duration(feed.ImportFrequencyHours) * time.Hour)

	// Import if we've passed the next import time
	return now.After(nextImportTime)
}