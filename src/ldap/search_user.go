package ldap

import (
	"fmt"
	"strconv"

	"gopkg.in/ldap.v2"

	"git.xenonstack.com/stacklabs/stacklabs-auth/config"
)

//SearchLDAPAccount function used for search account from ldap
func SearchLDAPAccount(mail string) (int, error) {
	ldapPort, err := strconv.Atoi(config.Conf.LDAP.Port)
	if err != nil {
		return 0, err
	}
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", config.Conf.LDAP.Host, ldapPort))
	if err != nil {
		return 0, err
	}
	defer l.Close()

	// First bind with a read only user
	err = l.Bind(config.Conf.LDAP.AdminDN, config.Conf.LDAP.AdminPass)
	if err != nil {
		return 0, err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		config.Conf.LDAP.UserParentDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(|(&(objectClass=*)(uid=%s))(mail=%s))", mail, mail),
		nil,
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return 0, err
	}

	return len(sr.Entries), err
}
