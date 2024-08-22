package service

import (
	"fmt"
	"math/rand"
	"smpp-otp/internal/config"
	smpp "smpp-otp/internal/infrastructure/interfaces"
	repository "smpp-otp/internal/repository/interfaces"
	"smpp-otp/pkg/lib/logger"
	"smpp-otp/pkg/lib/utils"
	"time"
)

type OTPService struct {
	repository repository.OTPRepository
	logger     *logger.Loggers
	cfg        *config.Config
	smppClient smpp.SMPPClient
}

func NewOTPService(repo repository.OTPRepository, logger *logger.Loggers, cfg *config.Config, smppClient smpp.SMPPClient) OTPService {
	return OTPService{repository: repo, logger: logger, cfg: cfg, smppClient: smppClient}
}

func GenerateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := fmt.Sprintf("%06d", r.Intn(1000000))
	return otp
}

func (s *OTPService) SaveAndSendOTP(phoneNumber string) error {
	otp := GenerateOTP()
	err := s.repository.SaveOTP(phoneNumber, otp)
	if err != nil {
		s.logger.ErrorLogger.Error("Error saving OTP to repository: %v", utils.Err(err))
		return err
	}

	err = s.smppClient.SendSMS(s.cfg, phoneNumber, otp)
	if err != nil {
		s.logger.ErrorLogger.Error("Error sending OTP via SMS: %v", utils.Err(err))
		return err
	}

	return nil
}

func (s *OTPService) ValidateOTP(phoneNumber string, otp string) error {
	storedOTP, err := s.repository.GetOTP(phoneNumber)
	if err != nil {
		if err.Error() == "redis: nil" {
			return fmt.Errorf("OTP not found or expired")
		}
		s.logger.ErrorLogger.Error("Error retrieving OTP from repository: %v", utils.Err(err))
		return err
	}

	if storedOTP != otp {
		return fmt.Errorf("OTP does not match")
	}

	return nil
}
