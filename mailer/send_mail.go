package mailer

import (
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMail(sender string, senderMail string, recepient string, recepientMail string, subject string, content string) {
	from := mail.NewEmail(sender, senderMail)
	to := mail.NewEmail(recepient, recepientMail)
	message := mail.NewSingleEmail(from, subject, to, content, content)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(response.StatusCode)
		log.Println(response.Body)
		log.Println(response.Headers)
	}
}
