-- Create tag table
CREATE TABLE tag (
  uuid TEXT PRIMARY KEY,
  creator_id INTEGER NOT NULL,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  name TEXT NOT NULL,
  abbreviation TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  FOREIGN KEY (creator_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE INDEX idx_tag_abbreviation ON tag(abbreviation);
CREATE INDEX idx_tag_creator_id ON tag(creator_id);
