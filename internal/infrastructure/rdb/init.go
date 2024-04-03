package rdb

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"
	"os"
	"poroto.app/poroto/planner/internal/domain/utils"
)

func InitDB(debugMode bool) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s&interpolateParams=%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		"Asia%2FTokyo",
		true,
	)

	if os.Getenv("ENV") != "development" {
		dsn += fmt.Sprintf("&tls=%s", "tidb")

		err := mysql.RegisterTLSConfig("tidb", &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: os.Getenv("DB_HOST"),
		})
		if err != nil {
			return nil, fmt.Errorf("error while registering tls config: %v\n", err)
		}
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	boil.SetDB(db)

	if debugMode {
		logger, err := utils.NewLogger(utils.LoggerOption{
			Tag: "sqlboiler",
		})
		if err != nil {
			return nil, fmt.Errorf("error while creating logger: %v\n", err)
		}

		stdOutLogger, err := zap.NewStdLogAt(logger, zap.InfoLevel)
		if err != nil {
			return nil, fmt.Errorf("error while creating logger: %v\n", err)
		}

		boil.DebugMode = true
		boil.DebugWriter = stdOutLogger.Writer()
	}

	return db, nil
}
