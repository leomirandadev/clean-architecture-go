package mail

import (
	"fmt"
	"net/smtp"
)

type Options struct {
	User        string `mapstructure:"user" validate:"required"`
	Password    string `mapstructure:"password" validate:"required"`
	Host        string `mapstructure:"host" validate:"required"`
	NameSender  string `mapstructure:"name_sender" validate:"required"`
	EmailSender string `mapstructure:"email_sender" validate:"required"`
	Port        string `mapstructure:"port" validate:"required"`
}

func NewSMTP(opts Options) MailSender {
	return &implSmtp{opts}
}

type implSmtp struct {
	Options
}

func (s implSmtp) Send(sendTo []string, subject, body string) error {

	auth := smtp.PlainAuth(s.NameSender, s.User, s.Password, s.Host)
	msg := s.buildMessage(sendTo, subject, body)

	err := smtp.SendMail(s.Host+":"+s.Port, auth, s.EmailSender, sendTo, []byte(msg))

	return err
}

func (s implSmtp) buildMessage(sendTo []string, subject, body string) string {
	msg := "MIME-version: 1.0\nContent-Type: text/html; charset=iso-8859-1\r\n"
	msg += fmt.Sprintf("From: %s <%s>\r\n", s.NameSender, s.EmailSender)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += fmt.Sprintf("\r\n%s\r\n", body)

	return msg
}
