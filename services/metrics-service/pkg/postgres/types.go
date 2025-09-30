package postgres

import "github.com/jackc/pgx/v5"

type TxPG pgx.Tx

type Config struct {
	Host                 string
	Port                 int
	User                 string
	Password             string
	DBName               string
	SSLMode              string
	PostgresQueryTimeout int
}
