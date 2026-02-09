# Security Policy

## Reporting Security Vulnerabilities

If you discover a security vulnerability in this project, please report it by creating a private security advisory on GitHub or by emailing the maintainers directly. Please do **not** create a public issue for security vulnerabilities.

## Security Best Practices

This project follows security best practices to protect sensitive information:

### Secrets Management

1. **Never commit credentials to the repository**
   - API keys, passwords, tokens, and private keys should never be committed to version control
   - Use `.gitignore` to exclude sensitive files (`.env`, `*.pem`, `*.key`, etc.)

2. **Use GitHub Secrets for CI/CD**
   - All deployment credentials are stored in GitHub Secrets
   - Required secrets: `SSH_KEY`, `EC2_HOST`, `EC2_USER`
   - Configure secrets at: Settings → Secrets and variables → Actions

3. **Use Environment Variables**
   - Application configuration (like `PORT`) uses environment variables
   - Store sensitive configuration in `.env` files (which are gitignored)
   - Never hardcode credentials in source code

### Deployment Security

1. **SSH Key Security**
   - Private SSH keys are stored only in GitHub Secrets
   - Keys are written to disk temporarily during deployment and cleaned up automatically
   - Use strong key pairs (RSA 2048+ or Ed25519)

2. **Server Hardening**
   - Keep your EC2 instance and dependencies up to date
   - Use security groups to limit network access
   - Enable only necessary ports (8080 for API, 22 for SSH)
   - Use fail2ban or similar tools to prevent brute force attacks

3. **Application Security**
   - The systemd service runs with `NoNewPrivileges=true` and `PrivateTmp=true`
   - Regular security updates via `go mod tidy` and dependency audits
   - CORS is configured to allow cross-origin requests (adjust for production)

## What's Public (and That's OK)

The following information is public and acceptable in a public repository:

- **Server IP/Hostname**: Public-facing servers have public IPs. While we've moved these to secrets for better security posture, exposing a public server IP is not a vulnerability.
- **SSH Username**: Standard usernames like `ec2-user`, `ubuntu`, etc. are not secrets
- **Port Numbers**: Standard ports (8080, 80, 22) are not secrets
- **Architecture**: Server architecture (amd64, arm64) is not sensitive

## Dependency Security

- Dependencies are managed via `go.mod` and should be kept up to date
- Run `go mod tidy` and `go list -m -u all` regularly to check for updates
- Monitor GitHub security advisories for vulnerabilities

## Code Security

- Code is automatically scanned by GitHub's security features
- Follow Go security best practices
- Avoid SQL injection, XSS, and other common vulnerabilities
- Validate and sanitize all user inputs

## Supported Versions

We recommend always using the latest version from the `main` branch, as it contains the most recent security updates.
