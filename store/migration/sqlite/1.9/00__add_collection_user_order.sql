-- Add user_order column to collection table for custom ordering
ALTER TABLE collection ADD COLUMN user_order INTEGER NOT NULL DEFAULT 0;
