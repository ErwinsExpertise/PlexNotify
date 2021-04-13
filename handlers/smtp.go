package handlers

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
)

func sendSMTP(content string) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	addr := host + ":" + port

	fromName := os.Getenv("SMTP_USER")
	target := os.Getenv("SMTP_TO_EMAIL")

	fromEmail := fromName + "@" + host

	toNames := []string{target}
	toEmails := []string{target}

	subject := "Plex Server Notification:"
	body := content

	// Build RFC-2822 email
	toAddresses := []string{}
	for i, _ := range toEmails {
		to := mail.Address{toNames[i], toEmails[i]}
		toAddresses = append(toAddresses, to.String())
	}
	toHeader := strings.Join(toAddresses, ", ")
	from := mail.Address{fromName, fromEmail}
	fromHeader := from.String()
	subjectHeader := subject
	header := make(map[string]string)

	header["To"] = toHeader
	header["From"] = fromHeader
	header["Subject"] = subjectHeader
	header["Content-Type"] = `text/plain; charset="UTF-8"`

	msg := ""
	for k, v := range header {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body
	bMsg := []byte(msg)
	// Send using local postfix service
	c, err := smtp.Dial(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	if err = c.Mail(fromHeader); err != nil {
		fmt.Println(err)
		return
	}
	for _, addr := range toEmails {
		if err = c.Rcpt(addr); err != nil {
			fmt.Println(err)
			return
		}
	}
	w, err := c.Data()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = w.Write(bMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = w.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.Quit()

	if err != nil {
		fmt.Println(err)
		return
	}
}
