package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"letterboxd-cineville/db"
	"letterboxd-cineville/model"
	"log/slog"
	"os"

	"github.com/go-gomail/gomail"
)

func generateToken() (string, error) {
	tokenBytes := make([]byte, 16)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenBytes), nil
}

// TODO:
// func sendConfirmationEmail(email, token string) error {
// 	confirmationLink := fmt.Sprintf("https://yourdomain.com/confirm?token=%s", token)
// 	message := fmt.Sprintf("Please confirm your email by clicking on the following link: %s", confirmationLink)
// 	fmt.Println(message)
//
// 	// Set up your email settings and send the email
// 	// (using gomail or another library)
// 	return nil
// }

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
	smtpHost := "mail.smtp2go.com"
	smtpPort := 2525
	smtpUser := os.Getenv("SMTPUSER")
	smtpPass := os.Getenv("SMTPPASS")

	// Create the email message
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@artsbentley.com") // Sender email
	m.SetHeader("To", email)                        // Recipient email
	m.SetHeader("Subject", "Please Confirm Your Email Address")

	// Set both plain text and HTML bodies
	m.SetBody("text/plain", plainTextBody)  // Plain text version
	m.AddAlternative("text/html", htmlBody) // HTML version with clickable link

	// Set up the SMTP dialer
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}
	return nil
}

type Service struct {
	store  *db.Store
	logger *slog.Logger
}

func NewService(store *db.Store) *Service {
	return &Service{
		store:  store,
		logger: slog.Default(),
	}
}

func (s *Service) GetAllUsers() ([]model.User, error) {
	users, err := s.store.GetAllUsers()
	return users, err
}

func (s *Service) ConfirmUserEmail(token string) error {
	id, err := s.store.GetUserIDByToken(token)
	if err != nil {
		return err
	}

	return s.store.ConfirmUserEmail(id)
}

func (s *Service) CreateNewUser(email string, username string) error {
	token, err := generateToken()
	if err != nil {
		return err
	}

	user := model.User{
		Email:              email,
		LetterboxdUsername: username,
		Token:              token,
		Watchlist:          make([]string, 0),
	}

	err = s.store.CreateNewUser(user)
	// TODO: handle properly
	if err != nil {
		return err
	}

	err = sendConfirmationEmail(user.Email, token)
	if err != nil {
		return err
	}

	return nil
}
