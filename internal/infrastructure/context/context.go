package context

import (
	"context"

	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

const (
	ctxKeyLog = "log"
	ctxKeyDb  = "db"
)

// LogFromContext retrieve the identity.XRHID data struct from the context.
func LogFromContext(ctx context.Context) *slog.Logger {
	if data := ctx.Value(ctxKeyLog); data != nil {
		if log, ok := data.(*slog.Logger); ok {
			return log
		}
	}
	return slog.Default()
}

// DBFromContext retrieve the gorm.DB data struct from the context.
func DBFromContext(ctx context.Context) *gorm.DB {
	if data := ctx.Value(ctxKeyDb); data != nil {
		if db, ok := data.(*gorm.DB); ok {
			return db
		}
	}
	return nil
}
