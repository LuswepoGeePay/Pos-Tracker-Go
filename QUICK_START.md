# Quick Start Guide - Environment Variables Setup

## ⚡ Quick Setup (5 minutes)

### 1. Copy the template
```bash
cp .env.example .env
```

### 2. Edit .env with your values
```bash
nano .env  # or use your favorite editor
```

### 3. Required values to change:
```
DATABASE_URL="your-database-url"
JWT_SECRET="your-256-bit-secret"
SMTP_HOST="smtp.gmail.com"
SMTP_USERNAME="your-email@gmail.com"
SMTP_PASSWORD="your-app-password"
SMTP_FROM_ADDRESS="noreply@yourdomain.com"
```

### 4. Load and run
```bash
source .env
go run main.go
```

---

## 📋 Common SMTP Providers

### Gmail
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-16-char-app-password
# Get app password: https://support.google.com/accounts/answer/185833
```

### Outlook/Office 365
```
SMTP_HOST=smtp.office365.com
SMTP_PORT=587
SMTP_USERNAME=your-email@outlook.com
SMTP_PASSWORD=your-email-password
```

### SendGrid
```
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=SG.xxxxxx-your-api-key
```

### AWS SES
```
SMTP_HOST=email-smtp.us-east-1.amazonaws.com
SMTP_PORT=587
SMTP_USERNAME=your-ses-username
SMTP_PASSWORD=your-ses-password
```

---

## 🔑 Generate Secure JWT Secret

**Option 1: Using OpenSSL**
```bash
openssl rand -base64 32
```

**Option 2: Using Python**
```bash
python3 -c "import secrets; print(secrets.token_urlsafe(32))"
```

**Option 3: Using Go**
```bash
go run -c "import crypto/rand, encoding/base64; b := make([]byte, 32); rand.Read(b); println(base64.StdEncoding.EncodeToString(b))"
```

---

## 🐳 Docker/Docker Compose Setup

### docker-compose.yml
```yaml
version: '3.8'
services:
  pos-backend:
    build: .
    ports:
      - "8050:8050"
    environment:
      - SERVER_PORT=8050
      - ENVIRONMENT=production
      - DATABASE_URL=user:pass@tcp(db:3306)/posmaster?charset=utf8mb4&parseTime=True&loc=Local
      - JWT_SECRET=your-secure-secret
      - SMTP_HOST=smtp.gmail.com
      - SMTP_PORT=587
      - SMTP_USERNAME=your-email@gmail.com
      - SMTP_PASSWORD=your-app-password
      - SMTP_FROM_ADDRESS=noreply@yourdomain.com
    depends_on:
      - db

  db:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=rootpass
      - MYSQL_DATABASE=posmaster
```

---

## 🚀 Deployment Checklist

### Before going live:
- [ ] Change `ENVIRONMENT=production`
- [ ] Use strong JWT_SECRET (32+ characters)
- [ ] Use app-specific password for email (not your main password)
- [ ] Verify DATABASE_URL points to production database
- [ ] Test email sending works
- [ ] Check logs for any configuration errors
- [ ] Set up SSL/TLS (CERT_FILE and KEY_FILE)
- [ ] Review and secure all environment variables

### Production Security Tips:
- [ ] Never commit .env to git (use .gitignore)
- [ ] Use secret management (AWS Secrets Manager, Azure Key Vault, etc.)
- [ ] Rotate secrets regularly
- [ ] Use strong, unique passwords
- [ ] Enable HTTPS/SSL in production
- [ ] Monitor logs for security issues

---

## 🧪 Testing Email Configuration

### Test endpoint (add to your routes if needed):
```go
r.POST("/test-email", func(c *gin.Context) {
    email := c.Query("email")
    if email == "" {
        c.JSON(400, gin.H{"error": "email query parameter required"})
        return
    }
    
    err := emailservices.TestSendEmail(email)
    if err != nil {
        c.JSON(500, gin.H{"error": fmt.Sprintf("Email failed: %v", err)})
        return
    }
    
    c.JSON(200, gin.H{"message": "Test email sent successfully"})
})

// Then test with:
// curl "http://localhost:8050/test-email?email=yourtest@example.com"
```

---

## 🔍 Troubleshooting

### Issue: "Required environment variable not set"
**Solution:** Check that all required variables are exported in your .env file and properly loaded.

### Issue: Email not sending
**Solution:** 
- Verify SMTP credentials are correct
- For Gmail, use app-specific password (not main password)
- Check firewall allows outgoing port 587
- Review application logs for specific error

### Issue: Database connection failed
**Solution:**
- Verify DATABASE_URL format is correct
- Check database host, port, username, password
- Ensure database user has necessary privileges
- Verify database server is running and accessible

### Issue: JWT token invalid
**Solution:**
- Ensure JWT_SECRET is set correctly
- Check JWT_SECRET is same across all server instances
- Verify token hasn't expired

---

## 📚 More Information

See the following files for detailed documentation:
- [`ENVIRONMENT_CONFIG.md`](ENVIRONMENT_CONFIG.md) - Complete configuration guide
- [`CHANGES_SUMMARY.md`](CHANGES_SUMMARY.md) - Detailed list of all changes
- [`.env.example`](.env.example) - Environment variables template

---

## 🆘 Need Help?

1. Check the logs for error messages: `tail -f pos_master.log`
2. Review `ENVIRONMENT_CONFIG.md` for configuration details
3. Test email with the test endpoint
4. Ensure all required environment variables are set
5. Check that services (MySQL, SMTP) are accessible from your server

---

## ✅ Verification Checklist

After setup, verify:
- [ ] Server starts without errors
- [ ] Can connect to database
- [ ] JWT token generation works
- [ ] User registration succeeds
- [ ] Welcome email is sent
- [ ] All logs show "info" level messages only

---

**Last Updated:** March 2026
**Version:** 1.0
