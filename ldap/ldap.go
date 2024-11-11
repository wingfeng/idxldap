package ldap

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/jimlambrt/gldap"
	"github.com/wingfeng/idx/service"
	"github.com/wingfeng/idxldap/conf"
)

var userService *service.UserService

func StartLdapServer(service *service.UserService) {
	userService = service
	slog.Info("LDAP 服务器配置", "IPAddress", conf.Options.LDAP.IPAddress, "Port", conf.Options.LDAP.Port)
	slog.Info("LDAP 服务器启动成功", "BaseDN", conf.Options.LDAP.BaseDN)
	link := fmt.Sprintf("%s:%d", conf.Options.LDAP.IPAddress, conf.Options.LDAP.Port)
	ldapServer, err := gldap.NewServer()
	if err != nil {
		log.Fatalf("unable to create server: %s", err.Error())
	}
	r, err := gldap.NewMux()
	if err != nil {
		log.Fatalf("unable to create router: %s", err.Error())
	}
	r.Bind(bindHandler)
	r.Search(searchHandler)
	r.Add(func(w *gldap.ResponseWriter, r *gldap.Request) {
		slog.Debug("Add request")
		resp := r.NewResponse(gldap.WithResponseCode(gldap.ResultNotSupported))
		defer func() {
			w.Write(resp)
		}()
	})
	r.Delete(func(w *gldap.ResponseWriter, r *gldap.Request) {
		slog.Debug("Delete request")

		resp := r.NewResponse(gldap.WithResponseCode(gldap.ResultNotSupported))
		defer func() {
			w.Write(resp)
		}()
	})
	r.Unbind(func(w *gldap.ResponseWriter, r *gldap.Request) {
		slog.Debug("Unbind request")
		um, err := r.GetUnbindMessage()
		if err != nil {
			slog.Error("UnbindHandler Error", "error", err)

			return
		}

		slog.Debug("UnbindHandler:", "message", um)
	})
	ldapServer.Router(r)
	ldapServer.Run(link) // listen on port 10389

}
