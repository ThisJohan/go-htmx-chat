package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (p PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.Database, p.SSLMode)
}

func OpenDB(cfg PostgresConfig) (*sqlx.DB, error) {
	fmt.Println(cfg.String())
	return sqlx.Open("postgres", cfg.String())
}
