package jobs

import (
	"log/slog"

	"github.com/eng-yousef-khaled/expenses_api/internal/adapters/gomailing"
	"github.com/gocraft/work"
)

type JobHandler struct {
	Mail gomailing.Mailer
}

func (j *JobHandler) ProccessSendMail(job *work.Job) error {
	email := job.ArgString("email")
	message := job.ArgString("message")
	data := map[string]any{"Name": email, "message": message}
	slog.Info("Processing email task after delay", "to", email)
	err := j.Mail.SendMail(gomailing.Request{
		To:       []string{email},
		Subject:  "Code Verification",
		Body:     &message,
		MailType: gomailing.Text,
		Data:     data,
	})
	if err != nil {
		slog.Error("ProcessSendMail Failed", "error", err)
	}
	return err
}
