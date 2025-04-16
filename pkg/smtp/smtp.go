package smtp

import (
	"fmt"
	"os"

	"github.com/go-gomail/gomail"
)

func NewClient() *gomail.Dialer {
	smtpHost := "mail.smtp2go.com"
	smtpPort := 2525
	smtpUser := os.Getenv("SMTPUSER")
	smtpPass := os.Getenv("SMTPPASS")

	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	return dialer
}

func sendConfirmationEmail(email, token string) error {
	// Generate the confirmation link
	confirmationLink := fmt.Sprintf("http://localhost/confirm?token=%s", token)

	// Plain text and HTML message bodies
	plainTextBody := fmt.Sprintf("Please confirm your email by clicking on the following link:\n%s", confirmationLink)
	htmlBody := fmt.Sprintf(`
        <p>Please confirm your email by clicking on the following link:</p>
        <p><a href="%s">Confirm your email</a></p>
        <p>If you cannot click the link, copy and paste this URL into your browser:</p>
        <p>%s</p>`, confirmationLink, confirmationLink)

	// SMTP2GO settings

	// Create the email message
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@artsbentley.com") // Sender email
	m.SetHeader("To", email)                        // Recipient email
	m.SetHeader("Subject", "Please Confirm Your Email Address")

	// Set both plain text and HTML bodies
	m.SetBody("text/plain", plainTextBody)  // Plain text version
	m.AddAlternative("text/html", htmlBody) // HTML version with clickable link

	// Set up the SMTP dialer

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}
	return nil
}
