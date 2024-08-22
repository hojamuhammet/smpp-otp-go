package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"smpp-otp/internal/config"
	"smpp-otp/internal/delivery/routers"
	smpp "smpp-otp/internal/infrastructure"
	"smpp-otp/internal/repository"
	"smpp-otp/internal/service"
	db "smpp-otp/pkg/database"
	"smpp-otp/pkg/lib/logger"
	"smpp-otp/pkg/lib/utils"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()

	logger, err := logger.SetupLogger(cfg.Env)
	if err != nil {
		slog.Error("failed to set up logger: %v", utils.Err(err))
		os.Exit(1)
	}

	logger.InfoLogger.Info("Server is up and running")

	database, err := db.InitDB(cfg)
	if err != nil {
		logger.ErrorLogger.Error("failed to initialize database: %v", utils.Err(err))
		os.Exit(1)
	}

	smppClient, err := smpp.NewSMPPClient(cfg)
	if err != nil {
		logger.ErrorLogger.Error("failed to initialize SMPP client: %v", utils.Err(err))
		os.Exit(1)
	}

	repo := repository.NewOTPRepository(database.GetClient(), logger)
	otpService := service.NewOTPService(repo, logger, cfg, smppClient)

	r := routers.SetupOTPRoutes(otpService, logger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		logger.InfoLogger.Info("Shutting down the server gracefully...")
		if err := database.Close(); err != nil {
			logger.ErrorLogger.Error("Error closing database:", utils.Err(err))
		}
		os.Exit(0)
	}()

	err = http.ListenAndServe(cfg.HTTPServer.Address, r)
	if err != nil {
		logger.ErrorLogger.Error("Server failed to start:", utils.Err(err))
	}
}
