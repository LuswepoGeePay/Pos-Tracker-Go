package emailservices

// EmailTemplate defines the structure for email templates
type EmailTemplate struct {
	Subject string
	Body    string
}

// GetEmailTemplates returns all available email templates for reference
func GetEmailTemplates() map[string]EmailTemplate {
	return map[string]EmailTemplate{
		"welcome": {
			Subject: "Welcome to POS System",
			Body:    getWelcomeEmailTemplate("John Doe"),
		},
		"password_reset": {
			Subject: "Reset Your POS System Password",
			Body:    getPasswordResetEmailTemplate("https://example.com/reset-password"),
		},
		"verification": {
			Subject: "Verify Your POS System Account",
			Body:    getVerificationEmailTemplate("https://example.com/verify-email"),
		},
	}
}

// TestSendEmail is a utility function for testing email sending
// Use only in development environment
func TestSendEmail(toEmail string) error {
	testEmail := &EmailRequest{
		To:      toEmail,
		Subject: "POS System - Test Email",
		Body:    getTestEmailTemplate(),
	}
	return SendEmail(testEmail)
}

// getTestEmailTemplate returns a simple test email template
func getTestEmailTemplate() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 8px; }
        .header { background-color: #2196F3; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { padding: 20px; }
        .footer { background-color: #f5f5f5; padding: 20px; text-align: center; font-size: 12px; border-radius: 0 0 8px 8px; }
        .success { color: #4CAF50; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Email Configuration Test</h1>
        </div>
        <div class="content">
            <p><span class="success">✓ Success!</span></p>
            <p>Your email configuration is working correctly.</p>
            <p>This is a test email from the POS System.</p>
            <h3>Next Steps:</h3>
            <ul>
                <li>Verify email templates display correctly</li>
                <li>Check for any formatting issues</li>
                <li>Test other email types (welcome, password reset, etc.)</li>
            </ul>
        </div>
        <div class="footer">
            <p>&copy; 2026 POS System. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`
}
