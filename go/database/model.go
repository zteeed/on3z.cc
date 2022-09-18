package database

import (
	"database/sql"
)

type Database struct {
	DB                                          *sql.DB
	StatementSelectShortURLMapping              *sql.Stmt
	StatementSelectShortURLMappingAuthenticated *sql.Stmt
	StatementAddShortURLMapping                 *sql.Stmt
	StatementDeleteAll                          *sql.Stmt
}

type ShortURLMapping struct {
	ShortURL string
	LongURL  string
	Auth0Sub *string
}
