package mailer

const (
	senderName = "GopherSocial"
	maxRetries = 3
)

type Client interface {
	Send(templateFile, username, email string, data any) error
}
