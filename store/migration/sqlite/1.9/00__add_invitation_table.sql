-- Create invitation table
CREATE TABLE invitation (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  from_path TEXT NOT NULL,
  to_path TEXT NOT NULL,
  accepted_at TEXT,
  deleted_at TEXT
);

CREATE INDEX idx_invitation_from_path ON invitation(from_path);
CREATE INDEX idx_invitation_to_path ON invitation(to_path);
CREATE INDEX idx_invitation_accepted_at ON invitation(accepted_at);
CREATE INDEX idx_invitation_deleted_at ON invitation(deleted_at);