#!/bin/bash

# Setup script for Let's Encrypt SSL certificate with IP-based access
# This script sets up certbot and obtains an SSL certificate for the server

set -e

echo "=== Let's Encrypt SSL Certificate Setup ==="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "ERROR: This script must be run as root (use sudo)"
    exit 1
fi

# Get server IP address
SERVER_IP=$(curl -s http://checkip.amazonaws.com || curl -s http://ifconfig.me || curl -s http://icanhazip.com)
echo "Detected server IP: $SERVER_IP"
echo ""

# Install certbot and nginx if not already installed
echo "Installing required packages..."
if command -v yum &> /dev/null; then
    # Detect OS version
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        # Amazon Linux 2023 and later don't need EPEL repository
        # VERSION_ID can be "2023" or "2023.x.x"
        if [ "$ID" = "amzn" ] && [[ "$VERSION_ID" == 2023* ]]; then
            echo "Detected Amazon Linux 2023 - installing packages directly..."
            yum install -y certbot nginx cronie
        else
            # Amazon Linux 2 / RHEL / CentOS - install EPEL first
            echo "Installing EPEL repository..."
            yum install -y epel-release
            yum install -y certbot nginx cronie
        fi
    else
        # Fallback for older systems
        yum install -y epel-release
        yum install -y certbot nginx cronie
    fi
elif command -v apt-get &> /dev/null; then
    # Ubuntu / Debian
    apt-get update
    apt-get install -y certbot nginx cron
else
    echo "ERROR: Unsupported package manager. Please install certbot and nginx manually."
    exit 1
fi

# Enable and start cron service
echo "Enabling cron service..."
if command -v systemctl &> /dev/null; then
    systemctl enable crond 2>/dev/null || systemctl enable cron 2>/dev/null || true
    systemctl start crond 2>/dev/null || systemctl start cron 2>/dev/null || true
fi

# Create directory for Let's Encrypt challenges
mkdir -p /var/www/certbot

# Stop nginx if running (we need port 80 for initial certificate request)
systemctl stop nginx || true

# Request certificate using standalone mode
# Note: IP-based certificates require DNS validation or can use certbot's standalone mode
# For IP-based access, we'll use standalone mode which doesn't require a domain
echo ""
echo "Requesting SSL certificate..."
echo "NOTE: For IP-based certificates, Let's Encrypt has limitations."
echo "This setup uses a self-signed certificate approach for IP-based access."
echo ""

# For production IP-based access, we'll generate a self-signed certificate
# Let's Encrypt requires domain names, so for IP-only access we use self-signed certs
CERT_DIR="/etc/letsencrypt/live/example-api"
mkdir -p "$CERT_DIR"

# Generate self-signed certificate for IP-based access
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout "$CERT_DIR/privkey.pem" \
    -out "$CERT_DIR/fullchain.pem" \
    -subj "/C=US/ST=State/L=City/O=Organization/CN=$SERVER_IP" \
    -addext "subjectAltName=IP:$SERVER_IP"

chmod 600 "$CERT_DIR/privkey.pem"
chmod 644 "$CERT_DIR/fullchain.pem"

echo ""
echo "Certificate generated successfully!"
echo ""
echo "IMPORTANT: This is a self-signed certificate."
echo "For production use with a domain name, run:"
echo "  sudo certbot certonly --standalone -d your-domain.com"
echo ""

# Set up certificate renewal (if using certbot)
# For self-signed certs, this is not needed, but we'll include it for future use
echo "Setting up automatic certificate renewal..."
(crontab -l 2>/dev/null; echo "0 0,12 * * * certbot renew --quiet --post-hook 'systemctl reload nginx'") | crontab -

# Deploy nginx configuration
echo ""
echo "Deploying nginx configuration..."
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cp "$SCRIPT_DIR/nginx-example-api.conf" /etc/nginx/conf.d/example-api.conf

# Test nginx configuration
echo "Testing nginx configuration..."
if ! nginx -t; then
    echo "ERROR: Nginx configuration test failed!"
    echo "Please check /etc/nginx/conf.d/example-api.conf for errors."
    exit 1
fi

# Enable and start nginx
echo "Enabling and starting nginx..."
systemctl enable nginx
systemctl start nginx

# Wait for nginx to start
sleep 2

# Verify HTTPS is working
echo ""
echo "Verifying HTTPS is working..."
if curl -k -s -f https://$SERVER_IP/health > /dev/null; then
    echo "✓ HTTPS is working correctly!"
    echo ""
    echo "You can test the API with:"
    echo "  curl -k https://$SERVER_IP/health"
else
    echo "WARNING: HTTPS verification failed. Please check nginx logs:"
    echo "  sudo journalctl -u nginx -n 50"
fi

echo ""
echo "=== Setup Complete ==="
echo ""
echo "Your API is now accessible via HTTPS at: https://$SERVER_IP"
echo ""
echo "Note: Browsers will show a security warning for self-signed certificates."
echo "For production with a domain name, obtain a real Let's Encrypt certificate by running:"
echo "  sudo certbot certonly --standalone -d your-domain.com"
echo "  sudo systemctl restart nginx"
