#!/usr/bin/env python3
"""Check the current database schema and verify columns."""
import sqlite3

DB_PATH = "data/monotreme_prod.db"

def main():
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()

    print("="*60)
    print("Database Schema Check")
    print("="*60)
    print()

    # Get shortcut table schema
    print("SHORTCUT TABLE SCHEMA:")
    print("-" * 60)
    cursor.execute("PRAGMA table_info(shortcut)")
    columns = cursor.fetchall()
    for col in columns:
        cid, name, type_, notnull, default, pk = col
        print(f"{cid:2d}. {name:20s} {type_:15s} "
              f"{'NOT NULL' if notnull else 'NULL':8s} "
              f"DEFAULT: {default if default else 'None':10s} "
              f"{'PK' if pk else ''}")

    print()
    print("="*60)
    print("MIGRATION HISTORY:")
    print("-" * 60)
    try:
        cursor.execute("SELECT version, created_ts FROM migration_history ORDER BY created_ts")
        migrations = cursor.fetchall()
        if migrations:
            for version, created_ts in migrations:
                print(f"  {version} (applied at: {created_ts})")
        else:
            print("  No migrations recorded")
    except sqlite3.OperationalError as e:
        print(f"  Error reading migration_history: {e}")

    print()
    print("="*60)
    print("SAMPLE QUERY TEST:")
    print("-" * 60)
    try:
        cursor.execute("SELECT id, name, is_favorite FROM shortcut LIMIT 3")
        shortcuts = cursor.fetchall()
        print(f"Successfully queried {len(shortcuts)} shortcuts")
        for shortcut_id, name, is_favorite in shortcuts:
            print(f"  ID: {shortcut_id}, Name: {name}, is_favorite: {is_favorite}")
    except sqlite3.OperationalError as e:
        print(f"✗ Error querying shortcuts: {e}")

    print()
    print("="*60)
    print("WAL MODE CHECK:")
    print("-" * 60)
    cursor.execute("PRAGMA journal_mode")
    journal_mode = cursor.fetchone()[0]
    print(f"Journal mode: {journal_mode}")

    if journal_mode.upper() == "WAL":
        print("\nWAL mode is active. Checkpointing WAL...")
        cursor.execute("PRAGMA wal_checkpoint(FULL)")
        result = cursor.fetchone()
        print(f"Checkpoint result: {result}")

    conn.close()
    print("\n" + "="*60)

if __name__ == "__main__":
    main()
