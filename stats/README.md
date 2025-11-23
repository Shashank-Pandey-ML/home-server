# Stats Service

A lightweight Go microservice that collects and reports system statistics including CPU, memory, disk, and network metrics.

## Features

- Real-time CPU usage and core-level statistics
- Memory usage and availability
- Disk usage for physical filesystems
- Network interface information
- System uptime and load averages
- Host system integration (works in Docker with proper mounts)

## API Endpoints

### `GET /stats`
Returns comprehensive system statistics in JSON format.

**Response Example:**
```json
{
  "timestamp": "2025-11-23T10:30:00Z",
  "hostname": "server-01",
  "platform": "darwin",
  "os": "darwin",
  "architecture": "arm64",
  "cpu_count": 8,
  "cpu": {
    "usage_percent": 45.2,
    "per_core_usage": [42.1, 48.3, ...],
    "model_name": "Apple M1"
  },
  "memory": {
    "total_gb": 16.0,
    "used_gb": 10.2,
    "available_gb": 5.8,
    "used_percent": 63.75
  },
  ...
}
```

### `GET /health`
Health check endpoint.

## Docker Configuration

The service requires access to host system information when running in Docker:

```yaml
volumes:
  - /proc:/host/proc:ro
  - /sys:/host/sys:ro
  - /etc/os-release:/host/etc/os-release:ro
  - /etc/hostname:/host/etc/hostname:ro
environment:
  - HOST_PROC=/host/proc
  - HOST_SYS=/host/sys
  - HOST_ETC=/host/etc
```

## Security

- All host mounts are read-only (`:ro`)
- No root filesystem access
- Only system information directories exposed
- Minimal attack surface

## Building

```bash
make build-stats
```

## Running Locally

```bash
cd stats-service
go run app/main.go
```

Service starts on port 8080 by default.
