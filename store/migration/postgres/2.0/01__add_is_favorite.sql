-- Add is_favorite column to shortcut table
ALTER TABLE shortcut ADD COLUMN is_favorite BOOLEAN NOT NULL DEFAULT FALSE;
