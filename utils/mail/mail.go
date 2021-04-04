package mail

type MailSender interface {
	Send(sendTo []string, msg []byte) error
}
