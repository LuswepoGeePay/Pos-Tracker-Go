package emailservices

import (
	"fmt"
	"net/smtp"
	"os"
)

// EmailRequest holds email sending request data
type EmailRequest struct {
	To      string
	Subject string
	Body    string
}

// SendEmail sends an email with HTML content
func SendEmail(req *EmailRequest) error {
	// Setup SMTP authentication
	auth := smtp.PlainAuth("", os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST"))

	// Compose email headers
	from := fmt.Sprintf("%s <%s>", os.Getenv("SMTP_FROM_NAME"), os.Getenv("SMTP_FROM_ADDRESS"))
	headers := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nContent-Type: text/html; charset=\"UTF-8\"\n", from, req.To, req.Subject)

	// Complete message with headers
	message := headers + "\n" + req.Body

	// Send email
	addr := fmt.Sprintf("%s:%s", os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"))
	err := smtp.SendMail(addr, auth, os.Getenv("SMTP_FROM_ADDRESS"), []string{req.To}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendWelcomeEmail sends a welcome email to new user
func SendWelcomeEmail(toEmail, userName string) error {
	htmlBody := getWelcomeEmailTemplate(userName)
	req := &EmailRequest{
		To:      toEmail,
		Subject: "Welcome to POS Tracker System",
		Body:    htmlBody,
	}
	return SendEmail(req)
}

// SendPasswordResetEmail sends a password reset email
func SendPasswordResetEmail(toEmail, resetURL string) error {
	htmlBody := getPasswordResetEmailTemplate(resetURL)
	req := &EmailRequest{
		To:      toEmail,
		Subject: "Reset Your POS Tracker System Password",
		Body:    htmlBody,
	}
	return SendEmail(req)
}

// SendAccountVerificationEmail sends account verification email
func SendAccountVerificationEmail(toEmail, verificationURL string) error {
	htmlBody := getVerificationEmailTemplate(verificationURL)
	req := &EmailRequest{
		To:      toEmail,
		Subject: "Verify Your POS Tracker System Account",
		Body:    htmlBody,
	}
	return SendEmail(req)
}

// getWelcomeEmailTemplate returns HTML template for welcome email
func getWelcomeEmailTemplate(userName string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 8px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { padding: 20px; }
        .footer { background-color: #f5f5f5; padding: 20px; text-align: center; font-size: 12px; border-radius: 0 0 8px 8px; }
        .button { background-color: #4CAF50; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to POS Tracker System</h1>
        </div>
        <div class="content">
            <p>Hello <strong>%s</strong>,</p>
            <p>Thank you for registering with our POS Tracker System. We're excited to have you on board!</p>
            <p>Your account is now active and ready to use. You can log in with your credentials at any time.</p>
            <h3>What's Next?</h3>
            <ul>
                <li>Complete your business profile</li>
                <li>Set up your devices</li>
                <li>Configure your applications</li>
                <li>Start tracking your sales</li>
            </ul>
            <p>If you have any questions, our support team is here to help.</p>
        </div>
        <div class="footer">
            <p>&copy; 2026 POS Tracker System. All rights reserved.</p>
            <p>This is an automated message. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`, userName)
}

// getPasswordResetEmailTemplate returns HTML template for password reset email
func getPasswordResetEmailTemplate(resetURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 8px; }
        .header { background-color: #FF9800; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { padding: 20px; }
        .footer { background-color: #f5f5f5; padding: 20px; text-align: center; font-size: 12px; border-radius: 0 0 8px 8px; }
        .button { background-color: #FF9800; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block; margin: 20px 0; }
        .warning { background-color: #fff3cd; border: 1px solid #ffc107; padding: 10px; border-radius: 4px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <p>We received a request to reset your password for your POS Tracker System account.</p>
            <p>Click the button below to create a new password. This link will expire in 24 hours.</p>
            <a href="%s" class="button">Reset My Password</a>
            <div class="warning">
                <strong>Didn't request a password reset?</strong>
                <p>If you didn't request this, you can safely ignore this email. Your password will not change.</p>
            </div>
            <p>If the button above doesn't work, copy and paste this link into your browser:</p>
            <p style="word-break: break-all;"><small>%s</small></p>
        </div>
        <div class="footer">
            <p>&copy; 2026 POS Tracker System. All rights reserved.</p>
            <p>This is an automated message. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`, resetURL, resetURL)
}

// getVerificationEmailTemplate returns HTML template for account verification email
func getVerificationEmailTemplate(verificationURL string) string {
	return fmt.Sprintf(`
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
        .button { background-color: #2196F3; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Verify Your Email Address</h1>
        </div>
        <div class="content">
            <p>Thank you for signing up! Please verify your email address to activate your account.</p>
            <p>Click the button below to verify your email:</p>
            <a href="%s" class="button">Verify Email Address</a>
            <p>If the button above doesn't work, copy and paste this link into your browser:</p>
            <p style="word-break: break-all;"><small>%s</small></p>
            <p>This verification link will expire in 48 hours.</p>
        </div>
        <div class="footer">
            <p>&copy; 2026 POS Tracker System. All rights reserved.</p>
            <p>This is an automated message. Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`, verificationURL, verificationURL)
}

// SendCustomEmail sends a custom email with provided HTML body
func SendCustomEmail(toEmail, subject, htmlBody string) error {
	req := &EmailRequest{
		To:      toEmail,
		Subject: subject,
		Body:    htmlBody,
	}
	return SendEmail(req)
}
