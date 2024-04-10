package logger

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type loggerContextKey string

const LoggerContextKey = loggerContextKey("DoContextKey")

func For(ctx context.Context) (*logrus.Entry, error) {
	logger := ctx.Value(LoggerContextKey)
	if logger == nil {
		err := fmt.Errorf("could not retrieve logger")
		return nil, err
	}

	ll, ok := logger.(*logrus.Entry)
	if !ok {
		err := fmt.Errorf("logrus.Entry has wrong type")
		return nil, err
	}
	return ll.WithContext(ctx), nil
}

func WithContext(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, LoggerContextKey, logger)
}
