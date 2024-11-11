package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/wingfeng/idx/repo"
	"github.com/wingfeng/idx/service"
	"github.com/wingfeng/idxldap/conf"
	"github.com/wingfeng/idxldap/ldap"
)

func main() {
	//配置Log
	logLevel := slog.LevelWarn
	switch strings.ToLower(conf.Options.LogLevel) {
	case "debug":
		logLevel = slog.LevelDebug

	case "info":
		logLevel = slog.LevelInfo

	case "warn":
		logLevel = slog.LevelWarn

	}
	slog.Info("Set log level", "Level", logLevel)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(handler))
	userRepo := repo.NewUserRepository()

	us := service.NewUserService(userRepo)
	ldap.StartLdapServer(us)
	select {}
}
