# Environment Configuration Guide

This document explains how to configure the POS System backend using environment variables instead of hardcoded values.

## Overview of Changes

The backend has been refactored to use environment variables for all configuration values. This includes:

1. **Database Configuration** - Connection string now loaded from env
2. **JWT Configuration** - Secret and expiration time from env
3. **Server Configuration** - Port and environment type from env
4. **Email Configuration** - SMTP settings from env
5. **SSL/TLS Configuration** - Certificate paths from env
6. **Logging Configuration** - Log file path and level from env

## Environment Variables

### Server Configuration
- `SERVER_PORT` - Port to run the server on (default: 8050)
- `ENVIRONMENT` - Environment type: production, staging, development (default: development)

### Database Configuration
- `DATABASE_URL` - Full database connection string
  - Example: `tracker_user:tracker_user@tcp(10.139.40.25:3306)/posmaster?charset=utf8mb4&parseTime=True&loc=Local`

### JWT Configuration
- `JWT_SECRET` - Secret key for signing JWT tokens (REQUIRED - must be set in production)
- `JWT_EXPIRATION_HOURS` - Token expiration time in hours (default: 24)

### Email Configuration (SMTP)
- `SMTP_HOST` - SMTP server hostname (e.g., smtp.gmail.com)
- `SMTP_PORT` - SMTP server port (default: 587)
- `SMTP_USERNAME` - Email account username
- `SMTP_PASSWORD` - Email account password or app password
- `SMTP_FROM_ADDRESS` - Sender email address
- `SMTP_FROM_NAME` - Display name for sender (default: POS System)

### SSL/TLS Configuration
- `CERT_FILE` - Path to SSL certificate file (optional)
- `KEY_FILE` - Path to SSL private key file (optional)

### Logging Configuration
- `LOG_FILE` - Path to log file (default: pos_master.log)
- `LOG_LEVEL` - Log level: debug, info, warn, error (default: info)

## Setting Up Environment Variables

### Method 1: Using .env File
Create a `.env` file in the root directory with your configuration:

```bash
# Server Configuration
export SERVER_PORT=8050
export ENVIRONMENT=production

# Database Configuration
export DATABASE_URL="user:password@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

# JWT Configuration
export JWT_SECRET="your-very-secure-and-long-secret-key-here"
export JWT_EXPIRATION_HOURS=24

# Email Configuration
export SMTP_HOST="smtp.gmail.com"
export SMTP_PORT=587
export SMTP_USERNAME="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"
export SMTP_FROM_ADDRESS="noreply@yourdomain.com"
export SMTP_FROM_NAME="POS System"

# SSL/TLS Configuration
export CERT_FILE="/path/to/cert.pem"
export KEY_FILE="/path/to/key.pem"

# Logging Configuration
export LOG_FILE="pos_master.log"
export LOG_LEVEL="info"
```

Then load it before running:
```bash
source .env
go run main.go
```

### Method 2: System Environment Variables
Set environment variables directly in your system:

```bash
export SERVER_PORT=8050
export DATABASE_URL="user:password@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
export JWT_SECRET="your-secure-secret-key"
export SMTP_HOST="smtp.gmail.com"
export SMTP_PORT=587
export SMTP_USERNAME="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"
export SMTP_FROM_ADDRESS="noreply@yourdomain.com"
go run main.go
```

### Method 3: Docker or Docker Compose
Pass environment variables to your container:

```dockerfile
ENV SERVER_PORT=8050
ENV DATABASE_URL="user:password@tcp(db:3306)/posmaster?charset=utf8mb4&parseTime=True&loc=Local"
ENV JWT_SECRET="your-secure-secret-key"
ENV SMTP_HOST="smtp.gmail.com"
ENV SMTP_PORT=587
ENV SMTP_USERNAME="your-email@gmail.com"
ENV SMTP_PASSWORD="your-app-password"
ENV SMTP_FROM_ADDRESS="noreply@yourdomain.com"
```

Or in docker-compose.yml:
```yaml
environment:
  - SERVER_PORT=8050
  - DATABASE_URL=user:password@tcp(db:3306)/posmaster?charset=utf8mb4&parseTime=True&loc=Local
  - JWT_SECRET=your-secure-secret-key
  - SMTP_HOST=smtp.gmail.com
  - SMTP_PORT=587
  - SMTP_USERNAME=your-email@gmail.com
  - SMTP_PASSWORD=your-app-password
  - SMTP_FROM_ADDRESS=noreply@yourdomain.com
```

