# Monotreme CLI Management Tool

A comprehensive command-line interface for managing your Monotreme application.

## Quick Start

Run the CLI from the project directory:

```bash
./monotreme
```

## Features

### 🔍 Status Monitoring
- Real-time container status (Running/Stopped/Not Created)
- Container uptime information
- Database statistics:
  - Number of shortcuts
  - Number of users
  - Number of collections
- Application URL display

### 🚀 Service Management

#### Start All Services
Starts both Monotreme and Redis containers using Docker Compose.

#### Stop All Services
Gracefully stops all containers and cleans up resources.

#### Restart All Services
Restarts all running services without rebuilding.

#### Rebuild and Start
Rebuilds the Docker image (using `rebuild_and_start.sh` if available) and starts services.

### ⚙️ Individual Service Management

Manage Monotreme and Redis containers independently:
- Start individual services
- Stop individual services
- Restart individual services

### 📋 Log Viewing

View real-time logs with options:
- Monotreme logs only
- Redis logs only
- All logs combined

Press `Ctrl+C` to exit log view.

## Menu Structure

```
MAIN MENU:
  1) Show Status
  2) Start All Services
  3) Stop All Services
  4) Restart All Services
  5) Rebuild and Start
  6) Manage Individual Services
  7) View Logs
  0) Exit

SERVICE MANAGEMENT:
  1) Start/Stop/Restart Monotreme
  2) Start/Stop/Restart Redis
  0) Back to Main Menu

LOGS VIEWER:
  1) Monotreme Logs
  2) Redis Logs
  3) All Logs
  0) Back to Main Menu
```

## Requirements

- Docker and Docker Compose installed and running
- `sqlite3` command-line tool (optional, for database statistics)
  - Install on Ubuntu/Debian: `sudo apt-get install sqlite3`

## Color-Coded Output

- 🟢 **Green**: Running/Success
- 🟡 **Yellow**: Stopped/Warning
- 🔴 **Red**: Error/Not Created
- 🔵 **Blue**: Information
- 🔷 **Cyan**: Menus/Headers

## Usage Examples

### Check Status
1. Run `./monotreme`
2. Select option `1` to view current status

### Start Services
1. Run `./monotreme`
2. Select option `2` to start all services
3. Services will start and status will be displayed

### View Logs
1. Run `./monotreme`
2. Select option `7` for logs menu
3. Choose which logs to view
4. Press `Ctrl+C` to return to menu

### Restart Individual Service
1. Run `./monotreme`
2. Select option `6` for service management
3. Choose the service to restart
4. Select restart option

## Troubleshooting

### Docker Not Running
If you see "Docker is not running or not installed", ensure:
- Docker daemon is running: `sudo systemctl start docker`
- Your user has Docker permissions: `sudo usermod -aG docker $USER`

### Database Statistics Unavailable
Install sqlite3:
```bash
sudo apt-get install sqlite3
```

### Permission Denied
Make the script executable:
```bash
chmod +x monotreme
```

## Tips

- The CLI automatically detects the current state of all services
- Menu options change dynamically based on service status
- All operations provide feedback and confirmation
- Use option `0` at any menu level to go back or exit
