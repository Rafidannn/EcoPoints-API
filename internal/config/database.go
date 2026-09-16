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

		// Drop foreign keys and indexes first so MySQL allows dropping waste_type_id column
		log.Println("Dropping legacy foreign keys and columns from waste_deposits...")
		_ = db.Exec("ALTER TABLE `waste_deposits` DROP FOREIGN KEY `fk_waste_deposits_waste_type`").Error
		_ = db.Exec("ALTER TABLE `waste_deposits` DROP FOREIGN KEY `waste_deposits_waste_type_id_foreign`").Error
		_ = db.Exec("ALTER TABLE `waste_deposits` DROP INDEX `fk_waste_deposits_waste_type`").Error
		_ = db.Exec("ALTER TABLE `waste_deposits` DROP INDEX `waste_deposits_waste_type_id_foreign`").Error

		_ = db.Exec("ALTER TABLE `waste_deposits` DROP COLUMN `waste_type_id`").Error
		_ = db.Exec("ALTER TABLE `waste_deposits` DROP COLUMN `weight_kg`").Error
		_ = db.Exec("ALTER TABLE `waste_deposits` DROP COLUMN `original_weight_kg`").Error
		_ = db.Exec("ALTER TABLE `waste_deposits` DROP COLUMN `actual_weight_kg`").Error

		// If column still exists for any reason, make it nullable
		if db.Migrator().HasColumn("waste_deposits", "waste_type_id") {
			_ = db.Exec("ALTER TABLE `waste_deposits` MODIFY COLUMN `waste_type_id` bigint unsigned NULL DEFAULT NULL").Error
		}
	}

	// Drop old triggers on waste_deposits if any
	type TriggerRow struct {
		TriggerName string `gorm:"column:TRIGGER_NAME"`
	}
	var triggers []TriggerRow
	_ = db.Raw("SELECT TRIGGER_NAME FROM information_schema.TRIGGERS WHERE EVENT_OBJECT_SCHEMA = ? AND EVENT_OBJECT_TABLE = 'waste_deposits'", cfg.DBDatabase).Scan(&triggers).Error
	for _, t := range triggers {
		log.Printf("Dropping legacy trigger: %s", t.TriggerName)
		_ = db.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS `%s`", t.TriggerName)).Error
	}

	if err := db.Model(&model.Reward{}).
		Where("category IS NULL OR category = ''").
		Update("category", "Voucher").Error; err != nil {
		return nil, fmt.Errorf("failed to backfill reward categories: %w", err)
	}
	return db, nil
}

