# Docker Cleaner 🧹

A powerful CLI tool to clean up Docker resources and free up disk space. Built with Go and designed for simplicity and efficiency.

## Features

- 🐳 **Clean Containers**: Remove stopped containers or all containers
- 🖼️ **Clean Images**: Remove dangling images or all unused images
- 💾 **Clean Volumes**: Remove unused volumes
- 🌐 **Clean Networks**: Remove unused networks
- 📊 **Detailed Reports**: See exactly what was cleaned and how much space was reclaimed
- 🚀 **Fast & Efficient**: Built with Go for maximum performance
- 🔒 **Safe**: Asks for confirmation before destructive operations

## Installation

### Download Binary

Download the latest binary for your platform from the [releases page](https://github.com/zahidhasann88/docker-cleaner/releases).

#### Linux/macOS
```bash
# Download and install
curl -L https://github.com/zahidhasann88/docker-cleaner/releases/latest/download/docker-cleaner-linux-amd64 -o docker-cleaner
chmod +x docker-cleaner
sudo mv docker-cleaner /usr/local/bin/
```

#### Windows
```powershell
# Download from releases page or use PowerShell
Invoke-WebRequest -Uri "https://github.com/zahidhasann88/docker-cleaner/releases/latest/download/docker-cleaner-windows-amd64.exe" -OutFile "docker-cleaner.exe"
```

### Docker

```bash
# Run directly with Docker
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock zahidhasann88/docker-cleaner:latest

# Or use GitHub Container Registry
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock ghcr.io/zahidhasann88/docker-cleaner:latest
```

### Go Install

```bash
go install github.com/zahidhasann88/docker-cleaner@latest
```

## Usage

### List Docker Resources
```bash
docker-cleaner list
```

### Clean Specific Resources
```bash
# Clean only containers
docker-cleaner clean --containers

# Clean only images
docker-cleaner clean --images

# Clean only volumes
docker-cleaner clean --volumes

# Clean only networks
docker-cleaner clean --networks
```

### Clean All Resources
```bash
# Clean everything (with confirmation)
docker-cleaner clean --all

# Clean everything (force, no confirmation)
docker-cleaner clean --all --force
```

### Advanced Options
```bash
# Clean only dangling images (default)
docker-cleaner clean --images --dangling

# Clean all unused images
docker-cleaner clean --images --dangling=false

# Force removal of all containers (running and stopped)
docker-cleaner clean --containers --force
```

## Examples

### Basic Cleanup
```bash
$ docker-cleaner clean --containers --images
This will remove Docker resources. Continue? (y/N)
y
🧹 Cleaning containers...
   ✓ Removed 3 containers
🖼️  Cleaning images...
   ✓ Removed 5 images

📊 Cleanup Summary:
   Containers: 3
   Images: 5
   Volumes: 0
   Networks: 0
   Space reclaimed: 1.2 GB
```

### List Resources
```bash
$ docker-cleaner list
📋 Docker Resources Overview
==================================================

🐳 Containers (2 total):
   ID           Image                Status          Names
   ------------------------------------------------------------
   1b9d5c227153 mysql:8.0            Up 29 seconds   /db-test-1
   a1b2c3d4e5f6 nginx:latest         Exited (0)      /web-server

🖼️  Images (7 total):
   ID           Repository                     Tag        Size
   ----------------------------------------------------------------------
   sha256:0c211 mysql                          8.0        736.1 MB
   sha256:15390 nginx                          latest     142.8 MB
   sha256:2c232 <none>                         <none>     736.1 MB
```

## Commands

| Command | Description |
|---------|-------------|
| `clean` | Clean Docker resources |
| `list` | List all Docker resources |
| `version` | Show version information |
| `help` | Show help for any command |

## Flags

| Flag | Description |
|------|-------------|
| `--all, -a` | Clean all resources |
| `--containers, -c` | Clean containers |
| `--images, -i` | Clean images |
| `--volumes, -v` | Clean volumes |
| `--networks, -n` | Clean networks |
| `--force, -f` | Force removal without confirmation |
| `--dangling` | Only remove dangling images (default: true) |

## Docker Usage

### Run with Docker
```bash
# Basic usage
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock zahidhasann88/docker-cleaner list

# Clean all resources
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock zahidhasann88/docker-cleaner clean --all --force
```

### Docker Compose
```yaml
version: '3.8'
services:
  docker-cleaner:
    image: zahidhasann88/docker-cleaner:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: clean --all --force
```

## Building from Source

### Prerequisites
- Go 1.21 or later
- Docker (for Docker image)

### Build
```bash
# Clone the repository
git clone https://github.com/zahidhasann88/docker-cleaner.git
cd docker-cleaner

# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `make test` and `make lint`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 🐛 [Report bugs](https://github.com/zahidhasann88/docker-cleaner/issues)
- 💡 [Request features](https://github.com/zahidhasann88/docker-cleaner/issues)
- 📖 [Documentation](https://github.com/zahidhasann88/docker-cleaner/wiki)

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for details about each release.

---

**⚠️ Warning**: This tool can remove Docker resources permanently. Always review what will be removed before confirming the operation, especially in production environments.