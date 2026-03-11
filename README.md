# Monotreme

**Monotreme** is an open source, self-hosted platform designed to help you organize, manage, and share your most important links. Easily create customizable, human-readable shortcuts to streamline your link management. Use tags to categorize your links and share them easily with your team or publicly.

Monotreme is a fork of [Slash](https://github.com/monotrememarks/slash) and was created to allow extensions to Slash's core model to fit one user's bizarre requirements. Monotreme is probably not ready for widespread use and if you're looking for a self-hosted bookmarking application then just use [Slash](https://github.com/monotrememarks/slash). It's better!


## Background

In today's workplace, essential information is often scattered across the cloud in the form of links. We understand the frustration of endlessly searching through emails, messages, and websites just to find the right link. Links are notorious for being unwieldy, complex, and easily lost in the shuffle. Remembering and sharing them can be a challenge.

That's why we developed Monotreme (forked from [Slash](https://github.com/monotrememarks/slash)), a solution that transforms these links into easily accessible, discoverable, and shareable shortcuts(e.g., `s/shortcut`). Say goodbye to link chaos and welcome the organizational ease of Monotreme into your daily online workflow.

## Features

- Create customizable `s/` short links for any URL.
- Share short links public or only with your teammates.
- View analytics on link traffic and sources.
- Easy access to your shortcuts with browser extension.
- Share your shortcuts with Collection to anyone, on any browser.
- Open source self-hosted solution.

## Deploy with Docker in seconds

```bash
docker run -d --name monotreme -p 5231:5231 -v ~/.monotreme/:/var/opt/monotreme bshort/monotreme:latest
```

Learn more in [Self-hosting Monotreme with Docker](https://github.com/bshort/monotreme/blob/main/docs/install.md).

---

## Getting Started

### 0. Prerequisites

Before running Monotreme, ensure you have the following installed:

#### Required
- **Docker** (version 20.10 or higher)
  ```bash
  # Check Docker version
  docker --version
  ```
- **Docker Compose** (version 2.0 or higher)
  ```bash
  # Check Docker Compose version
  docker compose version
  ```

#### Optional
- **sqlite3** (for database statistics and backups)
  ```bash
  # Install on Ubuntu/Debian
  sudo apt-get install sqlite3

  # Install on macOS
  brew install sqlite3
  ```
- **Node.js and npm** (only if building from source)
  - Node.js 18.x or higher
  - npm 9.x or higher

#### System Requirements
- **Memory**: Minimum 512MB RAM, 1GB+ recommended
- **Disk Space**: Minimum 500MB for application and database
- **Port Availability**: Port 5231 (Monotreme) and 6379 (Redis) must be available

#### How Monotreme Runs
Monotreme consists of two Docker containers orchestrated by Docker Compose:
- **monotreme**: The main application server (Go backend + React frontend)
  - Runs on port 5231
  - Stores data in SQLite database at `./data/monotreme_prod.db`
  - Supports Redis caching for improved performance
- **monotreme-redis**: Redis cache server
  - Runs on port 6379
  - Caches user shortcuts and frequently accessed data
  - Can be disabled via environment variables

The application data is persisted in the `./data` directory on your host machine, ensuring your shortcuts, users, and collections survive container restarts.

---

### 1. Starting the Application

#### Option A: Using the Monotreme CLI (Recommended)

The easiest way to manage Monotreme is using the built-in CLI tool:

```bash
# Make the CLI executable (first time only)
chmod +x monotreme

# Launch the interactive menu
./monotreme
```

From the menu, select option **2) Start All Services**.

The CLI will:
- Start both Monotreme and Redis containers
- Show you the application status
- Display the access URL (http://localhost:5231)

**Alternative CLI commands:**
```bash
# Quick status check
./monotreme  # Then select option 1

# View in menu for all management options
./monotreme
```

#### Option B: Using Docker Compose Directly

```bash
# Start services in detached mode
docker compose up -d

# Check if services are running
docker compose ps

# View logs
docker compose logs -f
```

#### Option C: Using the Startup Script

```bash
# Start without rebuilding
bash start.sh

# Or rebuild and start (includes npm install and fresh build)
bash rebuild_and_start.sh
```

#### First-Time Setup

After starting for the first time:

1. **Access the application**: Open http://localhost:5231 in your browser
2. **Create admin account**: You'll be prompted to create the first user account
3. **Configure settings**: Set your workspace preferences in Settings
4. **Create your first shortcut**: Click "New Shortcut" to get started

#### Verification

Confirm everything is running:
```bash
# Check container status
docker ps | grep monotreme

# Check application logs
docker compose logs monotreme | tail -20

# Test the endpoint
curl http://localhost:5231/healthz
# Should return: "Service ready."
```

---

### 2. Maintaining the Application

#### Using the Monotreme CLI

The CLI provides comprehensive maintenance options:

```bash
./monotreme
```

**Common maintenance tasks:**
- **Option 1**: Show Status - Monitor system health and database statistics
- **Option 4**: Restart All Services - Apply configuration changes
- **Option 5**: Rebuild and Start - Update to latest code changes
- **Option 6**: Manage Individual Services - Restart specific components
- **Option 7**: View Logs - Troubleshoot issues

#### Monitoring

**Check System Status:**
```bash
# Using CLI (recommended)
./monotreme  # Select option 1

# Or manually
docker compose ps
docker stats monotreme monotreme-redis
```

**View Logs:**
```bash
# Using CLI
./monotreme  # Select option 7

# Or manually - View all logs
docker compose logs -f

# View Monotreme logs only
docker compose logs -f monotreme

# View Redis logs only
docker compose logs -f monotreme-redis

# View last 100 lines
docker compose logs --tail=100 monotreme
```

**Database Statistics:**
```bash
# Using CLI (shows automatically in status)
./monotreme  # Select option 1

# Or query manually
sqlite3 data/monotreme_prod.db "SELECT COUNT(*) as shortcuts FROM shortcut WHERE row_status='NORMAL';"
sqlite3 data/monotreme_prod.db "SELECT COUNT(*) as users FROM user WHERE row_status='NORMAL';"
sqlite3 data/monotreme_prod.db "SELECT COUNT(*) as collections FROM collection;"
```

#### Updating Monotreme

**From Source Code:**
```bash
# Pull latest changes
git pull origin main

# Rebuild and restart
./monotreme  # Select option 5 (Rebuild and Start)

# Or use script
bash rebuild_and_start.sh
```

**From Docker Image:**
```bash
# Pull latest image
docker compose pull

# Restart with new image
docker compose up -d
```

#### Performance Tuning

**Redis Cache Configuration:**

Edit `.env` or `docker-compose.yml`:
```env
# Enable/disable Redis caching
REDIS_ENABLED=true

# Redis connection details
REDIS_HOST=monotreme-redis
REDIS_PORT=6379
```

**Database Optimization:**
```bash
# Vacuum database to reclaim space
docker compose exec monotreme sqlite3 /var/opt/monotreme/monotreme_prod.db "VACUUM;"

# Check database size
ls -lh data/monotreme_prod.db
```

#### Troubleshooting

**Common Issues:**

1. **Port already in use:**
   ```bash
   # Check what's using port 5231
   sudo lsof -i :5231

   # Change port in docker-compose.yml
   # Change "5231:5231" to "8080:5231"
   ```

2. **Database locked errors:**
   ```bash
   # Stop all services
   ./monotreme  # Option 3

   # Remove lock files
   rm -f data/monotreme_prod.db-shm data/monotreme_prod.db-wal

   # Restart
   ./monotreme  # Option 2
   ```

3. **Redis connection errors:**
   ```bash
   # Check Redis is running
   docker compose ps monotreme-redis

   # Restart Redis
   ./monotreme  # Option 6, then restart Redis
   ```

4. **Application won't start:**
   ```bash
   # Check logs for errors
   docker compose logs monotreme

   # Rebuild from scratch
   docker compose down
   docker compose build --no-cache
   docker compose up -d
   ```

---

### 3. Backing Up the Application

#### Full Backup Strategy

Monotreme stores all data in the `./data` directory. A complete backup includes:
- SQLite database (`monotreme_prod.db`)
- Database write-ahead log files (`.db-shm`, `.db-wal`)
- Environment configuration (`.env`)
- Docker configuration (`docker-compose.yml`)

#### Automated Backup Script

Create a backup script (`backup.sh`):

```bash
#!/bin/bash
# Save this as backup.sh and make executable: chmod +x backup.sh

BACKUP_DIR="$HOME/monotreme-backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_NAME="monotreme_backup_${TIMESTAMP}"

# Create backup directory
mkdir -p "$BACKUP_DIR/$BACKUP_NAME"

# Stop containers to ensure clean backup
echo "Stopping Monotreme..."
docker compose down

# Backup database files
echo "Backing up database..."
cp -r data/ "$BACKUP_DIR/$BACKUP_NAME/"

# Backup configuration files
echo "Backing up configuration..."
cp docker-compose.yml "$BACKUP_DIR/$BACKUP_NAME/"
[ -f .env ] && cp .env "$BACKUP_DIR/$BACKUP_NAME/"

# Create compressed archive
echo "Creating archive..."
cd "$BACKUP_DIR"
tar -czf "${BACKUP_NAME}.tar.gz" "$BACKUP_NAME"
rm -rf "$BACKUP_NAME"

# Restart containers
echo "Restarting Monotreme..."
cd - > /dev/null
docker compose up -d

echo "Backup complete: $BACKUP_DIR/${BACKUP_NAME}.tar.gz"
```

#### Manual Backup Methods

**Quick Database Backup (Hot Backup):**
```bash
# Backup without stopping services
sqlite3 data/monotreme_prod.db ".backup data/monotreme_backup_$(date +%Y%m%d).db"
```

**Cold Backup (Safer, requires downtime):**
```bash
# Stop services
./monotreme  # Option 3 (Stop All Services)

# Create backup
mkdir -p ~/monotreme-backups
cp -r data ~/monotreme-backups/monotreme_data_$(date +%Y%m%d_%H%M%S)

# Restart services
./monotreme  # Option 2 (Start All Services)
```

**Backup to Remote Server:**
```bash
# Using rsync
rsync -avz data/ user@backup-server:/backups/monotreme/

# Using scp
tar -czf - data/ | ssh user@backup-server "cat > /backups/monotreme_$(date +%Y%m%d).tar.gz"
```

#### Automated Scheduled Backups

**Using Cron (Linux/macOS):**

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /home/bshort/Projects/monotreme/backup.sh

# Add weekly backup on Sunday at 3 AM
0 3 * * 0 /home/bshort/Projects/monotreme/backup.sh
```

**Backup Retention Script:**

```bash
#!/bin/bash
# Keep only last 7 daily backups and 4 weekly backups

BACKUP_DIR="$HOME/monotreme-backups"

# Delete backups older than 7 days
find "$BACKUP_DIR" -name "monotreme_backup_*.tar.gz" -mtime +7 -delete

# Keep only 4 most recent backups
ls -t "$BACKUP_DIR"/monotreme_backup_*.tar.gz | tail -n +5 | xargs -r rm
```

#### Restoring from Backup

**Restore Database:**

```bash
# Stop services
./monotreme  # Option 3

# Backup current data (just in case)
mv data data.old

# Extract backup
tar -xzf ~/monotreme-backups/monotreme_backup_YYYYMMDD_HHMMSS.tar.gz
cp -r monotreme_backup_YYYYMMDD_HHMMSS/data ./

# Restart services
./monotreme  # Option 2
```

**Restore Specific Database:**

```bash
# Stop services
docker compose down

# Restore database file
cp ~/monotreme-backups/monotreme_backup_20231016.db data/monotreme_prod.db

# Remove WAL files (will be recreated)
rm -f data/monotreme_prod.db-shm data/monotreme_prod.db-wal

# Start services
docker compose up -d
```

#### Testing Backups

Regularly verify your backups work:

```bash
# Test database integrity
sqlite3 ~/monotreme-backups/monotreme_backup_20231016.db "PRAGMA integrity_check;"

# Should output: ok

# Test restoration in separate directory
mkdir -p /tmp/monotreme-test
cd /tmp/monotreme-test
tar -xzf ~/monotreme-backups/monotreme_backup_YYYYMMDD_HHMMSS.tar.gz
sqlite3 monotreme_backup_*/data/monotreme_prod.db ".tables"
```

---

### 4. Stopping the Application

#### Using the Monotreme CLI (Recommended)

```bash
./monotreme
```

Select **Option 3) Stop All Services**

This will:
- Gracefully stop both Monotreme and Redis containers
- Clean up network resources
- Preserve all data in the `./data` directory

#### Using Docker Compose

**Stop and remove containers:**
```bash
# Stop containers (removes containers but keeps data)
docker compose down
```

**Stop without removing containers:**
```bash
# Just stop services (containers still exist)
docker compose stop
```

**Stop individual service:**
```bash
# Stop only Monotreme
docker compose stop monotreme

# Stop only Redis
docker compose stop monotreme-redis
```

#### Using the Stop Script

```bash
bash stop.sh
```

#### Force Stop (Emergency)

If containers won't stop gracefully:

```bash
# Force kill containers
docker compose kill

# Clean up
docker compose down
```

#### Complete Cleanup

**Remove containers and volumes:**
```bash
# WARNING: This removes the database!
# Backup first!
docker compose down -v

# Remove images as well
docker compose down --rmi all
```

**Remove everything (including data):**
```bash
# WARNING: DESTRUCTIVE - Backup first!

# Stop and remove containers
docker compose down

# Remove data directory
rm -rf data/

# Remove Docker images
docker rmi monotreme:local
docker rmi redis:7-alpine
```

#### Verification

Confirm services are stopped:

```bash
# Check no containers running
docker ps | grep monotreme
# Should return nothing

# Verify with CLI
./monotreme  # Option 1 - Should show "NOT RUNNING"

# Check port is free
sudo lsof -i :5231
# Should return nothing
```

#### Temporary Stop (Maintenance Mode)

For brief maintenance:

```bash
# Stop services
docker compose stop

# ... perform maintenance ...

# Resume services (quick restart)
docker compose start
```

---

## Additional Resources

- **CLI Documentation**: See [MONOTREME_CLI.md](./MONOTREME_CLI.md) for detailed CLI usage
- **Installation Guide**: [docs/install.md](./docs/install.md)
- **Contributing**: See [CONTRIBUTING.md](./CONTRIBUTING.md)
- **Troubleshooting**: Check [GitHub Issues](https://github.com/bshort/monotreme/issues)

