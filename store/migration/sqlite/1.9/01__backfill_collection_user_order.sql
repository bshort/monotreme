-- Backfill user_order for existing collections based on created_ts (newest first)
-- This ensures collections created before the user_order feature get proper sequential values

-- Create a temporary table with row numbers based on creation date
WITH numbered_collections AS (
  SELECT
    id,
    ROW_NUMBER() OVER (ORDER BY created_ts DESC) - 1 AS new_order
  FROM collection
)
UPDATE collection
SET user_order = (
  SELECT new_order
  FROM numbered_collections
  WHERE numbered_collections.id = collection.id
);
