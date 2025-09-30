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