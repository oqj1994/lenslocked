package M

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	UserName string
	Password string
	DBName   string
	SSLMODE  string
}

func DefaultConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		UserName: "root",
		Password: "123456",
		DBName:   "lenslocked",
		SSLMODE:  "disable",
	}
}

func (c PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.Host,
		c.Port, c.UserName, c.Password, c.DBName, c.SSLMODE)
}

func Open(cfg PostgresConfig) (*sql.DB, error) {
	return sql.Open("pgx", cfg.String())
}
