// Package stats provides a runner to collect workspace statistics measurements hourly.
package stats

import (
	"context"
	"log/slog"
	"time"

	"github.com/bshort/monotreme/store"
)

type Runner struct {
	Store *store.Store
}

func NewRunner(store *store.Store) *Runner {
	return &Runner{
		Store: store,
	}
}

// Schedule stats collection every hour.
const runnerInterval = time.Hour

func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(runnerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.RunOnce(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (r *Runner) RunOnce(ctx context.Context) {
	if err := r.collectStats(ctx); err != nil {
		slog.Error("failed to collect stats measurement", "error", err)
	}
}

func (r *Runner) collectStats(ctx context.Context) error {
	now := time.Now().Unix()

	// Get total shortcuts count
	shortcuts, err := r.Store.ListShortcuts(ctx, &store.FindShortcut{})
	if err != nil {
		return err
	}
	totalShortcuts := int32(len(shortcuts))

	// Get total users count
	users, err := r.Store.ListUsers(ctx, &store.FindUser{})
	if err != nil {
		return err
	}
	totalUsers := int32(len(users))

	// Get total collections count
	collections, err := r.Store.ListCollections(ctx, &store.FindCollection{})
	if err != nil {
		return err
	}
	totalCollections := int32(len(collections))

	// Get total hits count (shortcut views)
	activities, err := r.Store.ListActivities(ctx, &store.FindActivity{
		Type: store.ActivityShortcutView,
	})
	if err != nil {
		return err
	}
	totalHits := int32(len(activities))

	// Create stats measurement record
	measurement := &store.StatsMeasurement{
		MeasuredTs:       now,
		ShortcutsCount:   totalShortcuts,
		UsersCount:       totalUsers,
		CollectionsCount: totalCollections,
		HitsCount:        totalHits,
	}

	_, err = r.Store.CreateStatsMeasurement(ctx, measurement)
	if err != nil {
		return err
	}

	slog.Info("stats measurement collected",
		"shortcuts", totalShortcuts,
		"users", totalUsers,
		"collections", totalCollections,
		"hits", totalHits)

	return nil
}