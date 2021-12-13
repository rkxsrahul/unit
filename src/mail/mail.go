package mail

import (
	"log"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/verifyToken"
)

// SendForgotPassMail is a method for sending reset password link in mail
func SendForgotPassMail(account database.Accounts, ws string) {
	// map saving name of user and reset password link for forgot password
	mapd := map[string]interface{}{
		"Name":             account.Name,
		"VerificationCode": "https://" + ws + "." + config.Conf.Address.HostAddress + "/reset-password?token=" + verifyToken.CheckSentToken(account.Userid, "forgot_pass"),
	}

	// readtoml file to fetch template path, subject and images path to be passed in mail
	tmplPath, subject, images := ReadToml("forgotPassword")

	// parse email template
	tmpl := EmailTemplate(tmplPath, mapd)

	//finally send mail
	go SendMail(account.Email, subject, tmpl, images)
}

// SendLoginLink is a method to send login link in mail
func SendLoginLink(account database.Accounts, ws, email, name string) {

	// map saving name of user and reset password link for forgot password
	mapd := map[string]interface{}{
		"Owneremail":       email,
		"Ownername":        name,
		"Workspace":        ws,
		"VerificationCode": "https://" + ws + "." + config.Conf.Address.HostAddress + "/login",
	}

	// readtoml file to fetch template path, subject and images path to be passed in mail
	tmplPath, subject, images := ReadToml("loginLink")

	//dynamic subject
	subject = "You're added to " + ws + " on " + subject
	// parse email template
	tmpl := EmailTemplate(tmplPath, mapd)

	//finally send mail
	go SendMail(account.Email, subject, tmpl, images)
}

// SendInviteLink is a method to send invite link in mail
func SendInviteLink(account database.Accounts, ws, email, name string) {
	log.Println("=--=-=-=-=-=-=-=-=", email, name)
	// map saving name of user and reset password link for forgot password
	mapd := map[string]interface{}{
		"Owneremail":       email,
		"Ownername":        name,
		"Workspace":        ws,
		"VerificationCode": "https://" + ws + "." + config.Conf.Address.HostAddress + "/member-registration?token=" + verifyToken.CheckSentToken(account.Userid, "invite_link"),
	}

	// readtoml file to fetch template path, subject and images path to be passed in mail
	tmplPath, subject, images := ReadToml("inviteLink")

	//dynamic subject
	subject = "You’re invited to join " + ws + " on " + subject

	// parse email template
	tmpl := EmailTemplate(tmplPath, mapd)

	//finally send mail
	go SendMail(account.Email, subject, tmpl, images)
}

// SendRecoveryMail is a method to send recover workspace link in mail
func SendRecoveryMail(account database.Accounts) {
	// map saving name of user and reset password link for forgot password
	mapd := map[string]interface{}{
		"Name":             account.Name,
		"VerificationCode": config.Conf.Address.FrontEndAddress + "/your-workspaces?token=" + verifyToken.CheckSentToken(account.Userid, "recover_workspace"),
	}

	// readtoml file to fetch template path, subject and images path to be passed in mail
	tmplPath, subject, images := ReadToml("recoverWorkspace")

	// parse email template
	tmpl := EmailTemplate(tmplPath, mapd)

	//finally send mail
	go SendMail(account.Email, subject, tmpl, images)
}

// SendVerifyMail is a method for sending verification code in mail
func SendVerifyMail(account database.Accounts) {
	// map saving name of user and verification code for email verification
	mapd := map[string]interface{}{
		"Name":             account.Name,
		"VerificationCode": verifyToken.CheckSentToken(account.Userid, "email_verification"),
	}

	// readtoml file to fetch template path, subject and images path to be passed in mail
	tmplPath, subject, images := ReadToml("verification")

	// parse email template
	tmpl := EmailTemplate(tmplPath, mapd)
	//finally send mail
	go SendMail(account.Email, subject, tmpl, images)
}

// SendForgotPassMail is a method for sending reset password link in mail
func SendForgotPassMailAccount(account database.Accounts) {
	// map saving name of user and reset password link for forgot password
	mapd := map[string]interface{}{
		"Name":             account.Name,
		"VerificationCode": config.Conf.Address.FrontEndAddress + "/reset-password?token=" + verifyToken.CheckSentToken(account.Userid, "forgot_pass"),
	}

	// readtoml file to fetch template path, subject and images path to be passed in mail
	tmplPath, subject, images := ReadToml("forgotPassword")

	// parse email template
	tmpl := EmailTemplate(tmplPath, mapd)
	//finally send mail
	go SendMail(account.Email, subject, tmpl, images)
}
