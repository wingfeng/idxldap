package ldap

import (
	"log/slog"
	"strings"

	"github.com/jimlambrt/gldap"
)

func bindHandler(w *gldap.ResponseWriter, r *gldap.Request) {
	resp := r.NewBindResponse(
		gldap.WithResponseCode(gldap.ResultInvalidCredentials),
	)
	defer func() {
		w.Write(resp)
	}()

	m, err := r.GetSimpleBindMessage()
	if err != nil {
		slog.Error("bindHandler Error", "error", err)

		return
	}

	var userName string
	if strings.Contains(m.UserName, "@") {
		userName = strings.Split(m.UserName, "@")[0]
	} else if strings.Contains(m.UserName, "\\") {
		userName = strings.Split(m.UserName, "\\")[1]
	} else if strings.Contains(m.UserName, ",") {
		arr := strings.Split(m.UserName, ",")
		for _, v := range arr {
			if strings.HasPrefix(v, "uid=") {
				userName = v[4:]
				break
			}
		}
	} else {
		userName = m.UserName
	}
	slog.Debug("bindHandler:", "user name", userName)
	if userService.VerifyPassword(userName, string(m.Password)) {
		resp.SetResultCode(gldap.ResultSuccess)
		slog.Debug("LDAP 用户认证成功", "username", m.UserName)
		return
	}
}
