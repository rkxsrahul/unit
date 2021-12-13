package ldap

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/ldap.v2"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
	"git.xenonstack.com/stacklabs/stacklabs-auth/src/database"
)

//CreateLDAPAccount function for create new account on ldap
func CreateLDAPAccount(acs database.Accounts, password string) error {

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
	names := strings.Split(acs.Name, " ")
	if len(names) < 2 {
		names = append(names, acs.Name)
	}
	addReq := ldap.NewAddRequest(fmt.Sprintf("uid=%s,%s", acs.Email, config.Conf.LDAP.UserParentDN))
	addReq.Attribute("objectClass", []string{"top", "inetOrgPerson", "posixAccount"})
	addReq.Attribute("userPassword", []string{password})
	addReq.Attribute("cn", []string{"neuron_labs-" + names[0]})
	addReq.Attribute("sn", []string{names[1]})
	addReq.Attribute("mail", []string{acs.Email})
	addReq.Attribute("uidNumber", []string{acs.Userid})
	addReq.Attribute("homeDirectory", []string{"/home/" + acs.Userid})
	addReq.Attribute("gidNumber", []string{"1001"})
	addReq.Attribute("employeeType", []string{config.GroupName})
	err = l.Add(addReq)
	if err != nil {
		return err
	}

	return nil
}

//CreateLDAPGroup function used for create group on ldap
func CreateLDAPGroup(name, id string) error {
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

	addReq := ldap.NewAddRequest(fmt.Sprintf("cn=%s,%s", name, config.Conf.LDAP.GroupBase))
	addReq.Attribute("objectClass", []string{"top", "posixGroup"})
	addReq.Attribute("gidNumber", []string{id})
	err = l.Add(addReq)
	if err != nil {
		return err
	}

	return nil
}
