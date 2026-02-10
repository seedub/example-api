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
  # Local development (HTTP)
  curl http://localhost:8080/
  
  # Production with HTTPS
  curl https://YOUR_SERVER_IP/
  # Response: Hello from the example API!
  ```

### Health Check
- **GET** `/health` - Returns API health status
  ```bash
  # Local development (HTTP)
  curl http://localhost:8080/health
  
  # Production with HTTPS
  curl https://YOUR_SERVER_IP/health
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

To enable automated deployment to AWS EC2, configure the following secrets in your GitHub repository (Settings → Secrets and variables → Actions → New repository secret):

- **`SSH_KEY`**: AWS EC2 SSH private key for deploying to the EC2 instance
  - This is the private SSH key (PEM format) used to authenticate with the AWS EC2 instance
  - Value: The complete private key content (including `-----BEGIN RSA PRIVATE KEY-----` and `-----END RSA PRIVATE KEY-----` headers)

- **`EC2_HOST`**: The IP address or hostname of your EC2 instance
  - Example: `35.90.6.81` or `ec2-xx-xx-xx-xx.compute.amazonaws.com`
  
- **`EC2_USER`**: The SSH username for your EC2 instance  
  - Typically `ec2-user` for Amazon Linux, `ubuntu` for Ubuntu, etc.

#### Deployment Target

The deployment is configured via GitHub Secrets for security. Example configuration:

- **Host**: Your EC2 instance IP or hostname (configured in `EC2_HOST` secret)
- **Architecture**: The deployment automatically detects the EC2 architecture (amd64 or arm64)
- **User**: Your EC2 SSH user (configured in `EC2_USER` secret)
- **Port**: 8080 (configured in systemd service)
- **Shared with**: Can be shared with frontend applications like [example-ui](https://github.com/seedub/example-ui) (Nginx on port 80)

#### Manual Deployment

If you need to deploy manually:

```bash
# Build for your target architecture
make build-amd64  # For x86_64 instances
# OR
make build-arm64  # For ARM/Graviton instances

# Copy to EC2 (replace with your actual values)
scp -i /path/to/key.pem bin/example-api-amd64 YOUR_EC2_USER@YOUR_EC2_HOST:/tmp/example-api

# SSH and install
ssh -i /path/to/key.pem YOUR_EC2_USER@YOUR_EC2_HOST
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

## SSL/TLS Certificate Management

### Initial Setup

The deployment includes automated SSL/TLS setup using the `deployment/setup-letsencrypt.sh` script:

```bash
# After deployment, SSH to your server and run:
sudo /tmp/setup-letsencrypt.sh
```

### Certificate Types

**For IP-based Access (Default):**
- Generates a self-signed certificate
- Works immediately for testing
- Browsers will show security warnings
- Suitable for development and internal use

**For Domain-based Access (Production):**
```bash
# After running setup-letsencrypt.sh, obtain a real Let's Encrypt certificate:
sudo certbot certonly --standalone -d your-domain.com

# Update nginx configuration with your domain:
sudo nano /etc/nginx/conf.d/example-api.conf
# Change: server_name _; to server_name your-domain.com;

# Update certificate paths to use your domain:
# ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
# ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

# Test and reload nginx
sudo nginx -t
sudo systemctl reload nginx
```

### Certificate Renewal

Certificates are automatically renewed via cron job:
```bash
# Check renewal status
sudo certbot renew --dry-run

# Manual renewal (if needed)
sudo certbot renew

# Reload nginx after renewal
sudo systemctl reload nginx
```

### Certificate Verification

```bash
# Quick validation of entire HTTPS setup
sudo /tmp/validate-https.sh

# Check certificate expiration
echo | openssl s_client -servername YOUR_DOMAIN -connect YOUR_SERVER:443 2>/dev/null | openssl x509 -noout -dates

# Verify HTTPS is working
curl -I https://YOUR_SERVER/health

# For self-signed certificates, use -k flag
curl -Ik https://YOUR_SERVER/health
```

**For complete HTTPS setup documentation, see [HTTPS_SETUP.md](HTTPS_SETUP.md)**

## Environment Variables

- `PORT`: Server port (default: 8080)

## Deployment Architecture with HTTPS

This API is designed to run behind nginx with HTTPS support using Let's Encrypt SSL certificates.

### Server Configuration

**API (this repository):**
- Go HTTP server running on localhost port 8080
- Not directly exposed to the internet
- Managed by systemd service

**Nginx (reverse proxy):**
- Handles HTTPS/SSL termination on port 443
- Redirects HTTP (port 80) to HTTPS
- Proxies requests to the Go API on localhost:8080
- Configured via `deployment/nginx-example-api.conf`

**SSL/TLS:**
- Let's Encrypt certificates (for domain-based deployments)
- Self-signed certificates (for IP-based deployments)
- Automatic HTTP to HTTPS redirect

### Setting Up HTTPS

After deploying the API, set up HTTPS with the provided script:

```bash
# SSH to your EC2 instance
ssh -i your-key.pem ec2-user@YOUR_EC2_HOST

# Run the SSL setup script
sudo /tmp/setup-letsencrypt.sh
```

This script will:
1. Install nginx and certbot
2. Generate SSL certificates (self-signed for IP-based access)
3. Configure nginx with the provided configuration
4. Set up automatic certificate renewal

**Note:** For production use with a domain name, update the certificate after running the script:
```bash
sudo certbot certonly --standalone -d your-domain.com
# Then update /etc/nginx/conf.d/example-api.conf with your domain
```

### UI Configuration for Same-Server Deployment

The frontend application needs to be configured to call the API on the same server. Example for React:

```javascript
// For HTTPS with nginx reverse proxy (recommended)
const API_BASE_URL = `https://${window.location.hostname}`

// The API will be available at:
// - https://your-server/api/items
// - https://your-server/health
// - https://your-server/
```

**Security Group / Firewall Requirements:**
- Port 443 (HTTPS) must be open for external access
- Port 80 (HTTP) must be open for Let's Encrypt verification and redirect to HTTPS
- Port 8080 should NOT be exposed externally (only localhost access needed)

### Nginx Reverse Proxy Configuration

The nginx configuration (`deployment/nginx-example-api.conf`) provides:

- **HTTPS/SSL termination** on port 443
- **HTTP to HTTPS redirect** on port 80
- **Reverse proxy** to Go API on localhost:8080
- **Security headers** (HSTS, X-Frame-Options, etc.)
- **CORS headers** for cross-origin requests
- **Let's Encrypt support** for certificate renewal

### Testing the Integration

```bash
# Test API with HTTPS (replace YOUR_SERVER_IP with your server)
curl https://YOUR_SERVER_IP/api/items
curl https://YOUR_SERVER_IP/health

# For self-signed certificates, use -k flag
curl -k https://YOUR_SERVER_IP/api/items

# Test HTTP to HTTPS redirect
curl -I http://YOUR_SERVER_IP/health
# Should return 301 redirect to https://

# Test from the UI
open https://YOUR_EC2_HOST
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
