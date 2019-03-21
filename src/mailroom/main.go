package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mailgun/mailgun-go"
)

func main() {

	domainName := os.Getenv("MAILGUN_DOMAIN_NAME")
	mailgunAPIKey := os.Getenv("MAILGUN_API_KEY")

	sender := flag.String("sender", "", "The sender's email address.")
	subject := flag.String("subject", "", "The email subject.")
	body := flag.String("body", "", "The email body.")
	recipient := flag.String("recipient", "", "The email recipient(s).")

	mg := mailgun.NewMailgun(domainName, mailgunAPIKey)

	message := mg.NewMessage(*sender, *subject, *body, *recipient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
