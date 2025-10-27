-- migration_history
CREATE TABLE migration_history (
  version TEXT NOT NULL PRIMARY KEY,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- workspace_setting
CREATE TABLE workspace_setting (
  key TEXT NOT NULL UNIQUE,
  value TEXT NOT NULL
);

-- user
CREATE TABLE user (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  row_status TEXT NOT NULL CHECK (row_status IN ('NORMAL', 'ARCHIVED')) DEFAULT 'NORMAL',
  email TEXT NOT NULL UNIQUE,
  nickname TEXT NOT NULL,
  password_hash TEXT NOT NULL,
  profile_picture TEXT NOT NULL DEFAULT '',
  role TEXT NOT NULL CHECK (role IN ('ADMIN', 'USER')) DEFAULT 'USER',
  invited_by INTEGER REFERENCES user(id),
  uuid TEXT,
  locale TEXT NOT NULL DEFAULT 'EN',
  color_theme TEXT NOT NULL DEFAULT 'SYSTEM',
  default_visibility TEXT NOT NULL DEFAULT 'WORKSPACE',
  auto_generate_title BOOLEAN NOT NULL DEFAULT true,
  auto_generate_icon BOOLEAN NOT NULL DEFAULT true,
  auto_generate_name BOOLEAN NOT NULL DEFAULT true,
  edit_mode_preference TEXT NOT NULL DEFAULT 'FLYOUT'
);

CREATE INDEX idx_user_email ON user(email);

-- user_setting
CREATE TABLE user_setting (
  user_id INTEGER NOT NULL,
  key TEXT NOT NULL,
  value TEXT NOT NULL,
  UNIQUE(user_id, key)
);

-- shortcut
CREATE TABLE shortcut (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  creator_id INTEGER NOT NULL,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  row_status TEXT NOT NULL CHECK (row_status IN ('NORMAL', 'ARCHIVED')) DEFAULT 'NORMAL',
  name TEXT NOT NULL UNIQUE,
  link TEXT NOT NULL,
  title TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  visibility TEXT NOT NULL CHECK (visibility IN ('PRIVATE', 'WORKSPACE', 'PUBLIC')) DEFAULT 'PRIVATE',
  tag TEXT NOT NULL DEFAULT '',
  og_metadata TEXT NOT NULL DEFAULT '{}',
  uuid TEXT NOT NULL DEFAULT '',
  custom_icon TEXT NOT NULL DEFAULT '',
  user_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_shortcut_name ON shortcut(name);
CREATE INDEX idx_shortcut_uuid ON shortcut(uuid);

-- activity
CREATE TABLE activity (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  creator_id INTEGER NOT NULL,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  type TEXT NOT NULL DEFAULT '',
  level TEXT NOT NULL CHECK (level IN ('INFO', 'WARN', 'ERROR')) DEFAULT 'INFO',
  payload TEXT NOT NULL DEFAULT '{}'
);

-- collection
CREATE TABLE collection (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  creator_id INTEGER NOT NULL,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  name TEXT NOT NULL UNIQUE,
  title TEXT NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  shortcut_ids INTEGER[] NOT NULL,
  visibility TEXT NOT NULL CHECK (visibility IN ('PRIVATE', 'WORKSPACE', 'PUBLIC')) DEFAULT 'PRIVATE',
  custom_icon TEXT NOT NULL DEFAULT '',
  uuid TEXT
);

CREATE INDEX idx_collection_name ON collection(name);

-- stats_measurement
CREATE TABLE stats_measurement (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  measured_ts BIGINT NOT NULL,
  shortcuts_count INTEGER NOT NULL DEFAULT 0,
  users_count INTEGER NOT NULL DEFAULT 0,
  collections_count INTEGER NOT NULL DEFAULT 0,
  hits_count INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_stats_measurement_measured_ts ON stats_measurement(measured_ts);

-- rss_feed
CREATE TABLE rss_feed (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL DEFAULT '',
  creator_id INTEGER NOT NULL,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  updated_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  row_status TEXT NOT NULL CHECK (row_status IN ('NORMAL', 'ARCHIVED')) DEFAULT 'NORMAL',

  title TEXT NOT NULL DEFAULT '',
  url TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',

  auto_import BOOLEAN NOT NULL DEFAULT false,
  import_frequency_hours INTEGER NOT NULL DEFAULT 24,
  last_import_ts BIGINT,

  default_tags TEXT NOT NULL DEFAULT '',
  default_visibility TEXT NOT NULL CHECK (default_visibility IN ('PRIVATE', 'WORKSPACE', 'PUBLIC')) DEFAULT 'WORKSPACE',
  shortcut_prefix TEXT NOT NULL DEFAULT '',

  is_active BOOLEAN NOT NULL DEFAULT true,
  last_error TEXT NOT NULL DEFAULT '',
  total_imported INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_rss_feed_creator_id ON rss_feed(creator_id);
CREATE INDEX idx_rss_feed_auto_import ON rss_feed(auto_import, is_active);
CREATE INDEX idx_rss_feed_uuid ON rss_feed(uuid);

-- rss_feed_item
CREATE TABLE rss_feed_item (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL DEFAULT '',
  rss_feed_id INTEGER NOT NULL,
  created_ts BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
  row_status TEXT NOT NULL CHECK (row_status IN ('NORMAL', 'ARCHIVED')) DEFAULT 'NORMAL',

  item_guid TEXT NOT NULL,
  item_link TEXT NOT NULL,
  item_title TEXT NOT NULL DEFAULT '',
  item_description TEXT NOT NULL DEFAULT '',
  item_published_ts BIGINT,

  shortcut_id INTEGER,
  import_success BOOLEAN NOT NULL DEFAULT false,
  import_error TEXT NOT NULL DEFAULT '',

  FOREIGN KEY (rss_feed_id) REFERENCES rss_feed(id) ON DELETE CASCADE,
  UNIQUE (rss_feed_id, item_guid)
);

CREATE INDEX idx_rss_feed_item_rss_feed_id ON rss_feed_item(rss_feed_id);
CREATE INDEX idx_rss_feed_item_guid ON rss_feed_item(rss_feed_id, item_guid);
CREATE INDEX idx_rss_feed_item_uuid ON rss_feed_item(uuid);

-- invitation
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
-- tag
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

-- bookmark_tag junction table
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
