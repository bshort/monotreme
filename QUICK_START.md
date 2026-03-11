# Monotreme Quick Start Guide

## TL;DR - Get Started in 30 Seconds

```bash
# 1. Start the application
./monotreme

# 2. Select option 2 (Start All Services)

# 3. Open http://localhost:5231
```

---

## Common Commands

### Using the CLI (Easiest)

```bash
./monotreme    # Opens interactive menu
```

**Quick Actions:**
- Press `1` - Check status
- Press `2` - Start services
- Press `3` - Stop services
- Press `4` - Restart services
- Press `7` - View logs

### Using Docker Compose

```bash
# Start
docker compose up -d

# Stop
docker compose down

# View logs
docker compose logs -f

# Restart
docker compose restart
```

### Using Scripts

```bash
bash start.sh              # Start services
bash stop.sh               # Stop services
bash rebuild_and_start.sh  # Rebuild and start
```

---

## Quick Backup

```bash
# Hot backup (no downtime)
sqlite3 data/monotreme_prod.db ".backup data/backup_$(date +%Y%m%d).db"

# Cold backup (safer)
./monotreme  # Stop services (option 3)
cp -r data ~/monotreme-backups/backup_$(date +%Y%m%d)
./monotreme  # Start services (option 2)
```

---

## Troubleshooting One-Liners

```bash
# Service won't start?
docker compose down && docker compose up -d

# Database locked?
rm -f data/*.db-shm data/*.db-wal && docker compose restart

# Port conflict?
sudo lsof -i :5231

# Check logs
docker compose logs --tail=50 monotreme

# Full reset (WARNING: Loses data!)
docker compose down && rm -rf data/ && docker compose up -d
```

---

## Database Quick Queries

```bash
# Count shortcuts
sqlite3 data/monotreme_prod.db "SELECT COUNT(*) FROM shortcut WHERE row_status='NORMAL';"

# Count users
sqlite3 data/monotreme_prod.db "SELECT COUNT(*) FROM user WHERE row_status='NORMAL';"

# Count collections
sqlite3 data/monotreme_prod.db "SELECT COUNT(*) FROM collection;"

# Check database size
ls -lh data/monotreme_prod.db
```

---

## First Time Setup Checklist

- [ ] Install Docker and Docker Compose
- [ ] Clone repository or download release
- [ ] Run `chmod +x monotreme` (first time only)
- [ ] Run `./monotreme` and select option 2
- [ ] Open http://localhost:5231
- [ ] Create admin account
- [ ] Create first shortcut

---

## System Requirements

- **Docker**: 20.10+
- **RAM**: 512MB minimum, 1GB recommended
- **Disk**: 500MB minimum
- **Ports**: 5231 (Monotreme), 6379 (Redis)

---

## Need Help?

- **Full Documentation**: [README.md](./README.md)
- **CLI Guide**: [MONOTREME_CLI.md](./MONOTREME_CLI.md)
- **Issues**: https://github.com/bshort/monotreme/issues

---

## Daily Operations

### Morning Startup
```bash
./monotreme  # Option 2 (Start)
```

### End of Day Shutdown
```bash
./monotreme  # Option 3 (Stop)
```

### Weekly Maintenance
```bash
# Check status and stats
./monotreme  # Option 1

# Backup database
sqlite3 data/monotreme_prod.db ".backup data/weekly_backup.db"

# View logs for any issues
./monotreme  # Option 7
```

### Monthly Tasks
```bash
# Update to latest version
git pull origin main
./monotreme  # Option 5 (Rebuild and Start)

# Vacuum database
docker compose exec monotreme sqlite3 /var/opt/monotreme/monotreme_prod.db "VACUUM;"

# Clean old backups
find ~/monotreme-backups -mtime +30 -delete
```

---

## Configuration Files

- `docker-compose.yml` - Container orchestration
- `.env` - Environment variables (optional)
- `data/monotreme_prod.db` - Main database
- `monotreme` - CLI management tool

---

## URLs

- **Application**: http://localhost:5231
- **Health Check**: http://localhost:5231/healthz
- **API Docs**: http://localhost:5231/swagger

---

## Emergency Procedures

### Application Completely Broken
```bash
# 1. Stop everything
docker compose down

# 2. Backup data
cp -r data data.backup

# 3. Clean rebuild
docker compose build --no-cache
docker compose up -d

# 4. Check logs
docker compose logs -f
```

### Need to Restore from Backup
```bash
# 1. Stop services
docker compose down

# 2. Restore database
cp data.backup/monotreme_prod.db data/

# 3. Clean WAL files
rm -f data/*.db-shm data/*.db-wal

# 4. Restart
docker compose up -d
```

### Lost Admin Password
```bash
# Use the database to reset (requires sqlite3)
docker compose exec monotreme sqlite3 /var/opt/monotreme/monotreme_prod.db
# Then manually update user table (see full docs)
```

---

**Remember**: Always backup before major changes!
