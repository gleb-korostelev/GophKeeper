package initConnection

import (
	"database/sql"

	"github.com/gleb-korostelev/GophKeeper/config"
	"github.com/gleb-korostelev/GophKeeper/tools/closer"
	"github.com/gleb-korostelev/GophKeeper/tools/db"
	"github.com/gleb-korostelev/GophKeeper/tools/logger"
	"github.com/pressly/goose/v3"
)

func NewDBConn() db.IAdapter {
	dsn := config.GetConfigString(config.DBDSN)

	dbGoose, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		logger.Fatalf("can't open database goose driver adapter: %w", err.Error())
	}

	err = dbGoose.Ping()
	if err != nil {
		logger.Fatalf("can not ping database: %w", err.Error())
	}
	// Example of closer use case
	closer.Add(dbGoose)

	ad, err := db.NewAdapter(dbGoose, sql.LevelReadUncommitted)
	if err != nil {
		logger.Fatal("can't initialize database adapter: %w", err.Error())
	}
	return ad
}
