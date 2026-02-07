package main

import (
	"fmt"

	"inventory/internal/config"
	"inventory/internal/db"
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
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal("failed to load config", err)
	}
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
	r := gin.Default()
	db, err := db.NewDBWithSchema(conf.DB)
	err = db.MakeUserAdmin(conf.Admin.GetAdmin())
	if err != nil {
		log.Fatal("failed to add admin", err)
	}
	if err != nil {
		log.Fatal("failed to setup db ", err)
	}
	handlers.Handle(r, db)
	r.Run(conf.API.Addr())
	if err != nil {
		log.Fatal(err)
	}
}
