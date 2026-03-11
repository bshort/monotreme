# Monotreme Dashboard Guide

## Overview

The Monotreme Dashboard is a persistent, real-time TUI (Terminal User Interface) that provides live monitoring and management of your Monotreme application.

## Features

### 🔄 Real-Time Monitoring
- **Auto-refresh every 2 seconds** - Status updates automatically
- **Detects external changes** - See updates even when docker-compose is run outside the dashboard
- **Non-blocking interface** - Execute commands while monitoring continues
- **Live resource metrics** - CPU and memory usage for each container

### 📊 Dashboard Sections

1. **System Status**
   - Docker daemon status
   - Container states (Running/Stopped/Not Created)
   - Container uptime
   - CPU and memory usage per container

2. **Database Statistics**
   - Number of shortcuts
   - Number of users
   - Number of collections
   - Database file size

3. **Application Access**
   - Main application URL
   - Health check endpoint
   - API documentation URL
   (Only shown when Monotreme is running)

4. **Last Action**
   - Shows the most recent command executed
   - Timestamp of execution
   - Success/failure status

5. **Controls**
   - Context-sensitive commands
   - Change based on current system state

## Quick Start

```bash
# Launch the dashboard
./monotreme

# The dashboard will:
# 1. Display current status
# 2. Auto-refresh every 2 seconds
# 3. Wait for your commands
```

## Keyboard Commands

### Service Management

| Key | Action | Description |
|-----|--------|-------------|
| `s` | Start All | Start Monotreme and Redis containers |
| `x` | Stop All | Stop all containers |
| `r` | Restart All | Restart all containers |
| `b` | Build & Start | Rebuild images and start services |

### Individual Container Control

| Key | Action | Description |
|-----|--------|-------------|
| `m` | Toggle Monotreme | Start/stop Monotreme container |
| `n` | Restart Monotreme | Restart Monotreme container |
| `d` | Toggle Redis | Start/stop Redis container |
| `e` | Restart Redis | Restart Redis container |

### Monitoring & Utilities

| Key | Action | Description |
|-----|--------|-------------|
| `l` | Monotreme Logs | View real-time logs (Ctrl+C to return) |
| `k` | Redis Logs | View real-time logs (Ctrl+C to return) |
| `v` | Vacuum Database | Optimize database (VACUUM command) |
| `c` | Clear Last Action | Clear the last action display |
| `q` | Quit | Exit the dashboard |

## Dashboard Behavior

### Persistent Display
- The dashboard never clears between commands
- Actions are executed without reloading the entire interface
- Status updates happen in place

### Real-Time Updates
The dashboard polls the system every 2 seconds to update:
- Container status changes
- Uptime information
- Resource usage (CPU/Memory)
- Database statistics

This means changes made outside the dashboard (e.g., running `docker compose up -d` in another terminal) will be reflected automatically.

### Non-Blocking Actions
When you execute a command:
1. The action starts immediately
2. "Last Action" section shows "Running..."
3. Dashboard continues polling and updating
4. Status updates to ✓ Success or ✗ Failed when complete
5. You can execute other commands while actions run

### Context-Sensitive Controls
The available commands change based on the current state:

**When nothing is running:**
- Start All Services
- Build & Start

**When everything is running:**
- Stop All Services
- Restart All
- Individual container controls

**When partially running:**
- Start All Services
- Stop All Services
- Start/stop individual containers

## Example Workflow

### Starting from Scratch
```
1. Launch dashboard: ./monotreme
2. Status shows: ○ NOT CREATED
3. Press 's' to start all services
4. Watch containers transition to ● RUNNING
5. Application URLs appear automatically
```

### Daily Monitoring
```
1. Launch dashboard: ./monotreme
2. See current uptime and resource usage
3. Check database statistics
4. Leave running to monitor continuously
5. Press 'q' when done
```

### Troubleshooting
```
1. Launch dashboard: ./monotreme
2. Notice Monotreme is stopped
3. Press 'l' to view logs
4. Press Ctrl+C to return to dashboard
5. Press 'n' to restart Monotreme
6. Watch status update automatically
```

