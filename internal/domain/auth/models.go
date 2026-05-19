package auth

type RequestType int

const (
	HTML RequestType = iota
	Text
)

type Request struct {
	To       []string
	Subject  string
	Body     *string
	TempFile *string
	Data     interface{}
	MailType RequestType
}

type VerificationCodeMail struct {
	Email   string
	Message string
}
