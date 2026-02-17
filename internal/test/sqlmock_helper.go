package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewSqlMock create a sqlmock associated to the gorm.DB
// session parameters.
// Return A sqlmock and gorm.DB reference on success, or nil, nil and an
// error with the description about what happened.
//
// Example
//
//	mockDB, db, err := test.NewSqlMock(&gorm.Session{SkipHooks: true})
func NewSqlMock(session *gorm.Session) (sqlmock.Sqlmock, *gorm.DB, error) {
	sqlDB, sqlMock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		TranslateError:         true,
	})
	if err != nil {
		return nil, nil, err
	}
	if session != nil {
		gormDB = gormDB.Session(session)
	}

	return sqlMock, gormDB, nil
}
