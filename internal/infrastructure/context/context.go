package context

import (
	"context"
	"errors"

	"log/slog"

	common_err "github.com/avisiedo/go-microservice-1/internal/errors/common"
	"gorm.io/gorm"
)

type contextKey int

const (
	ctxKeyLog contextKey = 1000
	ctxKeyDB  contextKey = 1001
)

// LogFromContext retrieve the identity.XRHID data struct from the context.
func LogFromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}
	if data := ctx.Value(ctxKeyLog); data != nil {
		if log, ok := data.(*slog.Logger); ok {
			return log
		}
	}
	return slog.Default()
}

// DBFromContext retrieve the gorm.DB data struct from the context.
func DBFromContext(ctx context.Context) (*gorm.DB, error) {
	if ctx == nil {
		return nil, common_err.ErrNil("ctx")
	}
	if data := ctx.Value(ctxKeyDB); data != nil {
		if db, ok := data.(*gorm.DB); ok {
			return db, nil
		}
	}
	return nil, errors.New("database not found in context")
}

// WithDB create a new context with the provided database value.
func WithDB(ctx context.Context, value *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKeyDB, value)
}

// WithLog create a new context with the provided log value.
func WithLog(ctx context.Context, value *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLog, value)
}
