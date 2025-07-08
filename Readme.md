# Docker Cleaner

A CLI tool to clean up Docker resources and reclaim disk space.

## Features

- 🧹 Clean stopped containers
- 🖼️ Remove unused images (dangling by default)
- 💾 Clean unused volumes
- 🌐 Remove unused networks
- 📊 Show cleanup statistics
- 📋 List all Docker resources
- 🔧 Flexible command-line options

## Installation

### Download Pre-built Binaries

Download the latest release from the [releases page](https://github.com/yourusername/docker-cleaner/releases).

### Install from Source

```bash
go install github.com/yourusername/docker-cleaner@latest
```

### Build from Source

```bash
git clone https://github.com/yourusername/docker-cleaner.git
cd docker-cleaner
make build
```

## Usage

### Basic Usage

```bash
# Clean stopped containers and dangling images (default)
docker-cleaner clean

# Clean all resources with confirmation
docker-cleaner clean --all

# Clean specific resources
docker-cleaner clean --containers --images --volumes

# Force clean without confirmation
docker-cleaner clean --all --force
```

### List Resources

```bash
# List all Docker resources
docker-cleaner list
```

### Version Information

```bash
# Show version information
docker-cleaner version
```

## Commands

### `clean`

Clean Docker resources to free up disk space.

**Flags:**
- `-a, --all`: Clean all resources (containers, images, volumes, networks)
- `-c, --containers`: Clean containers only
- `-i, --images`: Clean images only
- `-v, --volumes`: Clean volumes only
- `-n, --networks`: Clean networks only
- `-f, --force`: Force removal without confirmation
- `--dangling`: Only remove dangling images (default: true)

**Examples:**

```bash
# Clean stopped containers and dangling images
docker-cleaner clean

# Clean all resources
docker-cleaner clean --all

# Clean only containers
docker-cleaner clean --containers

# Clean images and volumes without confirmation
docker-cleaner clean --images --volumes --force

# Clean all images (not just dangling)
docker-cleaner clean --images --dangling=false
```

### `list`

List Docker resources with detailed information.

**Example:**

```bash
docker-cleaner list
```

### `version`

Show version information.

**Example:**

```bash
docker-cleaner version
```

## Requirements

- Docker installed and running
- Docker socket accessible (usually `/var/run/docker.sock`)
- Appropriate permissions to manage Docker resources

## Docker Usage

You can also run docker-cleaner as a Docker container:

```bash
# Build the image
docker build -t docker-cleaner .

# Run the container (mount Docker socket)
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock docker-cleaner clean --all
```

## Development

### Prerequisites

- Go 1.21 or later
- Docker for testing

### Setup

```bash
git clone https://github.com/yourusername/docker-cleaner.git
cd docker-cleaner
make deps
```

### Available Make Targets

```bash
make help          # Show all available targets
make build         # Build the binary
make build-all     # Build for multiple platforms
make test          # Run tests
make lint          # Run linter
make fmt           # Format code
make clean         # Clean build artifacts
make install       # Install the binary
make release       # Create release packages
```

### Testing

```bash
make test
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Safety

This tool modifies Docker resources. Always review what will be removed before running cleanup commands, especially with the `--force` flag.

## Troubleshooting

### Permission Denied

If you get permission denied errors:

```bash
# Add your user to the docker group
sudo usermod -aG docker $USER

# Or run with sudo
sudo docker-cleaner clean
```

### Docker Socket Not Found

If Docker socket is not found:

```bash
# Check if Docker is running
docker info

# Check socket location
ls -la /var/run/docker.sock
```

## Roadmap

- [ ] Add dry-run mode
- [ ] Support for custom Docker socket paths
- [ ] Interactive mode for selective cleanup
- [ ] Configuration file support
- [ ] Backup before cleanup
- [ ] Scheduled cleanup
- [ ] Docker Compose integration