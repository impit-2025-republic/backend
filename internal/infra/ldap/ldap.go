package ldap

import (
	"b8boost/backend/config"
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

type LDAP struct {
	LDAPUserFilter string
	LDAPBaseDN     string
	LDAPServer     string
	LDAPPort       string
	LDAPBindDN     string
	LDAPBindPass   string
}

func NewLDAP(conf config.Config) LDAP {
	return LDAP{
		LDAPServer:   conf.LDAPServer,
		LDAPPort:     conf.LDAPPort,
		LDAPBindDN:   conf.LDAPBindDN,
		LDAPBindPass: conf.LDAPBindPass,
	}
}

func (la LDAP) Connect() *ldap.Conn {
	l, err := ldap.DialURL(fmt.Sprintf("%s:%s", la.LDAPServer, la.LDAPPort))
	if err != nil {
		panic(err)
	}

	err = l.Bind(la.LDAPBindDN, la.LDAPBindPass)
	if err != nil {
		panic(err)
	}
	return l
}

type LDAPUserData map[string][]string

func (l LDAP) FindTelegramID(tgID int64) LDAPUserData {
	conn := l.Connect()
	searchFilter := fmt.Sprintf("(&(objectClass=inetOrgPerson)(description=%s))", ldap.EscapeFilter(fmt.Sprintf("%d", tgID)))
	searchRequest := ldap.NewSearchRequest(
		"ou=users,dc=sso,dc=b8st,dc=ru",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		[]string{"dn", "cn", "description", "uid"},
		nil,
	)
	defer conn.Close()
	sr, err := conn.Search(searchRequest)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if len(sr.Entries) != 1 {
		fmt.Println("zero")
		return nil
	}

	userData := make(LDAPUserData)
	for _, attr := range sr.Entries[0].Attributes {
		userData[attr.Name] = attr.Values
	}

	userData["dn"] = []string{sr.Entries[0].DN}

	return userData
}

func GetFirstValueOrDefault(data LDAPUserData, key, defaultValue string) string {
	if values, exists := data[key]; exists && len(values) > 0 {
		return values[0]
	}
	return defaultValue
}
