package mail

type OVHMailConfig struct {
	Server   string
	Port     int16
	Email    string
	Password string
}
type Request struct {
	to       []string
	subject  string
	body     string
	fileName string
	data     interface{}
}

const MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func NewRequest(to []string, body string) *Request {
	return &Request{
		to:   to,
		body: body,
	}
}
