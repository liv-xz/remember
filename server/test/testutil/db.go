package testutil

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MockDB struct {
	DB    *gorm.DB
	Mock  sqlmock.Sqlmock
	SqlDB *sql.DB
}

func NewMockDB(t *testing.T) *MockDB {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("创建 sqlmock 失败: %v", err)
	}

	d := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})
	gormDB, err := gorm.Open(d, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		sqlDB.Close()
		t.Fatalf("打开 gorm 失败: %v", err)
	}

	return &MockDB{
		DB:    gormDB,
		Mock:  mock,
		SqlDB: sqlDB,
	}
}

func (m *MockDB) Close() {
	m.SqlDB.Close()
}
