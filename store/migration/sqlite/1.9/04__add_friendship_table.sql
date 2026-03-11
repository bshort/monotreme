-- friendship table for friend requests
CREATE TABLE friendship (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  friend_id INTEGER NOT NULL,
  status TEXT NOT NULL DEFAULT 'PENDING',
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  accepted_ts BIGINT,
  FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
  FOREIGN KEY (friend_id) REFERENCES user(id) ON DELETE CASCADE,
  UNIQUE(user_id, friend_id)
);

CREATE INDEX idx_friendship_user_id ON friendship(user_id);
CREATE INDEX idx_friendship_friend_id ON friendship(friend_id);
CREATE INDEX idx_friendship_status ON friendship(status);
