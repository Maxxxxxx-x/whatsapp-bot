package logger

import (
	"fmt"

	"github.com/rs/zerolog"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WaLogger struct {
	logger zerolog.Logger
	module string
}

func NewWaLogger(parent zerolog.Logger, module string) waLog.Logger {
	return &WaLogger{
		logger: parent.With().Str("module", module).Logger(),
		module: module,
	}
}

func (w *WaLogger) Errorf(msg string, args ...any) {
	w.logger.Error().Msgf(msg, args...)
}

func (w *WaLogger) Warnf(msg string, args ...any) {
	w.logger.Warn().Msgf(msg, args...)
}

func (w *WaLogger) Infof(msg string, args ...any) {
	w.logger.Info().Msgf(msg, args...)
}

func (w *WaLogger) Debugf(msg string, args ...any) {
	w.logger.Debug().Msgf(msg, args...)
}

func (w *WaLogger) Sub(module string) waLog.Logger {
	return &WaLogger{
		logger: w.logger.With().Str("module", fmt.Sprintf("%s/%s", w.module, module)).Logger(),
		module: fmt.Sprintf("%s/%s", w.module, module),
	}
}
