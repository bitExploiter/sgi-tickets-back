package toolbox

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"strconv"
	"sync"

	"gopkg.in/gomail.v2"
)

var mailWaitGroup sync.WaitGroup

func sendEmail(mailTo string, subject string, body string) {
	port, _ := strconv.Atoi(os.Getenv("MAIL_SERVER_PORT"))
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("MAIL_SERVER_FROM"))
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(os.Getenv("MAIL_SERVER_HOST"), port, os.Getenv("MAIL_SERVER_USER"), os.Getenv("MAIL_SERVER_PASSWORD"))
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
	}
}

func SendEmailAsync(mailTo string, subject string, body string) {
	mailWaitGroup.Add(1)
	go func() {
		sendEmail(mailTo, subject, body)
		mailWaitGroup.Done()
	}()
}

func SendNotificacionEmail(nombre string, mailTo string, asunto string, templatePath string, data interface{}) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println(err)
		return
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Println(err)
		return
	}
	SendEmailAsync(mailTo, asunto, body.String())
}
