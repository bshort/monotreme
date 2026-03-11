-- Manual migration to add stats_measurement table
-- Run this if you're getting "no such table: stats_measurement" error

-- For SQLite:
CREATE TABLE IF NOT EXISTS stats_measurement (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  measured_ts BIGINT NOT NULL,
  shortcuts_count INTEGER NOT NULL DEFAULT 0,
  users_count INTEGER NOT NULL DEFAULT 0,
  collections_count INTEGER NOT NULL DEFAULT 0,
  hits_count INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_stats_measurement_measured_ts ON stats_measurement(measured_ts);

-- For PostgreSQL, replace the above with:
-- CREATE TABLE IF NOT EXISTS stats_measurement (
--   id SERIAL PRIMARY KEY,
--   measured_ts BIGINT NOT NULL,
--   shortcuts_count INTEGER NOT NULL DEFAULT 0,
--   users_count INTEGER NOT NULL DEFAULT 0,
--   collections_count INTEGER NOT NULL DEFAULT 0,
--   hits_count INTEGER NOT NULL DEFAULT 0
-- );
--
-- CREATE INDEX IF NOT EXISTS idx_stats_measurement_measured_ts ON stats_measurement(measured_ts);