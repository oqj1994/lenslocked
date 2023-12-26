package M

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type Email struct {
	To      string
	From    string
	Subject string
	Text    string
	HTML    string
}

//Msg 辅助生成gomail.Message 结构体

type SMTPConfig struct {
	Host     string
	Port     int
	UserName string
	Password string
}

func NewEmailService(cfg SMTPConfig) EmailService {
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.UserName, cfg.Password)
	return EmailService{
		dialer: d,
	}
}

type EmailService struct {
	DefaultSender string
	dialer        *gomail.Dialer
}

func (e EmailService) Send(m Email) error {
	msg := gomail.NewMessage()
	msg.SetHeader("To", m.To)
	e.setForm(msg, m)
	msg.SetHeader("Subject", m.Subject)
	switch {
	case m.Text != "" && m.HTML != "":
		msg.SetBody("text/plain", m.Text)
		msg.AddAlternative("text/html", m.HTML)

	case m.Text != "":
		msg.SetBody("text/plain", m.Text)
	case m.HTML != "":
		msg.SetBody("text/html", m.HTML)

	}
	err := e.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}

func (e EmailService) ForgetPassword(userName, resetURL string) error {
	email := Email{
		To:      userName,
		Subject: "To reset your password",
		Text:    "To reset your password, please follow this link: " + resetURL,
		HTML:    fmt.Sprintf(`To reset your password, please follow this link: <a href="%s">reset password</a>`, resetURL),
	}
	err := e.Send(email)
	if err != nil {
		return fmt.Errorf("forget password : %w", err)
	}
	return nil
}

func (e EmailService) setForm(m *gomail.Message, email Email) {
	if email.From != "" {
		m.SetHeader("From", email.From)
		return
	}
	if e.DefaultSender != "" {
		m.SetHeader("From", e.DefaultSender)
		return
	}
	m.SetHeader("From", DefaultSender)

}
