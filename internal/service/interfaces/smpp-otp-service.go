package service

import (
	"smpp-otp/internal/config"
)

type OTPService interface {
	SaveAndSendOTP(cfg config.Config, phoneNumber string) error
	ValidateOTP(phoneNumber string, otp string) error
}
