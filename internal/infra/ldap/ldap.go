package ldap

import (
	"b8boost/backend/config"
	"fmt"
	"strconv"

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

func GetFirstValueOrDefaultInt(data LDAPUserData, key string, defaultValue int) int {
	if values, exists := data[key]; exists && len(values) > 0 {
		val, err := strconv.Atoi(values[0])
		if err != nil {
			return defaultValue
		}
		return val
	}
	return defaultValue
}

func GetFirstValueOrDefaultPtr(data LDAPUserData, key string, defaultValue *string) *string {
	if values, exists := data[key]; exists && len(values) > 0 {
		return &values[0]
	}
	return defaultValue
}

func (l LDAP) FetchAllUsers() ([]LDAPUserData, error) {
	conn := l.Connect()
	defer conn.Close()
	searchFilter := "(&(objectClass=inetOrgPerson))"
	searchRequest := ldap.NewSearchRequest(
		"ou=users,dc=sso,dc=b8st,dc=ru",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		[]string{"dn", "cn", "description", "uid", "mobile", "sn", "givenName", "createTimestamp", "modifyTimestamp"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("LDAP search error: %w", err)
	}

	var users []LDAPUserData
	for _, entry := range sr.Entries {
		userData := make(LDAPUserData)
		for _, attr := range entry.Attributes {
			userData[attr.Name] = attr.Values
		}
		userData["dn"] = []string{entry.DN}
		users = append(users, userData)
	}

	return users, nil
}
