-- stats_measurement table for tracking workspace statistics over time
CREATE TABLE stats_measurement (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  measured_ts BIGINT NOT NULL,
  shortcuts_count INTEGER NOT NULL DEFAULT 0,
  users_count INTEGER NOT NULL DEFAULT 0,
  collections_count INTEGER NOT NULL DEFAULT 0,
  hits_count INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_stats_measurement_measured_ts ON stats_measurement(measured_ts);