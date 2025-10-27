-- Create WORKSPACE (friends-only) bookmarks for test users
-- Note: Using WORKSPACE visibility since FRIENDS is not in the schema yet

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'alice-workspace-1',
  'https://example.com/alice/workspace1',
  'Alice Workspace Bookmark 1',
  'This is a workspace-only bookmark from Alice',
  'WORKSPACE',
  'test workspace-only',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'alice@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'bob-workspace-1',
  'https://example.com/bob/workspace1',
  'Bob Workspace Bookmark 1',
  'This is a workspace-only bookmark from Bob',
  'WORKSPACE',
  'test workspace-only',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'bob@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'carol-workspace-1',
  'https://example.com/carol/workspace1',
  'Carol Workspace Bookmark 1',
  'This is a workspace-only bookmark from Carol',
  'WORKSPACE',
  'test workspace-only',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'carol@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'dave-workspace-1',
  'https://example.com/dave/workspace1',
  'Dave Workspace Bookmark 1',
  'This is a workspace-only bookmark from Dave',
  'WORKSPACE',
  'test workspace-only',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'dave@test.com';

INSERT INTO shortcut (creator_id, name, link, title, description, visibility, tag, created_ts, updated_ts)
SELECT
  id,
  'eve-workspace-1',
  'https://example.com/eve/workspace1',
  'Eve Workspace Bookmark 1',
  'This is a workspace-only bookmark from Eve',
  'WORKSPACE',
  'test workspace-only',
  strftime('%s', 'now'),
  strftime('%s', 'now')
FROM user WHERE email = 'eve@test.com';
