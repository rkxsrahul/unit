package ldap

import (
	"fmt"
	"strconv"

	"gopkg.in/ldap.v2"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
)

//UpdateLDAPAccountPass function for update the account info on ldap
func UpdateLDAPAccountPass(email, name, password string) error {
	ldapPort, err := strconv.Atoi(config.Conf.LDAP.Port)
	if err != nil {
		return err
	}
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", config.Conf.LDAP.Host, ldapPort))
	if err != nil {
		return err
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(config.Conf.LDAP.AdminDN, config.Conf.LDAP.AdminPass)
	if err != nil {
		return err
	}

	modifyReq := ldap.NewModifyRequest(fmt.Sprintf("uid=%s,%s", email, config.Conf.LDAP.UserParentDN))
	if password != "" {
		modifyReq.Replace("userPassword", []string{password})
	}
	modifyReq.Replace("cn", []string{name})
	err = l.Modify(modifyReq)

	return err
}
