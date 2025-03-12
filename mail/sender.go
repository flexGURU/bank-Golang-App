package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress = "smtp.gmail.com"
	gmailAuthAdd = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		attachments []string,
	) error
}

type GmailSender struct {
	name string
	fromEmail string
	password string

}


func NewGmailSender(name string, fromEmail string, password string) EmailSender  {

	return &GmailSender{
		name: name,
		fromEmail: fromEmail,
		password: password,
	}
	
}


func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	attachments []string,
) error {
	e := email.NewEmail()
	
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmail)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc

	for _, file := range attachments {
		_, err := e.AttachFile(file)
		if err != nil {
			return fmt.Errorf("failed to attach file %s, %w", file, err )
		}
	}

	smtpAUTH :=smtp.PlainAuth("", sender.fromEmail, sender.password, smtpAuthAddress )

	return e.Send(gmailAuthAdd,smtpAUTH)

}