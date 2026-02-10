#!/bin/bash

# Validation script to check HTTPS configuration
# Run this after setting up HTTPS to verify everything is working

set -e

echo "=== HTTPS Configuration Validation ==="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get server IP
SERVER_IP=$(curl -s http://checkip.amazonaws.com || curl -s http://ifconfig.me || curl -s http://icanhazip.com)
echo "Server IP: $SERVER_IP"
echo ""

# Check 1: Nginx running
echo -n "Checking nginx status... "
if systemctl is-active --quiet nginx; then
    echo -e "${GREEN}✓ Running${NC}"
else
    echo -e "${RED}✗ Not running${NC}"
    echo "  Run: sudo systemctl start nginx"
fi

# Check 2: API running
echo -n "Checking example-api status... "
if systemctl is-active --quiet example-api; then
    echo -e "${GREEN}✓ Running${NC}"
else
    echo -e "${RED}✗ Not running${NC}"
    echo "  Run: sudo systemctl start example-api"
fi

# Check 3: Port 80 listening
echo -n "Checking port 80 (HTTP)... "
if sudo netstat -tlnp | grep -q ':80 '; then
    echo -e "${GREEN}✓ Listening${NC}"
else
    echo -e "${RED}✗ Not listening${NC}"
fi

# Check 4: Port 443 listening
echo -n "Checking port 443 (HTTPS)... "
if sudo netstat -tlnp | grep -q ':443 '; then
    echo -e "${GREEN}✓ Listening${NC}"
else
    echo -e "${RED}✗ Not listening${NC}"
fi

# Check 5: Port 8080 (API) listening on localhost only
echo -n "Checking port 8080 (API on localhost)... "
if netstat -tln | grep -q '127.0.0.1:8080'; then
    echo -e "${GREEN}✓ Listening on localhost${NC}"
elif netstat -tln | grep -q '0.0.0.0:8080'; then
    echo -e "${YELLOW}⚠ Listening on all interfaces (should be localhost only)${NC}"
else
    echo -e "${RED}✗ Not listening${NC}"
fi

# Check 6: SSL certificate exists
echo -n "Checking SSL certificate... "
if [ -f /etc/letsencrypt/live/example-api/fullchain.pem ]; then
    echo -e "${GREEN}✓ Found${NC}"
    # Check expiration
    EXPIRY=$(openssl x509 -enddate -noout -in /etc/letsencrypt/live/example-api/fullchain.pem | cut -d= -f2)
    echo "  Expires: $EXPIRY"
else
    echo -e "${RED}✗ Not found${NC}"
    echo "  Run: sudo /tmp/setup-letsencrypt.sh"
fi

# Check 7: Nginx configuration
echo -n "Checking nginx configuration syntax... "
if sudo nginx -t > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Valid${NC}"
else
    echo -e "${RED}✗ Invalid${NC}"
    echo "  Run: sudo nginx -t"
fi

echo ""
echo "=== Functional Tests ==="
echo ""

# Test 8: API health check (direct)
echo -n "Testing API directly (localhost:8080)... "
if curl -s http://localhost:8080/health | grep -q "healthy"; then
    echo -e "${GREEN}✓ Working${NC}"
else
    echo -e "${RED}✗ Failed${NC}"
fi

# Test 9: HTTPS health check
echo -n "Testing HTTPS endpoint... "
if curl -k -s https://localhost/health | grep -q "healthy"; then
    echo -e "${GREEN}✓ Working${NC}"
else
    echo -e "${RED}✗ Failed${NC}"
fi

# Test 10: HTTP to HTTPS redirect
echo -n "Testing HTTP to HTTPS redirect... "
REDIRECT=$(curl -s -o /dev/null -w "%{http_code}" http://localhost/health)
if [ "$REDIRECT" = "301" ]; then
    echo -e "${GREEN}✓ Redirecting (301)${NC}"
else
    echo -e "${YELLOW}⚠ Status: $REDIRECT (expected 301)${NC}"
fi

# Test 11: CORS headers
echo -n "Testing CORS headers... "
CORS_HEADER=$(curl -s -I https://localhost/health -k | grep -i "access-control-allow-origin")
if [ -n "$CORS_HEADER" ]; then
    echo -e "${GREEN}✓ Present${NC}"
else
    echo -e "${YELLOW}⚠ Not found${NC}"
fi

# Test 12: Security headers
echo -n "Testing security headers (HSTS)... "
HSTS_HEADER=$(curl -s -I https://localhost/health -k | grep -i "strict-transport-security")
if [ -n "$HSTS_HEADER" ]; then
    echo -e "${GREEN}✓ Present${NC}"
else
    echo -e "${YELLOW}⚠ Not found${NC}"
fi

echo ""
echo "=== Validation Complete ==="
echo ""
echo "To test from external machine:"
echo "  curl -k https://$SERVER_IP/health"
echo "  curl -k https://$SERVER_IP/api/items"
echo ""
echo "Note: Use -k flag for self-signed certificates"
echo "For production, obtain a real certificate with: sudo certbot certonly --standalone -d your-domain.com"
