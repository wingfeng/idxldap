package ldap

import (
	"fmt"
	"log/slog"

	"github.com/jimlambrt/gldap"
	"github.com/wingfeng/idxldap/conf"
)

func searchHandler(w *gldap.ResponseWriter, r *gldap.Request) {
	resp := r.NewSearchDoneResponse()
	defer func() {
		w.Write(resp)
	}()
	m, err := r.GetSearchMessage()
	if err != nil {
		slog.Error("searchHandler Error", "error", err)

		return
	}
	slog.Debug("searchHandler:", "message", m)

	slog.Debug("search scope:", "scope", m.Scope, "filter", m.Filter)

	exp, err := parseFilter(m.Filter)
	if err != nil {
		slog.Error("Parse filter Error", "error", err)
	}

	users, err := userService.SearchUsers(exp, 0, 10)
	if err != nil {
		slog.Error("search users from db error", "error", err)
		resp.SetResultCode(gldap.ResultLocalError)
		return
	}
	slog.Debug("searching for users", "users", users)

	for _, user := range users {
		id := fmt.Sprintf("uid=%s,ou=people,%s", user.GetUserName(), conf.Options.LDAP.BaseDN)
		entry := r.NewSearchResponseEntry(id, gldap.WithAttributes(map[string][]string{
			"objectClass":       {"top", "person", "organizationalPerson", "inetOrgPerson"},
			"cn":                {user.GetUserName()},
			"sn":                {user.GetId()},
			"mail":              {user.GetEmail()},
			"userPassword":      {user.GetPasswordHash()},
			"uid":               {user.GetUserName()},
			"entryUUID":         {user.GetId()},
			"krb5PrincipalName": {user.GetUserName()},
			"createdTimeStamp":  {"2023-01-01T00:00:00Z"},
			"modifyTimeStamp":   {"2023-01-01T00:00:00Z"},
		}))
		w.Write(entry)
	}

	resp.SetResultCode(gldap.ResultSuccess)

	// if m.BaseDN == "ou=users,dc=idx,dc=com" {
	// 	entry := r.NewSearchResponseEntry(
	// 		"ou=people,cn=example,dc=org",
	// 		gldap.WithAttributes(map[string][]string{
	// 			"objectclass": {"organizationalUnit"},
	// 			"ou":          {"people"},
	// 		}),
	// 	)
	// 	w.Write(entry)
	// 	resp.SetResultCode(gldap.ResultSuccess)
	// }

}
