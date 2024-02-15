package db

import (
	"context"

	"gorm.io/gorm"
)

const contextKeyDb = "db"

func DbFromContext(ctx context.Context) *gorm.DB {
	if result, ok := ctx.Value(contextKeyDb).(*gorm.DB); ok {
		return result
	}
	return nil
}

func ContextWithDb(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, contextKeyDb, db)
}
