package config

import (
	"fmt"
	"log"
	"time"

	"ecopoints-go-api/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *Config) (*gorm.DB, error) {
	// DSN format: username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBDatabase,
	)

	logLevel := logger.Warn
	if cfg.AppEnv == "development" || cfg.AppEnv == "local" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %w", err)
	}

	// Set connection pool parameters - kept low for shared hosting limits
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Database connection to Laravel MySQL database successfully established!")
	if err := db.AutoMigrate(&model.User{}, &model.WasteDeposit{}, &model.WasteDepositItem{}, &model.Reward{}, &model.RewardRedemption{}); err != nil {
		return nil, fmt.Errorf("failed to update application schema: %w", err)
	}

	// Backfill existing waste_deposits to waste_deposit_items if old columns exist
	if db.Migrator().HasColumn("waste_deposits", "waste_type_id") {
		log.Println("Migrating legacy single-item waste_deposits into waste_deposit_items...")
		_ = db.Exec(`
			INSERT INTO waste_deposit_items (waste_deposit_id, waste_type_id, weight_kg, original_weight_kg, actual_weight_kg, created_at, updated_at)
			SELECT id, waste_type_id, weight_kg, original_weight_kg, actual_weight_kg, created_at, updated_at
			FROM waste_deposits
			WHERE id NOT IN (SELECT DISTINCT waste_deposit_id FROM waste_deposit_items WHERE waste_deposit_id IS NOT NULL)
		`).Error

		// Drop legacy columns from waste_deposits
		log.Println("Dropping legacy columns from waste_deposits...")
		_ = db.Migrator().DropColumn("waste_deposits", "waste_type_id")
		_ = db.Migrator().DropColumn("waste_deposits", "weight_kg")
		_ = db.Migrator().DropColumn("waste_deposits", "original_weight_kg")
		_ = db.Migrator().DropColumn("waste_deposits", "actual_weight_kg")
	}

	if err := db.Model(&model.Reward{}).
		Where("category IS NULL OR category = ''").
		Update("category", "Voucher").Error; err != nil {
		return nil, fmt.Errorf("failed to backfill reward categories: %w", err)
	}
	return db, nil
}

