package database

import (
	"log/slog"

	"github.com/VysMax/organizational-structure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func New(cfg *config.Config, log *slog.Logger) (*Db, error) {

	log.Info("Connecting to database...")

	gormDb, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Info("Connecction to database successful")

	return &Db{gormDb}, nil
}
