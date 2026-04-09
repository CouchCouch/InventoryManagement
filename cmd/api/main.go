package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"inventory/internal/auth"
	"inventory/internal/config"
	"inventory/internal/db"
	"inventory/internal/domain"
	"inventory/internal/handlers"
	"inventory/internal/logger"

	"github.com/gin-gonic/gin"
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
	
	// Load config first to get log level
	conf, err := config.GetConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	
	// Initialize logger with config settings
	isDev := logger.DetectEnvironment()
	log := logger.Initialize(conf.LogLevel, isDev)
	slog.SetDefault(log)
	
	slog.Info("Starting Inventory API", "environment", map[bool]string{true: "development", false: "production"}[isDev])
	
	r := gin.Default()
	authService := auth.NewAuthService(conf.Auth.JWTSecret, conf.Auth.JWTRrefreshSecret)
	database, err := db.NewDBWithSchema(conf.DB)
	if err != nil {
		slog.Error("failed to setup database", "error", err)
		os.Exit(1)
	}
	
	err = database.MakeUserAdmin(conf.Admin.GetAdmin())
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			slog.Info("Skipped adding user, user already exists")
		} else {
			slog.Error("failed to add admin", "error", err)
			os.Exit(1)
		}
	}
	
	handlers.Handle(r, database, authService)
	
	slog.Info("Server starting", "address", conf.API.Addr())
	err = r.Run(conf.API.Addr())
	if err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