### Working with External Commands
```
Terminal 1:
1. Launch dashboard: ./monotreme
2. See current status

Terminal 2:
1. Run: docker compose down
2. Switch to Terminal 1
3. Dashboard automatically shows: ○ NOT CREATED

Terminal 2:
1. Run: docker compose up -d
2. Switch to Terminal 1
3. Dashboard automatically shows: ● RUNNING
```

## Status Indicators

| Symbol | Color | Meaning |
|--------|-------|---------|
| ● | Green | Running |
| ◐ | Yellow | Stopped |
| ○ | Red | Not Created |
| ✓ | Green | Docker Available |
| ✗ | Red | Docker Not Available |

## Resource Metrics

The dashboard shows real-time resource usage:
- **CPU**: Percentage of CPU used by container
- **Memory**: Current/Total memory used

Example: `CPU: 2.5%  │  Memory: 45MiB / 1GiB`

## Action Status Codes

| Status | Meaning |
|--------|---------|
| `Running...` | Command is executing |
| `✓ Success` | Command completed successfully |
| `✗ Failed` | Command failed (check logs) |
| `✓ Build Complete` | Rebuild finished successfully |
| `✓ Database Vacuumed` | Database optimization complete |

## Tips & Tricks

### Tip 1: Leave It Running
The dashboard is designed to stay open. Leave it running in a terminal window to monitor your application continuously.

### Tip 2: Quick Status Check
Just want to check if services are running? Launch the dashboard, glance at the status indicators, press 'q' to exit. Takes 2 seconds.

### Tip 3: Log Monitoring
Press 'l' or 'k' to jump into logs, then Ctrl+C to return to the dashboard. No need to quit and restart.

### Tip 4: Clear Clutter
Finished with an action? Press 'c' to clear the "Last Action" section and clean up the display.

### Tip 5: Emergency Stop
Services acting up? Press 'x' to stop everything immediately. The dashboard continues running so you can restart when ready.

## Troubleshooting

### Dashboard Not Updating
- Check that Docker is running: `docker info`
- Verify containers exist: `docker ps -a`
- Ensure terminal supports ANSI colors

### Database Stats Show 0
- Install sqlite3: `sudo apt-get install sqlite3`
- Verify database file exists: `ls -lh data/monotreme_prod.db`

### Actions Fail Silently
- Check action log: `cat /tmp/monotreme_action.log`
- Verify Docker permissions: `docker ps`
- Ensure you're in project directory

### Terminal Garbled After Exit
- Run: `reset` to restore terminal
- Or: `clear` to just clear the screen

## Comparison with Old Menu System

| Feature | Old Menu | New Dashboard |
|---------|----------|---------------|
| Interface | Reloads each time | Persistent display |
| Updates | Manual refresh | Auto-refresh (2s) |
| External changes | Not detected | Detected automatically |
| Actions | Blocking | Non-blocking |
| Resource metrics | None | Live CPU/Memory |
| Navigation | Multi-level menus | Single-key commands |
| Status | Static | Real-time |

## Advanced Usage

### Custom Refresh Interval
Edit the script and change:
```bash
POLL_INTERVAL=2  # Change to desired seconds
```

### Running in tmux/screen
The dashboard works great in tmux or screen:
```bash
# In tmux
tmux new-session -s monotreme
./monotreme

# Detach: Ctrl+B, then D
# Reattach: tmux attach -t monotreme
```

### SSH Remote Monitoring
Monitor remote Monotreme installations:
```bash
ssh user@remote-server
cd /path/to/monotreme
./monotreme
```

## Requirements

- Bash 4.0+
- Docker & Docker Compose
- Terminal with ANSI color support
- sqlite3 (optional, for database stats)

## Exit

Press `q` to exit cleanly. The dashboard will:
1. Restore your cursor
2. Clear the screen
3. Show "Dashboard closed. Goodbye!"

---

**Remember**: The dashboard is a monitor AND controller. You can manage your entire Monotreme installation without ever leaving it!
