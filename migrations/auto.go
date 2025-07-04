package migrations

import (
	"os"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CheckForMigrations() error {

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		logger.Info("🚀 Starting migrations...")
		if err := RunMigrations(); err != nil {
			logger.Errorf("Error processing migrations: %v", err)
		}
		return nil
	}
	// if not "migrate" args[1]
	return nil
}

func RunMigrations() error {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		//DisableForeignKeyConstraintWhenMigrating: true, //временно игнорировать миграции в первый раз а потом их добавить
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate()
	if err != nil {
		return err
	}

	logger.Info("✅")
	return nil
}
