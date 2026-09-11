package repository

import (
	"errors"

	"ecopoints-go-api/internal/model"

	"gorm.io/gorm"
)

type WasteTypeRepository interface {
	Create(wasteType *model.WasteType) error
	Update(wasteType *model.WasteType) error
	Delete(id uint64) error
	FindAll() ([]model.WasteType, error)
	FindByID(id uint64) (*model.WasteType, error)
}

type wasteTypeRepository struct {
	db *gorm.DB
}

func NewWasteTypeRepository(db *gorm.DB) WasteTypeRepository {
	return &wasteTypeRepository{db: db}
}

func (r *wasteTypeRepository) Create(wasteType *model.WasteType) error {
	return r.db.Create(wasteType).Error
}

func (r *wasteTypeRepository) Update(wasteType *model.WasteType) error {
	return r.db.Save(wasteType).Error
}

func (r *wasteTypeRepository) Delete(id uint64) error {
	result := r.db.Delete(&model.WasteType{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *wasteTypeRepository) FindAll() ([]model.WasteType, error) {
	var wasteTypes []model.WasteType
	err := r.db.Order("id asc").Find(&wasteTypes).Error
	if err != nil {
		return nil, err
	}
	return wasteTypes, nil
}

func (r *wasteTypeRepository) FindByID(id uint64) (*model.WasteType, error) {
	var wasteType model.WasteType
	err := r.db.First(&wasteType, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &wasteType, nil
}
