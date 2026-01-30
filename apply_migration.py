#!/usr/bin/env python3
"""
Apply the is_favorite migration to the production database.
This script adds the is_favorite column to the shortcut table.
"""
import sqlite3
import sys
from pathlib import Path

DB_PATH = "data/monotreme_prod.db"
MIGRATION_VERSION = "2.0"

def check_column_exists(cursor, table_name, column_name):
    """Check if a column exists in a table."""
    cursor.execute(f"PRAGMA table_info({table_name})")
    columns = [row[1] for row in cursor.fetchall()]
    return column_name in columns

def check_migration_applied(cursor, version):
    """Check if a migration version has been applied."""
    try:
        cursor.execute("SELECT version FROM migration_history WHERE version = ?", (version,))
        return cursor.fetchone() is not None
    except sqlite3.OperationalError:
        # migration_history table might not exist yet
        return False

def apply_migration(db_path):
    """Apply the is_favorite migration to the database."""
    if not Path(db_path).exists():
        print(f"Error: Database not found at {db_path}")
        return False

    print(f"Connecting to database: {db_path}")
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    try:
        # Check if is_favorite column already exists
        if check_column_exists(cursor, "shortcut", "is_favorite"):
            print("✓ is_favorite column already exists in shortcut table")

            # Check if user_order column exists (migration 02)
            if not check_column_exists(cursor, "shortcut", "user_order"):
                print("✗ user_order column is missing. This is unexpected.")
                print("  The schema might be in an inconsistent state.")
                return False
            else:
                print("✓ user_order column exists in shortcut table")

            return True

        print("Applying migration 2.0/01__add_is_favorite.sql...")

        # Apply the migration - add is_favorite column
        cursor.execute("ALTER TABLE shortcut ADD COLUMN is_favorite BOOLEAN NOT NULL DEFAULT false")
        print("✓ Added is_favorite column")

        # Record migration in history
        cursor.execute(
            "INSERT INTO migration_history (version) VALUES (?)",
            (f"{MIGRATION_VERSION}/01__add_is_favorite.sql",)
        )
        print("✓ Recorded migration in history")

        # Check if user_order needs backfilling (migration 02)
        if check_column_exists(cursor, "shortcut", "user_order"):
            print("\nApplying migration 2.0/02__backfill_shortcut_user_order.sql...")

            # Backfill user_order for existing shortcuts
            cursor.execute("""
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
                )
            """)
            print(f"✓ Backfilled user_order for {cursor.rowcount} shortcuts")

            # Record migration in history
            cursor.execute(
                "INSERT INTO migration_history (version) VALUES (?)",
                (f"{MIGRATION_VERSION}/02__backfill_shortcut_user_order.sql",)
            )
            print("✓ Recorded backfill migration in history")
        else:
            print("⚠ user_order column doesn't exist - skipping backfill")

        # Commit all changes
        conn.commit()
        print("\n✓ Migration applied successfully!")

        # Verify the changes
        print("\nVerifying changes...")
        cursor.execute("PRAGMA table_info(shortcut)")
        columns = [row[1] for row in cursor.fetchall()]
        print(f"Shortcut table now has {len(columns)} columns:")
        print(f"  Columns: {', '.join(columns)}")

        return True

    except Exception as e:
        print(f"\n✗ Error applying migration: {e}")
        conn.rollback()
        return False
    finally:
        conn.close()

def main():
    print("="*60)
    print("Monotreme Database Migration Tool")
    print("Migration: Add is_favorite column to shortcut table")
    print("="*60)
    print()

    success = apply_migration(DB_PATH)

    if success:
        print("\n" + "="*60)
        print("SUCCESS: Migration completed!")
        print("="*60)
        print("\nNext steps:")
        print("1. Restart your Monotreme application")
        print("2. The error should be resolved")
        sys.exit(0)
    else:
        print("\n" + "="*60)
        print("FAILED: Migration could not be applied")
        print("="*60)
        print("\nTroubleshooting:")
        print("1. Make sure the database file exists")
        print("2. Check that you have write permissions")
        print("3. Verify the database is not locked by another process")
        sys.exit(1)

if __name__ == "__main__":
    main()
