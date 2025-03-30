package config

import (
	"be-yourmoments/user-svc/internal/helper/utils"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDB() *sqlx.DB {
	// Load database configuration from environment variables

	dbHost := utils.GetEnv("DB_HOST")
	dbPort := utils.GetEnv("DB_PORT")
	dbUser := utils.GetEnv("DB_USERNAME")
	dbPass := utils.GetEnv("DB_PASSWORD")
	dbName := utils.GetEnv("DB_NAME")
	dbSSLMode := utils.GetEnv("DB_SSLMODE")
	minConns := utils.GetEnv("DB_MIN_CONNS")
	maxConns := utils.GetEnv("DB_MAX_CONNS")
	maxIdleTime := utils.GetEnv("DB_MAX_IDLE_TIME")
	timeZone := utils.GetEnv("TZ")

	// Construct the connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s", dbUser, dbPass, dbHost, dbPort, dbName, dbSSLMode, timeZone)

	// Open the database connection
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Convert environment variables to proper types
	if minConns != "" {
		if v, err := strconv.Atoi(minConns); err == nil {
			db.SetMaxIdleConns(v)
		}
	}

	if maxConns != "" {
		if v, err := strconv.Atoi(maxConns); err == nil {
			db.SetMaxOpenConns(v)
		}
	}

	if maxIdleTime != "" {
		if v, err := time.ParseDuration(maxIdleTime); err == nil {
			db.SetConnMaxIdleTime(v)
		}
	}

	return db
}
