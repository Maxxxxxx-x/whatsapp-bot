package logger

import (
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/Maxxxxxx-x/whatsapp-bot/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	once sync.Once
	log  zerolog.Logger
)

func GetLogger(cfg *config.Config) zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		parsedLevel, err := zerolog.ParseLevel(cfg.LogLevel)
		if err != nil {
			parsedLevel = zerolog.InfoLevel
		}

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		if cfg.AppEnv != "dev" {
			fileLogger := &lumberjack.Logger{
				Filename:   cfg.LogFilePath(),
				MaxSize:    5,
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   true,
			}

			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		}

		var gitRevision string

		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		log = zerolog.New(output).
			Level(parsedLevel).
			With().
			Timestamp().
			Str("git_revision", gitRevision).
			Str("go_version", buildInfo.GoVersion).
			Logger()
	})

	return log
}
