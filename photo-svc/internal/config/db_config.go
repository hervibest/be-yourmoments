package config

import (
	logs "be-yourmoments/photo-svc/internal/helper/logger"
	"be-yourmoments/photo-svc/internal/helper/utils"
	"log"

	"fmt"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
)

var (
	dbHost    = utils.GetEnv("DB_HOST")
	dbPort    = utils.GetEnv("DB_PORT")
	dbUser    = utils.GetEnv("DB_USERNAME")
	dbPass    = utils.GetEnv("DB_PASSWORD")
	dbSSLMode = utils.GetEnv("DB_SSLMODE")
	dbName    = utils.GetEnv("DB_NAME")
	minConns  = utils.GetEnv("DB_MIN_CONNS")
	maxConns  = utils.GetEnv("DB_MAX_CONNS")
	// TimeOutDuration, _ = strconv.Atoi(utils.GetEnv("DB_CONNECTION_TIMEOUT"))
	maxIdleTime = utils.GetEnv("DB_MAX_IDLE_TIME")
	timeZone    = utils.GetEnv("TZ")
)

func NewPostgresDatabase() *sqlx.DB {
	logger := logs.New("database_connection")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s", dbUser, dbPass, dbHost, dbPort, dbName, dbSSLMode, timeZone)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

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

	logger.Log("Database connected on " + dsn)

	return db
}
