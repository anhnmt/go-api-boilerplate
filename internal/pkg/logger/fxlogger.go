package logger

import (
	"strings"

	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

var _ fxevent.Logger = (*zeroLogger)(nil)

// zeroLogger an Fx event logger that logs events using a zerolog logger.
// https://github.com/ipfans/fxlogger
type zeroLogger struct {
	log zerolog.Logger
}

func NewFxLogger(l zerolog.Logger) fxevent.Logger {
	return &zeroLogger{log: l}
}

func (l *zeroLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	// case *fxevent.OnStartExecuting:
	// l.log.Info().Str("callee", e.FunctionName).
	//     Str("caller", e.CallerName).
	//     Msg("OnStart hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.log.Err(e.Err).
				Str("function", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStart hook failed")
		} else {
			l.log.Info().
				Str("function", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStart hook executed")
		}
	// case *fxevent.OnStopExecuting:
	// l.log.Info().Str("callee", e.FunctionName).
	// 	Str("caller", e.CallerName).
	// 	Msg("OnStop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.log.Err(e.Err).
				Str("function", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStop hook failed")
		} else {
			l.log.Info().
				Str("function", e.FunctionName).
				Str("caller", e.CallerName).
				Str("runtime", e.Runtime.String()).
				Msg("OnStop hook executed")
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.log.Err(e.Err).
				Str("type", e.TypeName).
				Msg("Supplied")
		} else {
			l.log.Info().
				Str("type", e.TypeName).
				Msg("Supplied")
		}
	case *fxevent.Provided:
		if e.Err != nil {
			l.log.Err(e.Err).Msg("Error encountered while applying options")
		} else {
			for _, rtype := range e.OutputTypeNames {
				l.log.Info().
					Str("type", rtype).
					Str("constructor", e.ConstructorName).
					Msg("Provided")
			}
		}
	// case *fxevent.Invoking:
	// Do nothing. Will log on Invoked.
	case *fxevent.Invoked:
		if e.Err != nil {
			l.log.Err(e.Err).
				Str("function", e.FunctionName).
				Msg("Invoke failed")
		} else {
			l.log.Info().
				Str("function", e.FunctionName).
				Msg("Invoked")
		}
	case *fxevent.Stopping:
		l.log.Info().
			Str("signal", strings.ToUpper(e.Signal.String())).
			Msg("Received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.log.Err(e.Err).Msg("Stop failed")
		}
	// case *fxevent.RollingBack:
	// l.log.Err(e.StartErr).Msg("start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.log.Err(e.Err).Msg("Rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.log.Err(e.Err).Msg("Start failed")
		} else {
			l.log.Info().Msg("Started")
		}
		// case *fxevent.LoggerInitialized:
		//     if e.Err != nil {
		//         l.log.Err(e.Err).Msg("custom logger initialization failed")
		//     }
	}
}
