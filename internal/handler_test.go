package internal

import (
	"database/sql"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
	DB *sql.DB
}
