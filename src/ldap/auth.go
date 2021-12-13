package ldap

import (
	"fmt"
	"log"
	"strconv"

	"gopkg.in/ldap.v2"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
)

//Authentication function used for Authenticate ldap credentails
func Authentication(mailOrUID, password string) (bool, string, string) {

	ldapPort, err := strconv.Atoi(config.Conf.LDAP.Port)
	if err != nil {
		return false, err.Error(), ""
	}

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", config.Conf.LDAP.Host, ldapPort))
	if err != nil {
		return false, err.Error(), ""
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(config.Conf.LDAP.AdminDN, config.Conf.LDAP.AdminPass)
	if err != nil {
		return false, err.Error(), ""
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		config.Conf.LDAP.UserParentDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(|(&(objectClass=*)(uid=%s))(mail=%s))", mailOrUID, mailOrUID),
		nil,
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return false, err.Error(), ""
	}

	if len(sr.Entries) != 1 {
		log.Println("Entries length: ", len(sr.Entries))
		return false, "Invalid Username or Password.", ""
	}

	// Bind as the user to verify their password
	err = l.Bind(sr.Entries[0].DN, password)
	if err != nil {
		log.Println(err)
		return false, "Invalid Username or Password.", ""
	}

	return true, "User Logged in successfully.", sr.Entries[0].GetAttributeValue("uid")
}
