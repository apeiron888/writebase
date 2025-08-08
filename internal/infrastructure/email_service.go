package infrastructure

import (
	"fmt"
	"log"
	// "net/smtp"
	"strconv"

	gomail "gopkg.in/mail.v2"
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
	
	body := fmt.Sprintf("If you did not request update Email ignore these message\n\nClick this link to verify your account:\n\n%s/auth/verify-Update-Email?token=%s",baseUrl, code)
	return m.sendEmail(email, subject, body)

}



func (m *MailtrapService) sendEmail(to, subject, body string) error{
	message := gomail.NewMessage()

	message.SetHeader("From", m.from)
	message.SetHeader("To", to)
	message.SetHeader("subject", subject)
	message.SetBody("text/plain", body)
	portint,err:= strconv.Atoi(m.port)
	log.Printf("Error:%v", err)
	dialer:= gomail.NewDialer(m.host, portint, m.username, m.password)
	dialer.SSL= true
	if err := dialer.DialAndSend(message); err != nil {
        log.Printf("Error:%v", err)
		return err
        // panic(err)
    } else {
        fmt.Println("Email sent successfully!")
    }
	return nil
}


type EmailService struct {
    mailtrap *MailtrapService
	baseUrl string
}

func NewEmailService(mailtrap *MailtrapService, url string) *EmailService {
    return &EmailService{mailtrap: mailtrap, baseUrl: url}
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