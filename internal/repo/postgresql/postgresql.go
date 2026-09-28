package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	conn *sql.DB
}

func NewPostgresWithURL(url string) (*PostgreSQL, error) {
	psql := new(PostgreSQL)

	connection, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v\n", err)
	}

	psql.conn = connection

	if err := psql.conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping to db: %v\n", err)
	}

	return psql, nil
}
