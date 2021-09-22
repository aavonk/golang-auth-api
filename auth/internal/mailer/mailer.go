package mailer

type Mailer interface {
	Send(recipient, templateFile string, data interface{}) error
}
