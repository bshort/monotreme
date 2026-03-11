# Docker Deployment Guide

This project supports two different networking configurations:

## Local Development (Bridge Networking)

For local development on your computer, use bridge networking:

```bash
./start.sh
```

This will:
- Use `docker-compose.yml` (base configuration with bridge networking)
- Use `docker-compose.override.yml` (local overrides)
- Mount data to `./data`
- Expose services at `localhost:5231` and `localhost:6379`

## Production Deployment (MacVlan Networking)

For production deployment with MacVlan networking:

```bash
./start-prod.sh
```

This will:
- Use `docker-compose.yml` (base configuration)
- Use `docker-compose.prod.yml` (production overrides with MacVlan)
- Mount data to `/Primer/Media/data/monotreme`
- Assign static IPs:
  - monotreme: 192.168.1.136
  - monotreme-redis: 192.168.1.137

### Prerequisites for Production

Before running in production, ensure:
1. The MacVlan network `MyMacVlan` exists
2. The IP addresses (192.168.1.136 and 192.168.1.137) are available
3. The data directory `/Primer/Media/data/monotreme` exists

### Common Commands

```bash
# View logs (local)
docker compose logs -f

# View logs (production)
docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f

# Stop containers (local)
docker compose down

# Stop containers (production)
docker compose -f docker-compose.yml -f docker-compose.prod.yml down
```

## File Structure

- `docker-compose.yml` - Base configuration with bridge networking
- `docker-compose.override.yml` - Local development overrides (auto-loaded)
- `docker-compose.prod.yml` - Production overrides with MacVlan
- `start.sh` - Local development startup script
- `start-prod.sh` - Production deployment script
