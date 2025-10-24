-- Add user_order column to shortcut table for custom ordering
ALTER TABLE shortcut ADD COLUMN user_order INTEGER NOT NULL DEFAULT 0;
