package logger

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Params struct {
	fx.In

	Config Config
}

func New(p Params) zerolog.Logger {
	var writer []io.Writer

	// UNIX Time is faster and smaller than most timestamps
	if p.Config.Format == "json" {
		writer = append(writer, os.Stdout)
	} else {
		cw := &zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}

		writer = append(writer, cw)
	}

	if p.Config.File != "" {
		writer = append(writer, &lumberjack.Logger{
			Filename:   p.Config.File,
			MaxSize:    p.Config.MaxSize, // megabytes
			MaxBackups: p.Config.MaxBackups,
			MaxAge:     p.Config.MaxAge, // days
		})
	}

	level, err := zerolog.ParseLevel(p.Config.Level)
	if err == nil {
		zerolog.SetGlobalLevel(level)
	}

	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.InterfaceMarshalFunc = sonic.Marshal

	// Caller Marshal Function
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	l := zerolog.
		New(zerolog.MultiLevelWriter(writer...)).
		With().
		Str("service", p.Config.ServiceName).
		Timestamp().
		Caller().
		Logger()

	log.Logger = l
	return l
}
