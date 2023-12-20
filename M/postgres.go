package M

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"

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
		Password: "password",
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

func Migrate(db *sql.DB, dir string) error {
	goose.SetDialect("postgres")
	err := goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate :%w", err)
	}
	return nil
}

func MigrateFS(fs embed.FS, db *sql.DB, dir string) error {
	if dir == "" {
		dir = "."
	}
	goose.SetBaseFS(fs)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}
