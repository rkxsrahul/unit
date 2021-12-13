package mail

import (
	"log"
	"strconv"

	gomail "gopkg.in/gomail.v2"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
)

// ToggleMail is a function to change the value of Mail Service in config
func ToggleMail(value string) {
	config.MailService = value
}

// SendMail is a function for sending mail using smtp credentials
func SendMail(to, sub, template string, images []string) {
	//update Configuration
	config.SetConfig()
	// creating new message with default settings
	m := gomail.NewMessage()

	// setting mail headers from, to and subject
	m.SetHeader("From", config.Conf.Mail.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", sub)

	//path is from where main.go is running
	// embedding static images
	for i := 0; i < len(images); i++ {
		m.Embed(images[i])
	}

	// set body of mail
	m.SetBody("text/html", template)

	// port of smtp mail
	port, _ := strconv.Atoi(config.Conf.Mail.Port)
	//use port 465 for TLS, other than 465 it will send without TLS.
	// connect to smtp server using mail admin username and password
	d := gomail.NewPlainDialer(config.Conf.Mail.Host, port, config.Conf.Mail.User, config.Conf.Mail.Pass)

	if port == 465 {
		log.Println("ues")
		d.SSL = true
	}
	if config.MailService != "true" {
		// send above mail message
		err := d.DialAndSend(m)
		log.Println(err)
	}
}
