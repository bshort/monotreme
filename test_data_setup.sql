-- Create 5 test users
-- Note: password_hash is empty since these are test users
INSERT INTO user (email, nickname, password_hash, profile_picture, role, created_ts, updated_ts)
VALUES
  ('alice@test.com', 'Alice Smith', '', '', 'USER', strftime('%s', 'now'), strftime('%s', 'now')),
  ('bob@test.com', 'Bob Johnson', '', '', 'USER', strftime('%s', 'now'), strftime('%s', 'now')),
  ('carol@test.com', 'Carol Williams', '', '', 'USER', strftime('%s', 'now'), strftime('%s', 'now')),
  ('dave@test.com', 'Dave Brown', '', '', 'USER', strftime('%s', 'now'), strftime('%s', 'now')),
  ('eve@test.com', 'Eve Davis', '', '', 'USER', strftime('%s', 'now'), strftime('%s', 'now'));

-- Get bshort's user ID (should be based on uuid)
-- Create friendship relationships with bshort
-- Status = ACCEPTED for active friendships
INSERT INTO friendship (user_id, friend_id, status, created_ts, accepted_ts)
SELECT
  (SELECT id FROM user WHERE uuid = 'b287826f-6133-452c-be06-480622be512d'),
  id,
  'ACCEPTED',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user
WHERE email IN ('alice@test.com', 'bob@test.com', 'carol@test.com', 'dave@test.com', 'eve@test.com');

-- Create reciprocal friendships (friend_id -> user_id)
INSERT INTO friendship (user_id, friend_id, status, created_ts, accepted_ts)
SELECT
  id,
  (SELECT id FROM user WHERE uuid = 'b287826f-6133-452c-be06-480622be512d'),
  'ACCEPTED',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user
WHERE email IN ('alice@test.com', 'bob@test.com', 'carol@test.com', 'dave@test.com', 'eve@test.com');

-- Create test bookmarks for Alice (PUBLIC, PRIVATE, FRIENDS)
INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'alice-public-1',
  'https://example.com/alice/public1',
  'Alice Public Bookmark 1',
  'This is a public bookmark from Alice',
  'PUBLIC',
  'test personal',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'alice@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'alice-private-1',
  'https://example.com/alice/private1',
  'Alice Private Bookmark 1',
  'This is a private bookmark from Alice',
  'PRIVATE',
  'test secret',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'alice@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'alice-friends-1',
  'https://example.com/alice/friends1',
  'Alice Friends Bookmark 1',
  'This is a friends-only bookmark from Alice',
  'FRIENDS',
  'test friends-only',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'alice@test.com';

-- Create test bookmarks for Bob
INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'bob-public-1',
  'https://example.com/bob/public1',
  'Bob Public Bookmark 1',
  'This is a public bookmark from Bob',
  'PUBLIC',
  'test tech',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'bob@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'bob-private-1',
  'https://example.com/bob/private1',
  'Bob Private Bookmark 1',
  'This is a private bookmark from Bob',
  'PRIVATE',
  'test confidential',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'bob@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'bob-friends-1',
  'https://example.com/bob/friends1',
  'Bob Friends Bookmark 1',
  'This is a friends-only bookmark from Bob',
  'FRIENDS',
  'test friends-only',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'bob@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'bob-public-2',
  'https://example.com/bob/public2',
  'Bob Public Bookmark 2',
  'Another public bookmark from Bob',
  'PUBLIC',
  'test programming',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'bob@test.com';

-- Create test bookmarks for Carol
INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'carol-public-1',
  'https://example.com/carol/public1',
  'Carol Public Bookmark 1',
  'This is a public bookmark from Carol',
  'PUBLIC',
  'test design',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'carol@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'carol-private-1',
  'https://example.com/carol/private1',
  'Carol Private Bookmark 1',
  'This is a private bookmark from Carol',
  'PRIVATE',
  'test personal',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'carol@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'carol-friends-1',
  'https://example.com/carol/friends1',
  'Carol Friends Bookmark 1',
  'This is a friends-only bookmark from Carol',
  'FRIENDS',
  'test friends-only',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'carol@test.com';

-- Create test bookmarks for Dave
INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'dave-public-1',
  'https://example.com/dave/public1',
  'Dave Public Bookmark 1',
  'This is a public bookmark from Dave',
  'PUBLIC',
  'test development',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'dave@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'dave-private-1',
  'https://example.com/dave/private1',
  'Dave Private Bookmark 1',
  'This is a private bookmark from Dave',
  'PRIVATE',
  'test work',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'dave@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'dave-friends-1',
  'https://example.com/dave/friends1',
  'Dave Friends Bookmark 1',
  'This is a friends-only bookmark from Dave',
  'FRIENDS',
  'test friends-only',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'dave@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'dave-public-2',
  'https://example.com/dave/public2',
  'Dave Public Bookmark 2',
  'Another public bookmark from Dave',
  'PUBLIC',
  'test tutorial',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'dave@test.com';

-- Create test bookmarks for Eve
INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'eve-public-1',
  'https://example.com/eve/public1',
  'Eve Public Bookmark 1',
  'This is a public bookmark from Eve',
  'PUBLIC',
  'test science',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'eve@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'eve-private-1',
  'https://example.com/eve/private1',
  'Eve Private Bookmark 1',
  'This is a private bookmark from Eve',
  'PRIVATE',
  'test research',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'eve@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'eve-friends-1',
  'https://example.com/eve/friends1',
  'Eve Friends Bookmark 1',
  'This is a friends-only bookmark from Eve',
  'FRIENDS',
  'test friends-only',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'eve@test.com';
