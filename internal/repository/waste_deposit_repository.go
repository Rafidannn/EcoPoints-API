package repository

import (
	"errors"
	"fmt"
	"time"

	"ecopoints-go-api/internal/model"

	"gorm.io/gorm"
)

type WasteDepositRepository interface {
	Create(deposit *model.WasteDeposit) error
	GetByID(id uint64) (*model.WasteDeposit, error)
	GetByUserID(userID uint64) ([]model.WasteDeposit, error)
	GetAll(status string) ([]model.WasteDeposit, error)
	Verify(depositID uint64, verifierID uint64, actualWeight float64, notes *string) (*model.WasteDeposit, uint, error)
	Reject(depositID uint64, actorID uint64, notes *string) (*model.WasteDeposit, error)
	Cancel(depositID uint64, userID uint64, notes *string) (*model.WasteDeposit, error)
}

type wasteDepositRepository struct {
	db *gorm.DB
}

func NewWasteDepositRepository(db *gorm.DB) WasteDepositRepository {
	return &wasteDepositRepository{db: db}
}

func (r *wasteDepositRepository) Create(deposit *model.WasteDeposit) error {
	return r.db.Create(deposit).Error
}

func (r *wasteDepositRepository) GetByID(id uint64) (*model.WasteDeposit, error) {
	var deposit model.WasteDeposit
	err := r.db.
		Preload("User").
		Preload("WasteType").
		Preload("DropPoint").
		Preload("Verifier").
		First(&deposit, id).Error
	if err != nil {
		return nil, err
	}
	return &deposit, nil
}

func (r *wasteDepositRepository) GetByUserID(userID uint64) ([]model.WasteDeposit, error) {
	var deposits []model.WasteDeposit
	err := r.db.
		Preload("WasteType").
		Preload("DropPoint").
		Preload("Verifier").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&deposits).Error
	return deposits, err
}

func (r *wasteDepositRepository) GetAll(status string) ([]model.WasteDeposit, error) {
	var deposits []model.WasteDeposit
	query := r.db.
		Preload("User").
		Preload("WasteType").
		Preload("DropPoint").
		Preload("Verifier").
		Order("created_at DESC")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&deposits).Error
	return deposits, err
}

// Verify verifies a deposit, computes earned points, updates the user balance, and inserts a point_transaction
func (r *wasteDepositRepository) Verify(depositID uint64, verifierID uint64, actualWeight float64, notes *string) (*model.WasteDeposit, uint, error) {
	var deposit model.WasteDeposit
	var earnedPoints uint

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch deposit with lock
		if err := tx.Preload("WasteType").Preload("User").First(&deposit, depositID).Error; err != nil {
			return fmt.Errorf("setoran tidak ditemukan: %w", err)
		}

		if deposit.Status == "verified" {
			return errors.New("setoran ini sudah diverifikasi sebelumnya")
		}

		// 2. Compute points based on weight and waste type
		if actualWeight <= 0 {
			actualWeight = deposit.WeightKg
		}
		pointsPerKg := uint(0)
		if deposit.WasteType != nil {
			pointsPerKg = deposit.WasteType.PointsPerKg
		}
		earnedPoints = uint(actualWeight * float64(pointsPerKg))

		now := time.Now()
		deposit.WeightKg = actualWeight
		deposit.Status = "verified"
		deposit.VerifiedBy = &verifierID
		deposit.VerifiedAt = &now
		if notes != nil && *notes != "" {
			deposit.Notes = notes
		}

		// 3. Update deposit status
		if err := tx.Save(&deposit).Error; err != nil {
			return fmt.Errorf("gagal memperbarui setoran: %w", err)
		}

		// 4. Update user points balance and exp
		if err := tx.Model(&model.User{}).
			Where("id = ?", deposit.UserID).
			Updates(map[string]interface{}{
				"points_balance": gorm.Expr("points_balance + ?", earnedPoints),
				"exp":            gorm.Expr("exp + ?", earnedPoints),
			}).Error; err != nil {
			return fmt.Errorf("gagal menambah poin user: %w", err)
		}

		// 5. Create point_transactions record
		desc := fmt.Sprintf("Setor %s (%.1f kg)", deposit.WasteType.Name, actualWeight)
		refType := "waste_deposit"
		refID := deposit.ID
		pt := model.PointTransaction{
			UserID:        deposit.UserID,
			Type:          "credit",
			Amount:        earnedPoints,
			ReferenceType: &refType,
			ReferenceID:   &refID,
			Description:   &desc,
			CreatedAt:     &now,
			UpdatedAt:     &now,
		}
		if err := tx.Create(&pt).Error; err != nil {
			return fmt.Errorf("gagal mencatat transaksi poin: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	// Fetch refreshed record
	refreshed, _ := r.GetByID(depositID)
	if refreshed != nil {
		return refreshed, earnedPoints, nil
	}
	return &deposit, earnedPoints, nil
}

// Reject rejects a pending deposit (staff/admin only). No points are awarded.
func (r *wasteDepositRepository) Reject(depositID uint64, actorID uint64, notes *string) (*model.WasteDeposit, error) {
	var deposit model.WasteDeposit

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&deposit, depositID).Error; err != nil {
			return fmt.Errorf("setoran tidak ditemukan: %w", err)
		}

		if deposit.Status != "pending" {
			return errors.New("hanya setoran berstatus pending yang bisa ditolak")
		}

		now := time.Now()
		deposit.Status = "rejected"
		deposit.VerifiedBy = &actorID
		deposit.VerifiedAt = &now
		if notes != nil && *notes != "" {
			deposit.Notes = notes
		}

		if err := tx.Save(&deposit).Error; err != nil {
			return fmt.Errorf("gagal memperbarui setoran: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Fetch refreshed record with preloads
	refreshed, _ := r.GetByID(depositID)
	if refreshed != nil {
		return refreshed, nil
	}
	return &deposit, nil
}

// Cancel cancels a pending deposit by its owner (nasabah).
func (r *wasteDepositRepository) Cancel(depositID uint64, userID uint64, notes *string) (*model.WasteDeposit, error) {
	var deposit model.WasteDeposit

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&deposit, depositID).Error; err != nil {
			return fmt.Errorf("setoran tidak ditemukan: %w", err)
		}

		if deposit.UserID != userID {
			return errors.New("Anda tidak berhak membatalkan setoran ini")
		}

		if deposit.Status != "pending" {
			return errors.New("hanya setoran berstatus pending yang bisa dibatalkan")
		}

		deposit.Status = "cancelled"
		if notes != nil && *notes != "" {
			deposit.Notes = notes
		}

		if err := tx.Save(&deposit).Error; err != nil {
			return fmt.Errorf("gagal memperbarui setoran: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	refreshed, _ := r.GetByID(depositID)
	if refreshed != nil {
		return refreshed, nil
	}
	return &deposit, nil
}
