package smpp

import (
	"log"
	"log/slog"
	"os"
	"smpp-otp/internal/config"
	"smpp-otp/pkg/lib/logger"
	utils "smpp-otp/pkg/lib/status"
	"time"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

type SMPPClient struct {
	Transceiver *smpp.Transceiver
}

func NewSMPPClient(cfg *config.Config) (*SMPPClient, error) {
	logger, err := logger.SetupLogger(cfg.Env)
	if err != nil {
		slog.Error("failed to set up logger: %v", utils.Err(err))
		os.Exit(1)
	}

	smppCfg := cfg.SMPP
	var tx *smpp.Transceiver

	for {
		tx = &smpp.Transceiver{
			Addr:    smppCfg.Addr,
			User:    smppCfg.User,
			Passwd:  smppCfg.Pass,
			Handler: func(p pdu.Body) {},
		}

		connStatus := tx.Bind()
		status := <-connStatus
		if status.Status() == smpp.Connected {
			logger.InfoLogger.Info("Connected to SMPP server.")
			slog.Info("Connected to SMPP server.")
			break
		} else {
			logger.ErrorLogger.Error("Failed to connect:", utils.Err(status.Error()))
			slog.Error("Failed to establish OTP connection:", utils.Err(status.Error()))
			logger.InfoLogger.Info("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
		}
	}

	return &SMPPClient{Transceiver: tx}, nil
}

func (c *SMPPClient) SendSMS(cfg *config.Config, dest, text string) error {
	params := &smpp.ShortMessage{
		Src:      cfg.Src_Phone_Number,
		Dst:      dest,
		Text:     pdutext.Raw(text),
		Register: 0, // No delivery receipt
	}
	resp, err := c.Transceiver.Submit(params)
	if err != nil {
		log.Println("Failed to send SMS:", err)
		return err
	}

	log.Println("SMS sent successfully, response:", resp)
	return nil
}
