-- Add custom_icon column to shortcut and collection tables
-- This migration attempts to add the columns, handling the case where they might already exist

-- Add custom_icon column to shortcut table
ALTER TABLE shortcut ADD COLUMN custom_icon TEXT NOT NULL DEFAULT '';

-- Add custom_icon column to collection table
ALTER TABLE collection ADD COLUMN custom_icon TEXT NOT NULL DEFAULT '';