package main

import (
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	host := "www.google.com"
	port := "443"
	timeout := time.Duration(1 * time.Second)

	try := 0
	for {
		_, err := net.DialTimeout("tcp", host+":"+port, timeout)
		if err != nil {
			fmt.Printf("%s %s %s\n", host, "not responding", err.Error())
		} else {
			fmt.Printf("%s %s %s\n", host, "responding on port:", port)
		}
		try += 1
		if try > 7 {
			break
		}
	}

	mail := os.Getenv("MAIL_TO")
	fmt.Println(mail)
	err := SendEmail([]string{mail}, "", []byte{})
	if err != nil {
		panic(err)
	}
}

func SendEmail(recp []string, subject string, body []byte) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	from := "From: " + os.Getenv("MAIL_FROM") + "\n"
	to := "To: " + strings.Join(recp, ",") + "\n"
	sub := "Subject: " + subject + "\n"
	headers := from + to + sub + mime + "\n"

	msg := []byte(headers + string(body))

	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	portStr := os.Getenv("MAIL_PORT")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", username, password, host)

	err = smtp.SendMail(addr, auth, os.Getenv("MAIL_FROM"), recp, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
