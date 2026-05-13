package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"text/template"
)

type Service interface {
	SendMail(request Request) bool
}

type svc struct {
	mailConfig OVHMailConfig
}

func CreateService(config OVHMailConfig) Service {
	return &svc{
		mailConfig: config,
	}
}
func (s *svc) SendMail(request Request) bool {
	s.parseTemple(&request)
	return s.sendMail(request)
}

func (s *svc) parseTemple(r *Request) error {
	t, err := template.ParseFiles(r.fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err := t.Execute(buffer, r.data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (s *svc) sendMail(r Request) bool {
	SMTP := fmt.Sprintf("%s:%d", s.mailConfig.Server, s.mailConfig.Port)
	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.mailConfig.Server,
	}

	conn, err := tls.Dial("tcp", SMTP, tlsConfig)
	if err != nil {
		fmt.Printf("Connection error: %s", err)
		return false
	}
	defer conn.Close()

	// 3. Create the SMTP client
	client, err := smtp.NewClient(conn, s.mailConfig.Server)
	if err != nil {
		fmt.Printf("SMTP client error: %s", err)
		return false
	}
	defer client.Quit()

	// 4. Authenticate
	auth := smtp.PlainAuth("", s.mailConfig.Email, s.mailConfig.Password, s.mailConfig.Server)
	if err = client.Auth(auth); err != nil {
		fmt.Printf("Auth error: %s", err)
		return false
	}

	// 5. Send the mail commands
	if err = client.Mail(s.mailConfig.Email); err != nil {
		return false
	}
	for _, addr := range r.to {
		if err = client.Rcpt(addr); err != nil {
			return false
		}
	}

	// 6. Write the body
	w, err := client.Data()
	if err != nil {
		return false
	}
	_, err = w.Write([]byte(body))
	if err != nil {
		return false
	}
	err = w.Close()

	if err != nil {
		fmt.Printf("Sending mail result is: %s", err)
		return false
	}

	fmt.Println("Mail sent successfully!")
	return true
}
