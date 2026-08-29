package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Maxxxxxx-x/whatsapp-bot/internal/config"
	"github.com/Maxxxxxx-x/whatsapp-bot/internal/logger"

	"github.com/mdp/qrterminal/v3"
	"github.com/rs/zerolog"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	_ "modernc.org/sqlite"
)

type Service struct {
	Client    *whatsmeow.Client
	StartTime time.Time
	db        *sql.DB
	log       zerolog.Logger
}

func NewService(ctx context.Context, cfg *config.Config, log zerolog.Logger) (*Service, error) {
	waLogger := logger.NewWaLogger(log, "whatsmeow")

	if err := os.MkdirAll(filepath.Dir(cfg.WhatsappDBPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	dsn := fmt.Sprintf("%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", cfg.WhatsappDBPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open whatsapp db: %w", err)
	}

	db.SetMaxOpenConns(1)

	container := sqlstore.NewWithDB(db, "sqlite", waLogger)
	if err := container.Upgrade(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to upgrade whatsapp db schema: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to get device store: %w", err)
	}

	client := whatsmeow.NewClient(deviceStore, waLogger)

	return &Service{
		Client:    client,
		StartTime: time.Now().Add(-5 * time.Second),
		db:        db,
		log:       log,
	}, nil
}

func (s *Service) RegisterEventHandler(handler func(evt any)) {
	s.Client.AddEventHandler(handler)
}

func (s *Service) Connect(ctx context.Context) error {
	if s.Client.Store.ID == nil {
		qrChan, err := s.Client.GetQRChannel(ctx)
		if err != nil {
			return fmt.Errorf("Failed to get QR channel: %w", err)
		}

		if err := s.Client.Connect(); err != nil {
			return fmt.Errorf("failed to connect client: %w", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("\nScan QR code:")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				s.log.Info().Str("event", evt.Event).Msg("QR Event received")
			}
		}
	} else {
		if err := s.Client.Connect(); err != nil {
			return fmt.Errorf("failed to reconnect client: %w", err)
		}
		s.log.Info().Msg("reconnected to whatsapp using saved session")
	}
	return nil
}

func (s *Service) Close() {
	s.Client.Disconnect()
	if s.db != nil {
		_ = s.db.Close()
	}
	s.log.Info().Msg("whatsapp service stopped")
}
