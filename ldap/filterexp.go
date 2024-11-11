package ldap

import (
	"fmt"
	"slices"
	"strings"

	"github.com/wingfeng/idx/models"
)

// Parse ldap search filter to sql expression tree
// "search filter:
// (&(|(uid=admin*)(mail=admin*)(cn=admin*)(sn=admin*))(objectclass=inetOrgPerson)(objectclass=organizationalPerson))"
// "(&(objectClass=user)(|(sAMAccountName=中文名)(sAMAccountName=test2)))"
func parseFilter(filter string) (*models.Expression, error) {
	filter = strings.TrimSpace(filter)
	if !strings.HasPrefix(filter, "(") || !strings.HasSuffix(filter, ")") {
		return nil, fmt.Errorf("invalid filter: %s", filter)
	}
	filter = filter[1 : len(filter)-1]

	stack := &Stack{}
	root := &models.Expression{Operator: "and"}
	root.Filter = filter
	//	current := root
	deep := 0

	arr := []byte(filter)
	for _, c := range arr {
		stack.Push(c)
		if c == '(' {
			deep++
		}
		if c == ')' {
			deep--

			if deep == 0 {
				tmp := make([]byte, 0)
				for {
					v := stack.Pop()
					if v == ')' {
						deep++
					}
					if v == '(' {
						deep--
					}
					tmp = append(tmp, v)
					if deep == 0 || stack.Len() == 0 {
						break
					}

				}
				slices.Reverse(tmp)
				subFilter := string(tmp)

				subExpr, err := parseFilter(subFilter)
				if err != nil {
					return nil, err
				}
				//	subExpr.Parent = current
				root.Children = append(root.Children, *subExpr)
			}
		}
	}

	if strings.HasPrefix(filter, "&") {
		root.Operator = "and"
		//	filter = filter[1:]
	} else if strings.HasPrefix(filter, "|") {
		root.Operator = "or"
		//	filter = filter[1:]
	} else if strings.HasPrefix(filter, "!") {
		root.Operator = "not"
		//	filter = filter[1:]
	} else {
		// Leaf node
		parts := strings.SplitN(filter, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid leaf filter: %s", filter)
		}
		root.Column = mapColumnName(parts[0])
		val := parts[1]
		if strings.HasPrefix(val, "*") || strings.HasSuffix(val, "*") {
			val = strings.Replace(val, "*", "%", -1)
			root.Value = val
			root.Operator = "like"
		} else {
			root.Value = val
			root.Operator = "="
		}

	}

	return root, nil
}

//	func (e *Expression) IsUserRelated() bool {
//		for _, s := range e.Children {
//			if s.Column == "objectclass" && (s.Value == "inetOrgPerson" || s.Value == "organizationalPerson") {
//				return true
//			}
//		}
//		return false
//	}
func mapColumnName(columnName string) string {
	switch columnName {
	case "entryUUID":
		return "id"
	case "uid":
		return "user_name"
	case "cn":
		return "display_name"
	case "sn":
		return "user_name"
	case "mail":
		return "email"
	case "sAMAccountName":
		return "user_name"
		// case "userPassword":
		// 	return "password"
	}
	return columnName
}
