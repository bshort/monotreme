-- Create bookmark_tag junction table for many-to-many relationship
CREATE TABLE bookmark_tag (
  shortcut_id INTEGER NOT NULL,
  tag_uuid TEXT NOT NULL,
  user_id INTEGER NOT NULL,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  PRIMARY KEY (shortcut_id, tag_uuid),
  FOREIGN KEY (shortcut_id) REFERENCES shortcut(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_uuid) REFERENCES tag(uuid) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
);

CREATE INDEX idx_bookmark_tag_shortcut_id ON bookmark_tag(shortcut_id);
CREATE INDEX idx_bookmark_tag_tag_uuid ON bookmark_tag(tag_uuid);
CREATE INDEX idx_bookmark_tag_user_id ON bookmark_tag(user_id);
CREATE INDEX idx_bookmark_tag_user_tag ON bookmark_tag(user_id, tag_uuid);
