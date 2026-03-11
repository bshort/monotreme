-- following table for user following relationships
CREATE TABLE following (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  follower_id INTEGER NOT NULL,
  following_id INTEGER NOT NULL,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  FOREIGN KEY (follower_id) REFERENCES user(id) ON DELETE CASCADE,
  FOREIGN KEY (following_id) REFERENCES user(id) ON DELETE CASCADE,
  UNIQUE(follower_id, following_id)
);

CREATE INDEX idx_following_follower_id ON following(follower_id);
CREATE INDEX idx_following_following_id ON following(following_id);
