package vm

import (
	"fmt"
	"strings"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/mail"
)

func VMRequest(info database.VMRequestInfo) error {

	// map saving name of user and reset password link for forgot password
	mapd := map[string]interface{}{
		"Useremail":    info.UserEmail,
		"Username":     info.UserName,
		"Name":         info.Name,
		"Description":  info.Description,
		"Source":       info.Source,
		"Flavour":      info.Flavour,
		"InboundPorts": info.Ports,
		"Workspace":    info.Workspace,
	}

	db := config.DB

	db.Create(&info)

	// readtoml file to fetch template path, subject and images path to be passed in mail
	tmplPath, subject, images := mail.ReadToml("vmcreation")

	subject = subject + " | " + fmt.Sprint(info.ID)
	// parse email template
	tmpl := mail.EmailTemplate(tmplPath, mapd)
	list := strings.Split(config.Conf.Service.SupportEmails, ",")
	for i := 0; i < len(list); i++ {
		//finally send mail
		go mail.SendMail(list[i], subject, tmpl, images)
	}

	return nil
}
