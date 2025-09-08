package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/ayayaakasvin/oneflick-ticket/internal/config"
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/inner"
	_ "github.com/lib/pq"
)

const origin = "PostgreSQL"

type PostgreSQL struct {
	conn *sql.DB
}

func NewPostgreSQLConnection(dbConfig config.StorageConfig, shutdownChannel inner.ShutdownChannel) *PostgreSQL {
	psql := new(PostgreSQL)

	connection, err := sql.Open("postgres", connString(dbConfig))
	if err != nil {
		msg := fmt.Sprintf("failed to connect to db: %v\n", err)
		shutdownChannel.Send(inner.ShutdownMessage, origin, msg)
		return nil
	}

	psql.conn = connection

	if err := psql.conn.Ping(); err != nil {
		msg := fmt.Sprintf("failed to ping to db: %v\n", err)
		shutdownChannel.Send(inner.ShutdownMessage, origin, msg)
		return nil
	}

	return psql
}

func connString(dbConfig config.StorageConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DatabaseName)
}
