package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type logrusLog struct{}

func NewLogrusLog() Logger {
	return &logrusLog{}
}

// Error exibe detalhes do erro
func (*logrusLog) Error(args ...interface{}) {
	logrus.Error(args...)
}

// ErrorContext exibe detalhes do erro com o contexto
func (*logrusLog) ErrorContext(ctx context.Context, args ...interface{}) {
	logrus.WithContext(ctx).Error(args...)
}

// Info exibe detalhes do erro
func (*logrusLog) Info(args ...interface{}) {
	logrus.Info(args...)
}

// InfoContext exibe detalhes do erro com o contexto
func (*logrusLog) InfoContext(ctx context.Context, args ...interface{}) {
	logrus.WithContext(ctx).Info(args...)
}

// Fatal exibe detalhes do erro
func (*logrusLog) Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}
