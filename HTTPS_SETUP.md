# HTTPS Setup Guide

This guide explains how to set up HTTPS for the example-api using nginx and Let's Encrypt SSL certificates.

## Prerequisites

- A deployed example-api instance on EC2 or similar server
- SSH access to the server
- Root/sudo privileges
- Ports 80, 443 open in your firewall/security group

## Quick Setup

After deploying the API using the CI/CD pipeline, HTTPS setup is automated:

```bash
# SSH to your server
ssh -i your-key.pem ec2-user@YOUR_SERVER_IP

# Run the automated setup script
sudo /tmp/setup-letsencrypt.sh
```

This script will:
1. Install nginx and certbot (if not already installed)
2. Generate SSL certificates (self-signed for IP-based access)
3. Deploy the nginx configuration automatically
4. Test the nginx configuration
5. Enable and start nginx
6. Set up automatic certificate renewal
7. Verify HTTPS is working

**No manual steps required!** The script is fully automated.

## Setup Details

### 1. Install Dependencies

The setup script automatically installs:
- **nginx**: Web server and reverse proxy
- **certbot**: Let's Encrypt certificate manager

### 2. SSL Certificate Options

#### Option A: IP-based Access (Default)

For servers accessed by IP address, the script generates a self-signed certificate:

```bash
sudo /tmp/setup-letsencrypt.sh

# Validate the setup
sudo /tmp/validate-https.sh
```

**Pros:**
- Works immediately
- No domain name required
- Good for development/testing

**Cons:**
- Browsers show security warnings
- Not recommended for production

#### Option B: Domain-based Access (Recommended for Production)

For production use with a domain name:

```bash
# 1. First run the basic setup
sudo /tmp/setup-letsencrypt.sh

# 2. Point your domain to your server's IP
#    (Configure your DNS A record)

# 3. Get a real Let's Encrypt certificate
sudo certbot certonly --standalone -d your-domain.com

# 4. Update nginx configuration
sudo nano /etc/nginx/conf.d/example-api.conf
#    Change: server_name _; 
#    To:     server_name your-domain.com;

# 5. Update certificate paths in the same file
#    Change: /etc/letsencrypt/live/example-api/...
#    To:     /etc/letsencrypt/live/your-domain.com/...

# 6. Test and reload nginx
sudo nginx -t
sudo systemctl reload nginx
```

### 3. Nginx Configuration

The deployment includes a pre-configured nginx setup at `/etc/nginx/conf.d/example-api.conf` with:

- **HTTP (port 80)**: Redirects to HTTPS
- **HTTPS (port 443)**: Serves the API with SSL/TLS
- **Reverse Proxy**: Forwards requests to Go API on localhost:8080
- **Security Headers**: HSTS, X-Frame-Options, X-Content-Type-Options, etc.
- **CORS Support**: Cross-origin request headers

### 4. Certificate Renewal

Certificates are automatically renewed via cron job (set up by the script).

#### Check Renewal Status
```bash
sudo certbot renew --dry-run
```

#### Manual Renewal
```bash
sudo certbot renew
sudo systemctl reload nginx
```

#### View Certificate Expiration
```bash
echo | openssl s_client -servername YOUR_DOMAIN -connect YOUR_SERVER:443 2>/dev/null | openssl x509 -noout -dates
```

## Architecture

```
┌─────────────────┐
│   User/Client   │
└────────┬────────┘
         │ HTTPS (port 443)
         ▼
┌─────────────────┐
│  Nginx (TLS)    │  ← SSL/TLS Termination
│  Port 80, 443   │  ← HTTP → HTTPS Redirect
└────────┬────────┘
         │ HTTP (localhost:8080)
         ▼
┌─────────────────┐
│  Go API Server  │
│  Port 8080      │
└─────────────────┘
```

## Quick Validation

After setup, run the validation script to verify everything is working:

```bash
sudo /tmp/validate-https.sh
```

This script checks:
- ✓ Nginx and API services are running
- ✓ Ports 80 and 443 are listening
- ✓ SSL certificate exists and expiration date
- ✓ Nginx configuration is valid
- ✓ API responds on localhost:8080
- ✓ HTTPS endpoint is working
- ✓ HTTP to HTTPS redirect is active
- ✓ CORS and security headers are present

## Testing

### Test HTTPS Connection
```bash
# Basic health check
curl https://YOUR_SERVER_IP/health

# For self-signed certificates, use -k flag
curl -k https://YOUR_SERVER_IP/health

# Test API endpoint
curl -k https://YOUR_SERVER_IP/api/items
```

### Test HTTP to HTTPS Redirect
```bash
curl -I http://YOUR_SERVER_IP/health
# Should return: HTTP/1.1 301 Moved Permanently
# Location: https://YOUR_SERVER_IP/health
```

### Test from Browser
Open your browser and navigate to:
- `https://YOUR_SERVER_IP` (or your domain)
- For self-signed certs, you'll need to accept the security warning

### Verify Certificate
```bash
# Check certificate details
openssl s_client -servername YOUR_DOMAIN -connect YOUR_SERVER:443 -showcerts

# Check certificate chain
curl -vI https://YOUR_SERVER_IP 2>&1 | grep -A 10 "SSL connection"
```

## Security Group / Firewall

Ensure these ports are open:

| Port | Protocol | Purpose                          |
|------|----------|----------------------------------|
| 80   | HTTP     | Let's Encrypt & HTTP→HTTPS redirect |
| 443  | HTTPS    | Secure API access               |

**Port 8080 should NOT be exposed** - it's only accessed by nginx on localhost.

## Troubleshooting

### Nginx won't start
```bash
# Check nginx configuration
sudo nginx -t

# Check nginx error log
sudo tail -f /var/log/nginx/error.log

# Check nginx status
sudo systemctl status nginx
```

### Certificate errors
```bash
# Check certificate files exist
sudo ls -la /etc/letsencrypt/live/

# Check certificate permissions
sudo ls -la /etc/letsencrypt/live/example-api/

# Regenerate self-signed certificate
sudo /tmp/setup-letsencrypt.sh
```

### API not responding
```bash
# Check if Go API is running
sudo systemctl status example-api

# Check API logs
sudo journalctl -u example-api -f

# Test API directly (bypassing nginx)
curl http://localhost:8080/health
```

### Connection refused
```bash
# Check if nginx is running
sudo systemctl status nginx

# Check if ports are open
sudo netstat -tlnp | grep -E ':(80|443)'

# Check firewall/security group settings
```

## Manual Configuration Files

If you need to manually configure, here are the key files:

### Nginx Configuration
Location: `/etc/nginx/conf.d/example-api.conf`
```nginx
# See deployment/nginx-example-api.conf for full configuration
```

### Systemd Service
Location: `/etc/systemd/system/example-api.service`
```ini
# See deployment/example-api.service for configuration
```

### SSL Certificates
Location: `/etc/letsencrypt/live/example-api/` or `/etc/letsencrypt/live/your-domain.com/`
- `fullchain.pem`: Certificate chain
- `privkey.pem`: Private key

## Additional Resources

- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [Nginx SSL Configuration](https://nginx.org/en/docs/http/configuring_https_servers.html)
- [Certbot Documentation](https://certbot.eff.org/docs/)

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review nginx error logs: `/var/log/nginx/error.log`
3. Review API logs: `sudo journalctl -u example-api`
4. Open an issue on GitHub
