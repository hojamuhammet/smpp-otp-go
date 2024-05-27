package smpp

import "smpp-otp/internal/config"

type SMPPClient interface {
	SendSMS(cfg *config.Config, dest, text string) error
}
