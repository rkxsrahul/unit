package ldap

import (
	"fmt"
	"log"
	"strconv"

	"gopkg.in/ldap.v2"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
)

//DeleteLDAPAccount function used for delete account from ldap
func DeleteLDAPAccount(email string) error {
	ldapPort, err := strconv.Atoi(config.Conf.LDAP.Port)
	if err != nil {
		log.Println(err)
		return err
	}
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", config.Conf.LDAP.Host, ldapPort))
	if err != nil {
		log.Println(err)
		return err
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(config.Conf.LDAP.AdminDN, config.Conf.LDAP.AdminPass)
	if err != nil {
		log.Println(err)
		return err
	}

	delReq := ldap.NewDelRequest(fmt.Sprintf("uid=%s,%s", email, config.Conf.LDAP.UserParentDN), nil)
	err = l.Del(delReq)
	log.Println(err)

	return err
}
