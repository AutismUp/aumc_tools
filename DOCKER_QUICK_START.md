# Docker Test Environment - Quick Start Guide

## Prerequisites

✅ Colima installed and running  
✅ Docker and Docker Compose installed

## Quick Commands

### Start Everything
```bash
./docker-test.sh
```
This will:
1. Build the Linux binary
2. Build the Docker image
3. Start the container
4. Run all tests
5. Leave container running for interactive use

### Access Container Shell
```bash
docker-compose exec minecraft-test bash
```

### Run AUMC Commands
```bash
# From your Mac (outside container)
docker-compose exec minecraft-test aumc --help
docker-compose exec minecraft-test aumc init

# Or inside the container
docker-compose exec minecraft-test bash
# Now you're inside the container
aumc --help
aumc init
```

### Stop Container
```bash
docker-compose down
```

### Restart Container
```bash
docker-compose restart
```

### View Logs
```bash
docker-compose logs -f
```

### Rebuild Everything
```bash
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d
```

## Common Tasks

### Test a Code Change
```bash
# 1. Make your changes to the Go code
# 2. Run the test script
./docker-test.sh
```

### Create and Test Config
```bash
docker-compose exec minecraft-test bash -c "
  mkdir -p /tmp/mytest
  cd /tmp/mytest
  aumc init
  cat config.json
"
```

### Check MSM Status
```bash
docker-compose exec minecraft-test sudo msm version
docker-compose exec minecraft-test sudo msm jargroup list
```

### Verify Java Installation
```bash
docker-compose exec minecraft-test java -version
```

## Colima Management

### Start Colima
```bash
colima start --cpu 2 --memory 4
```

### Stop Colima
```bash
colima stop
```

### Check Colima Status
```bash
colima status
```

### Restart Colima
```bash
colima restart
```

## Troubleshooting

### "Cannot connect to Docker daemon"
```bash
colima start
```

### Container won't start
```bash
docker-compose down
docker-compose up -d
docker-compose logs
```

### Need to rebuild image
```bash
docker-compose build --no-cache
```

### Clean everything
```bash
docker-compose down -v
docker system prune -a
```

## File Locations

### On Your Mac
- Source code: `/Users/nicholashatch/dev/aumc_tools`
- Binary: `bin/aumc-linux-amd64`
- Docker files: `Dockerfile`, `docker-compose.yml`

### Inside Container
- Workspace (synced): `/workspace`
- Binary installed: `/usr/local/bin/aumc`
- MSM: `/opt/msm`
- BuildTools: `/opt/build_tools/BuildTools.jar`

## Port Forwarding

- **25565** - Minecraft server port (forwarded to host)

## Users in Container

- **testuser** - Main user (has sudo, in minecraft group)
- **minecraft** - MSM service user
- **root** - System admin

## Tips

1. The `/workspace` directory in the container is synced with your project directory
2. Changes to Go code require rebuilding the binary: `make build-linux`
3. The container persists data in a Docker volume: `minecraft-data`
4. Use `docker-compose exec` to run commands without entering the container
5. Use `docker-compose exec -it` for interactive commands

## Example Workflow

```bash
# 1. Start environment
./docker-test.sh

# 2. Make code changes in your editor
# (edit cmd/root.go or other files)

# 3. Rebuild and test
make build-linux
docker-compose exec minecraft-test sudo cp /workspace/bin/aumc-linux-amd64 /usr/local/bin/aumc
docker-compose exec minecraft-test aumc --help

# 4. Run specific tests
docker-compose exec minecraft-test bash /workspace/test-aumc.sh

# 5. When done
docker-compose down
```
