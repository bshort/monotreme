-- Backfill user_order for existing shortcuts based on created_ts (newest first)
-- This ensures shortcuts created before the user_order feature get proper sequential values

-- Create a temporary table with row numbers based on creation date
WITH numbered_shortcuts AS (
  SELECT
    id,
    ROW_NUMBER() OVER (ORDER BY created_ts DESC) - 1 AS new_order
  FROM shortcut
)
UPDATE shortcut
SET user_order = (
  SELECT new_order
  FROM numbered_shortcuts
  WHERE numbered_shortcuts.id = shortcut.id
);
