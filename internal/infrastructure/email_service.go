package infrastructure

import (
	"fmt"
	"net/smtp"
	// "os"
)

// MailtrapService implements IEmailService using Mailtrap's SMTP
type MailtrapService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewMailtrapService(host, port, username, password, from string) *MailtrapService {
	return &MailtrapService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// SendVerificationEmail sends a registration verification code/link
func (m *MailtrapService) SendVerificationEmail(email, code, baseUrl string) error {
	subject := "Verify Your Account"
	
	body := fmt.Sprintf("Click this link to verify your account:\n\n%s/auth/verify?code=%s",baseUrl, code)
	return m.sendEmail(email, subject, body)
}

// SendPasswordReset sends a reset password link/token
func (m *MailtrapService) SendPasswordReset(email, token, baseUrl string) error {
	subject := "Reset Your Password"
	body := fmt.Sprintf("Click this link to reset your password:\n\n%s/auth/reset-password?token=%s",baseUrl, token)
	return m.sendEmail(email, subject, body)
}
func (m *MailtrapService) SendUpdateVerificationEmail(email, code, baseUrl string) error {
	subject := "Verify Your New Eamil Account"
	
	body := fmt.Sprintf("If you did not request update Email ignore these message\n\nClick this link to verify your account:\n\n%s/auth/verify-Update-Email?code=%s",baseUrl, code)
	return m.sendEmail(email, subject, body)

}

// internal shared logic
func (m *MailtrapService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		to, subject, body,
	))

	address := fmt.Sprintf("%s:%s", m.host, m.port)

	return smtp.SendMail(address, auth, m.from, []string{to}, msg)
}


type EmailService struct {
    mailtrap *MailtrapService
	baseUrl string
}

func NewEmailService(mailtrap *MailtrapService, url string) *EmailService {
    return &EmailService{mailtrap: mailtrap}
}

func (e *EmailService) SendVerificationEmail(email, code string) error {
    return e.mailtrap.SendVerificationEmail(email, code, e.baseUrl)
}

func (e *EmailService) SendPasswordReset(email, token string) error {
    return e.mailtrap.SendPasswordReset(email, token, e.baseUrl)
}

func (e *EmailService) SendUpdateVerificationEmail(email, code string) error {
    return e.mailtrap.SendUpdateVerificationEmail(email, code, e.baseUrl)
}