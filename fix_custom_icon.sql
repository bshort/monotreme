-- Manual fix for missing custom_icon columns
-- Run this SQL script directly on your database to add the missing custom_icon columns

-- Add custom_icon column to shortcut table if it doesn't exist
ALTER TABLE shortcut ADD COLUMN custom_icon TEXT NOT NULL DEFAULT '';

-- Add custom_icon column to collection table if it doesn't exist
ALTER TABLE collection ADD COLUMN custom_icon TEXT NOT NULL DEFAULT '';

-- Verify the columns were added
PRAGMA table_info(shortcut);
PRAGMA table_info(collection);