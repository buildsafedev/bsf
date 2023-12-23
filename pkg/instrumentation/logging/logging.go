package logging

import (
	"context"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

type loggerKey struct{}

// SetupLogger sets up the logger.
func SetupLogger() *logrus.Entry {
	// TODO: Change log format to be appropriate for CLI.
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{
		// FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			return "", fileName
		},
	})

	return log.WithFields(logrus.Fields{})
}

// InjectLogger injects the logger into the context.
func InjectLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// WithLogger returns a logger from the context.
func WithLogger(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value(loggerKey{}).(*logrus.Entry)
	if !ok {
		return SetupLogger()
	}
	return logger
}
