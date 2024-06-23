package store

import (
	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/config"
	"github.com/JacobRWebb/InventoryManagement.Users.Api/pkg/store/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func open(cfg *config.Config) (db *gorm.DB, err error) {
	return gorm.Open(postgres.Open(cfg.DatabaseDSN))
}

func MustOpen(cfg *config.Config) *gorm.DB {
	db, err := open(cfg)

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&user.User{})

	if err != nil {
		panic(err)
	}

	return db
}
