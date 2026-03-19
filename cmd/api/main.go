package main

import (
	"errors"
	"fmt"

	"inventory/internal/auth"
	"inventory/internal/config"
	"inventory/internal/db"
	"inventory/internal/domain"
	"inventory/internal/handlers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const inventoryAPIText = `
  _____                      _                              _____ _____
 |_   _|                    | |                       /\   |  __ \_   _|
   | |  _ ____   _____ _ __ | |_ ___  _ __ _   _     /  \  | |__) || |
   | | | '_ \ \ / / _ \ '_ \| __/ _ \| '__| | | |   / /\ \ |  ___/ | |
  _| |_| | | \ V /  __/ | | | || (_) | |  | |_| |  / ____ \| |    _| |_
 |_____|_| |_|\_/ \___|_| |_|\__\___/|_|   \__, | /_/    \_\_|   |_____|
                                            __/ |
                                           |___/
	`

func main() {
	fmt.Println(inventoryAPIText)
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal("failed to load config", err)
	}
	r := gin.Default()
	authService := auth.NewAuthService(conf.Auth.JWTSecret, conf.Auth.JWTRrefreshSecret)
	db, err := db.NewDBWithSchema(conf.DB)
	if err != nil {
		log.Fatal("failed to setup db ", err)
		log.WithField("err", err).Fatal("failed to setup db")
	}
	err = db.MakeUserAdmin(conf.Admin.GetAdmin())
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			log.Info("Skipped adding user, user already exists")
		} else {
			log.WithField("err", err).Fatal("failed to add admin")
		}
	}
	handlers.Handle(r, db, authService)
	r.Run(conf.API.Addr())
	if err != nil {
		log.Fatal(err)
	}
}
