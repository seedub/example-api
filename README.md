# example-api

API backend to support the React frontend ([example-ui](https://github.com/seedub/example-ui))

## Overview

This is a lightweight REST API built with Go that provides backend services for the example-ui React application. The API is designed to be simple, fast, and easy to deploy.

## Features

- **Simple HTTP API**: RESTful endpoints for basic operations
- **CORS Support**: Configured to allow cross-origin requests from the UI
- **Health Check**: Built-in health check endpoint for monitoring
- **Comprehensive Testing**: 76% code coverage with unit tests
- **CI/CD Pipeline**: Automated testing, linting, building, and deployment
- **AWS Deployment**: Automated deployment to EC2 instances

## API Endpoints

### Root Endpoint
- **GET** `/` - Returns a simple text message for the UI
  ```bash
  curl http://localhost:8080/
  # Response: Hello from the example API!
  ```

### Health Check
- **GET** `/health` - Returns API health status
  ```bash
  curl http://localhost:8080/health
  # Response: {"status":"healthy","time":"2026-02-09T21:00:00Z"}
  ```

### Items API
- **GET** `/api/items` - Get all items
- **POST** `/api/items` - Create a new item
- **GET** `/api/items/:id` - Get a specific item
- **PUT** `/api/items/:id` - Update an item
- **DELETE** `/api/items/:id` - Delete an item

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/seedub/example-api.git
   cd example-api
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   go test ./... -v -cover
   ```

4. Build the application:
   ```bash
   # Build for current platform
   go build -o bin/example-api ./main.go
   
   # Or use Makefile for cross-compilation
   make build-amd64  # Build for x86_64/amd64
   make build-arm64  # Build for ARM64
   make build-all    # Build for all architectures
   ```

5. Run the application:
   ```bash
   ./bin/example-api
   # Or with custom port:
   PORT=3000 ./bin/example-api
   ```

The API will start on port 8080 by default (or the port specified in the `PORT` environment variable).

## Development

### Running Tests

```bash
# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Project Structure

```
.
├── api/                    # API handlers and routing
│   ├── handlers.go         # HTTP request handlers
│   ├── handlers_test.go    # Handler tests
│   ├── middleware.go       # Middleware (CORS, logging)
│   ├── middleware_test.go  # Middleware tests
│   ├── router.go           # Route configuration
│   └── router_test.go      # Router tests
├── models/                 # Data models
│   ├── item.go            # Item model
│   └── item_test.go       # Model tests
├── deployment/            # Deployment configurations
│   └── example-api.service # Systemd service file
├── .github/workflows/     # CI/CD workflows
│   └── ci-cd.yml          # GitHub Actions pipeline
├── main.go                # Application entry point
└── go.mod                 # Go module file
```

## Deployment

### CI/CD Pipeline

The project uses GitHub Actions for continuous integration and deployment:

1. **Test**: Runs all tests with race detection and coverage reporting
2. **Lint**: Checks code quality with golangci-lint
3. **Build**: Compiles the application binary
4. **Deploy**: Deploys to AWS EC2 (only on main branch)

### AWS Deployment

The application is automatically deployed to AWS EC2 when changes are pushed to the `main` branch.

The build process supports both **amd64** (x86_64) and **arm64** (ARM/Graviton) architectures. The deployment workflow automatically detects the EC2 instance architecture and deploys the correct binary.

#### Supported Architectures

- **amd64** (x86_64): Standard Intel/AMD processors
- **arm64** (aarch64): AWS Graviton processors

#### Required GitHub Secrets

- `SSH_KEY`: SSH private key for EC2 access

#### Deployment Target

- **Host**: 35.90.6.81 (us-west-2)
- **Architecture**: ARM64 (aarch64)
- **User**: ec2-user
- **Port**: 8080
- **Shared with**: [example-ui](https://github.com/seedub/example-ui) (Nginx on port 80)

#### Manual Deployment

If you need to deploy manually:

```bash
# Build for your target architecture
make build-amd64  # For x86_64 instances
# OR
make build-arm64  # For ARM/Graviton instances

# Copy to EC2 (adjust for your architecture)
scp -i /path/to/key.pem bin/example-api-amd64 ec2-user@35.90.6.81:/tmp/example-api

# SSH and install
ssh -i /path/to/key.pem ec2-user@35.90.6.81
sudo mv /tmp/example-api /usr/local/bin/example-api
sudo chmod +x /usr/local/bin/example-api
sudo systemctl restart example-api
```

### Systemd Service

The API runs as a systemd service on the EC2 instance. The service file is located in `deployment/example-api.service`.

```bash
# Check service status
sudo systemctl status example-api

# View logs
sudo journalctl -u example-api -f

# Restart service
sudo systemctl restart example-api
```

## Environment Variables

- `PORT`: Server port (default: 8080)

## Deployment Architecture (Same Server as UI)

This API runs on the same AWS EC2 instance (35.90.6.81, ARM64) as the [example-ui](https://github.com/seedub/example-ui) React application.

### Server Configuration

**API (this repository):**
- Go HTTP server running directly on port 8080
- No nginx needed - Go's built-in HTTP server is production-ready
- Managed by systemd service

**UI (example-ui repository):**
- Nginx on port 80 serving static React files
- Location: `/var/www/example-ui`

### UI Configuration for Same-Server Deployment

The example-ui needs to be configured to call the API on the same server. Update `src/App.jsx`:

```javascript
// For same-server deployment on port 8080
const API_BASE_URL = 'http://35.90.6.81:8080'

// OR use window location if you want it dynamic
const API_BASE_URL = `http://${window.location.hostname}:8080`
```

**Note:** Port 8080 must be opened in the AWS Security Group for external access.

### Alternative: Nginx Reverse Proxy (Optional)

If you prefer not to expose port 8080, you can configure Nginx to proxy API requests:

```nginx
# In /etc/nginx/conf.d/example-ui.conf
location /api/ {
    proxy_pass http://localhost:8080/api/;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}

location /health {
    proxy_pass http://localhost:8080/health;
}
```

Then the UI can use relative URLs or same-origin requests:
```javascript
const API_BASE_URL = ''  // Same origin
```

### Testing the Integration

```bash
# Test API directly
curl http://35.90.6.81:8080/api/items
curl http://35.90.6.81:8080/health

# Test from the UI
open http://35.90.6.81
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Testing

This project maintains high code coverage (>70% required by CI). All new features should include appropriate tests.

## License

This project is open source and available under the MIT License.

## Related Projects

- [example-ui](https://github.com/seedub/example-ui) - React frontend for this API 
