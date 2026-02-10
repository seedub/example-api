# GitHub Copilot Instructions for example-api

## Project Summary

This is a lightweight REST API built with Go that provides backend services for the example-ui React application. The API is designed to be simple, fast, and easy to deploy to AWS EC2 instances with automated HTTPS/SSL support via nginx reverse proxy.

## Codebase Overview

The project is a Go HTTP server that provides RESTful endpoints for managing items and health checks. It uses:
- **Go 1.24** - Programming language
- **Standard library** - HTTP server (no external frameworks)
- **Systemd** - Service management on Linux
- **Nginx** - Reverse proxy with SSL/TLS termination
- **GitHub Actions** - CI/CD pipeline

## Project Structure

```
.
├── api/                    # API handlers and routing
│   ├── handlers.go         # HTTP request handlers for items and health
│   ├── handlers_test.go    # Handler tests
│   ├── middleware.go       # Middleware (CORS, logging)
│   ├── middleware_test.go  # Middleware tests
│   ├── router.go           # Route configuration
│   └── router_test.go      # Router tests
├── models/                 # Data models
│   ├── item.go            # Item model definition
│   └── item_test.go       # Model tests
├── deployment/            # Deployment configurations
│   ├── example-api.service      # Systemd service file
│   ├── nginx-example-api.conf   # Nginx reverse proxy config
│   ├── setup-letsencrypt.sh     # SSL/TLS setup script
│   └── validate-https.sh        # HTTPS validation script
├── .github/workflows/     # CI/CD workflows
│   └── ci-cd.yml          # GitHub Actions pipeline
├── main.go                # Application entry point
├── go.mod                 # Go module file
└── README.md              # Project documentation
```

## API Endpoints

- `GET /` - Root endpoint returning a simple message
- `GET /health` - Health check endpoint with timestamp
- `GET /api/items` - Get all items
- `POST /api/items` - Create a new item
- `GET /api/items/:id` - Get a specific item
- `PUT /api/items/:id` - Update an item
- `DELETE /api/items/:id` - Delete an item

## Development Guidelines

### Building the Code

```bash
# Download dependencies
go mod download

# Build for current platform
go build -o bin/example-api ./main.go

# Or use Makefile for cross-compilation
make build-amd64  # Build for x86_64/amd64
make build-arm64  # Build for ARM64
make build-all    # Build for all architectures
```

### Running Tests

```bash
# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run tests with race detection (required for CI)
go test ./... -v -race -coverprofile=coverage.out -covermode=atomic

# Generate coverage report
go tool cover -func=coverage.out
go tool cover -html=coverage.out
```

### Code Quality Requirements

- **Coverage threshold**: Minimum 70% code coverage (enforced by CI)
- **Linting**: All code must pass `golangci-lint` checks
- **Testing**: All new features must include unit tests
- **Race detection**: All tests must pass with `-race` flag

### Formatting and Linting

```bash
# Format code (Go standard)
go fmt ./...

# Run linter (required before merge)
golangci-lint run

# Fix auto-fixable linting issues
golangci-lint run --fix
```

### Running the Application

```bash
# Default port (8080)
./bin/example-api

# Custom port
PORT=3000 ./bin/example-api

# Run with go run for development
go run main.go
```

## CI/CD Pipeline

The GitHub Actions workflow (`.github/workflows/ci-cd.yml`) runs on push and pull requests:

1. **Test**: Runs all tests with race detection and coverage reporting
2. **Lint**: Checks code quality with golangci-lint
3. **Build**: Compiles binaries for amd64 and arm64 architectures
4. **Deploy**: Deploys to AWS EC2 (only on main branch pushes)

## Deployment Architecture

- **API Server**: Go HTTP server on localhost:8080 (not exposed externally)
- **Nginx**: Reverse proxy on port 443 (HTTPS) and port 80 (HTTP redirect)
- **SSL/TLS**: Let's Encrypt certificates for production, self-signed for development
- **Service Management**: Systemd service for automatic start/restart

### Required GitHub Secrets for Deployment

- `SSH_KEY`: AWS EC2 SSH private key (PEM format)
- `EC2_HOST`: EC2 instance IP or hostname
- `EC2_USER`: SSH username (e.g., ec2-user, ubuntu)

## Key Technical Principles

1. **Simplicity**: Use Go standard library when possible, minimize external dependencies
2. **Test Coverage**: Maintain >70% code coverage for all code
3. **Cross-Platform**: Support both amd64 and arm64 architectures
4. **Security**: All production traffic must use HTTPS/TLS
5. **CORS**: Configured to support cross-origin requests from frontend
6. **Error Handling**: Return appropriate HTTP status codes and error messages
7. **Logging**: Log all requests and important operations
8. **Health Checks**: Provide `/health` endpoint for monitoring

## Contribution Requirements

Before submitting a pull request:

1. ✅ All tests pass: `go test ./... -v -race`
2. ✅ Coverage meets threshold: `go test ./... -cover` (>70%)
3. ✅ Code is formatted: `go fmt ./...`
4. ✅ Linting passes: `golangci-lint run`
5. ✅ No race conditions detected
6. ✅ New features include appropriate tests
7. ✅ Documentation is updated if needed

## Environment Variables

- `PORT`: Server port (default: 8080)

## Related Projects

- [example-ui](https://github.com/seedub/example-ui) - React frontend for this API