### Method 4: Cloud Platforms (AWS, Azure, GCP)
- **AWS:** Use AWS Secrets Manager or Parameter Store
- **Azure:** Use Azure Key Vault
- **GCP:** Use Secret Manager

## Email Service

The email service has been completely refactored to use HTML templates directly in code instead of reading from files.

### Available Email Functions

#### SendWelcomeEmail
Sent automatically when a new user registers.

```go
emailservices.SendWelcomeEmail(email, userName)
```

#### SendPasswordResetEmail
Send password reset email with reset link.

```go
emailservices.SendPasswordResetEmail(email, resetURL)
```

#### SendAccountVerificationEmail
Send account verification email.

```go
emailservices.SendAccountVerificationEmail(email, verificationURL)
```

#### SendCustomEmail
Send custom email with your own HTML body.

```go
emailservices.SendCustomEmail(email, subject, htmlBody)
```

### Email Templates

All email templates are now defined as HTML strings in the `services/emailservices/email_service.go` file. You can customize them by editing the template functions:

- `getWelcomeEmailTemplate(userName)` - Welcome email for new users
- `getPasswordResetEmailTemplate(resetURL)` - Password reset email
- `getVerificationEmailTemplate(verificationURL)` - Account verification email

## Configuration Initialization

The configuration is automatically loaded when the application starts:

```go
func main() {
    // Initialize configuration from environment variables
    config.InitConfig()
    
    // Rest of initialization...
    config.InitDB()
}
```

## Required Environment Variables

The following environment variables are **REQUIRED** and must be set:

1. `DATABASE_URL` - Database connection string
2. `JWT_SECRET` - JWT signing secret (must be secure and long in production)
3. `SMTP_HOST` - SMTP server hostname
4. `SMTP_USERNAME` - SMTP username
5. `SMTP_PASSWORD` - SMTP password
6. `SMTP_FROM_ADDRESS` - From email address

If any required variable is missing, the application will panic with an error message.

## Optional Environment Variables

These variables have sensible defaults:

- `SERVER_PORT` (default: 8050)
- `ENVIRONMENT` (default: development)
- `JWT_EXPIRATION_HOURS` (default: 24)
- `SMTP_PORT` (default: 587)
- `SMTP_FROM_NAME` (default: POS System)
- `LOG_FILE` (default: pos_master.log)
- `LOG_LEVEL` (default: info)
- `CERT_FILE` (optional)
- `KEY_FILE` (optional)

## Best Practices

1. **Never commit secrets to version control** - Use .gitignore to exclude .env files
2. **Use strong JWT secrets** - Use a random string of at least 32 characters
3. **Use app passwords** - For Gmail and other services, use app-specific passwords instead of your main password
4. **Rotate secrets regularly** - Change JWT secret and email passwords periodically in production
5. **Environment-specific configuration** - Use different values for development, staging, and production
6. **Log sensitive data carefully** - Never log passwords or tokens
7. **Use HTTPS in production** - Always enable SSL/TLS in production

## Troubleshooting

### Missing environment variables
If the application crashes with a message about missing environment variables, ensure all required variables are set.

### Email not sending
- Check SMTP credentials are correct
- Verify SMTP_HOST and SMTP_PORT are correct
- Check if your email provider requires app-specific passwords
- Verify the from_address is allowed by your email provider
- Check application logs for specific error messages

### JWT token errors
- Ensure JWT_SECRET is set correctly
- Verify JWT_SECRET is the same across all server instances
- Check token expiration time with JWT_EXPIRATION_HOURS

### Database connection issues
- Verify DATABASE_URL format is correct for your database
- Check database host, port, username, and password
- Ensure database user has necessary permissions

## Migration from Old System

If you were previously using hardcoded configuration in files, follow these steps:

1. Copy values from your old configuration to environment variables
2. Update the .env file with your values
3. Load environment variables before starting the application
4. Test thoroughly in a staging environment first
5. Deploy to production with appropriate environment variables

## Files Modified

- `config/env.go` - New configuration loading system
- `config/database.go` - Updated to use DATABASE_URL env variable
- `main.go` - Updated to initialize config and use SERVER_PORT
- `services/authservices/auth_service.go` - Updated to use JWT_SECRET from config
- `services/emailservices/email_service.go` - New email service with HTML templates
- `services/user_services/register_user_service.go` - Updated to send welcome emails
- `.env` - Updated with all environment variables
