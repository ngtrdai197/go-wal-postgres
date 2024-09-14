package db

import (
	"context"
	"database/sql"
	"go-wal/constant"

	"gorm.io/gorm"
)

func StringToNullString(s string, cannotNil bool) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  cannotNil,
	}
}

func GetTX(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(constant.TxKey).(*gorm.DB)
	if !ok {
		tx = db.WithContext(ctx)
	}
	return tx
}
