package mail

import (
	"net/smtp"
)

type SmtpConfigs struct {
	UserSender     string
	PasswordSender string
	HostSender     string
	NameSender     string
	EmailSender    string
	Port           string
}

func NewSMTP(opts SmtpConfigs) MailSender {
	return &opts
}

func (this *SmtpConfigs) Send(sendTo []string, msg []byte) error {

	auth := smtp.PlainAuth(this.NameSender, this.UserSender, this.PasswordSender, this.HostSender)
	err := smtp.SendMail(this.HostSender+":"+this.Port, auth, this.EmailSender, sendTo, msg)

	return err
}
