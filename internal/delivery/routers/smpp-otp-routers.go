package routers

import (
	"net/http"
	"smpp-otp/internal/delivery/handlers"
	"smpp-otp/internal/service"
	"smpp-otp/pkg/lib/logger"

	"github.com/go-chi/chi/v5"
)

func SetupOTPRoutes(otpService service.OTPService, logger *logger.Loggers) http.Handler {
	otpRouter := chi.NewRouter()
	otpHandler := handlers.NewOTPHandler(otpService)

	otpRouter.Post("/sendOTP", otpHandler.GenerateAndSaveOTPHandler)
	otpRouter.Post("/validateOTP", otpHandler.ValidateOTPHandler)

	return otpRouter
}
