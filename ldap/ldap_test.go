package ldap

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/go-ldap/ldap/v3"
	"github.com/stretchr/testify/assert"
	"github.com/wingfeng/idx/models"
)

var url = "ldap://127.0.0.1:10389"

func TestConn_Bind(t *testing.T) {
	l, err := ldap.DialURL(url)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = l.Bind("cn=admin,dc=example,dc=com", "password1")
	if err != nil {
		log.Fatal(err)
	}
}

func TestConn_SearchAsync(t *testing.T) {
	l, err := ldap.DialURL(url)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	filter := "(&(|(uid=admin*)(mail=admin*)(cn=admin*)(sn=admin*))(objectclass=inetOrgPerson)(objectclass=organizationalPerson))"
	searchRequest := ldap.NewSearchRequest(
		"dc=example,dc=com", // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,               // The filter to apply
		[]string{"dn", "cn"}, // A list attributes to retrieve
		nil,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := l.SearchAsync(ctx, searchRequest, 64)
	for r.Next() {
		entry := r.Entry()
		fmt.Printf("%s has DN %s\n", entry.GetAttributeValue("cn"), entry.DN)
	}
	err = r.Err()
	assert.NoError(t, err)

}

func TestOp(t *testing.T) {
	//filter := "(&(objectClass=user)(|(sAMAccountName=中文名)(sAMAccountName=test2)))"
	filter := "(&(|(uid=admin*)(mail=admin*)(cn=admin*)(sn=admin*))(objectclass=inetOrgPerson)(objectclass=organizationalPerson))"
	exp, err := parseFilter(filter)
	assert.NoError(t, err)
	if err == nil {
		printTree(exp, "+")
	}

	//	t.Logf("op %s", exp.Value)
	//	t.Logf("operands len（%d） : %v", len(exp.Children), exp.Children)
}
func printTree(expr *models.Expression, indent string) {
	//fmt.Println(indent, "Expression:"+expr.Filter)
	fmt.Printf("%sType: %s, Column: %s, Value: %s, Operator: %s\n", indent, expr.Operator, expr.Column, expr.Value, expr.Operator)
	for _, child := range expr.Children {
		printTree(&child, indent+"  ")
	}
}
