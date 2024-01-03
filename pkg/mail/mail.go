package mail

type MailSender interface {
	Send(sendTo []string, subject, body string) error
}
