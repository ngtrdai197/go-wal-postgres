package logger

import (
	"context"
	"fmt"
	"go-wal/constant"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type TracingHook struct{}

const (
	// TraceIDKey is a context key for the request ID
	TraceIDKey string = "trace_id"
)

func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	traceID := GetTraceIDFromContext(ctx) // as per your tracing framework
	if traceID != "" {
		e.Str(TraceIDKey, traceID)
	}
}

func GetTraceIDFromContext(ctx context.Context) string {
	traceID, ok := ctx.Value(constant.XRequestId).(string)
	if !ok {
		return ""
	}
	return traceID
}

func InitGlobalLogger() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Hook(zerolog.HookFunc(func(e *zerolog.Event, level zerolog.Level, msg string) {
		if level == zerolog.ErrorLevel {
			e.Str("stack", fmt.Sprintf("%+v", errors.WithStack(errors.New(msg))))
		}
	})).Hook(TracingHook{})
}

func Info(ctx context.Context) *zerolog.Event {
	return log.Info().Ctx(ctx)
}

func Error(ctx context.Context) *zerolog.Event {
	return log.Error().Ctx(ctx)
}

func Debug(ctx context.Context) *zerolog.Event {
	return log.Debug().Ctx(ctx)
}

func Warn(ctx context.Context) *zerolog.Event {
	return log.Warn().Ctx(ctx)
}

func Fatal(ctx context.Context) *zerolog.Event {
	return log.Fatal().Ctx(ctx)
}

func Panic(ctx context.Context) *zerolog.Event {
	return log.Panic().Ctx(ctx)
}
