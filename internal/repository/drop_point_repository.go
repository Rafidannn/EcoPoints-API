package repository

import (
	"errors"

	"ecopoints-go-api/internal/model"

	"gorm.io/gorm"
)

type DropPointRepository interface {
	Create(dropPoint *model.DropPoint) error
	Update(dropPoint *model.DropPoint) error
	Delete(id uint64) error
	FindAll() ([]model.DropPoint, error)
	FindByID(id uint64) (*model.DropPoint, error)
}

type dropPointRepository struct {
	db *gorm.DB
}

func NewDropPointRepository(db *gorm.DB) DropPointRepository {
	return &dropPointRepository{db: db}
}

func (r *dropPointRepository) Create(dropPoint *model.DropPoint) error {
	return r.db.Create(dropPoint).Error
}

func (r *dropPointRepository) Update(dropPoint *model.DropPoint) error {
	return r.db.Save(dropPoint).Error
}

func (r *dropPointRepository) Delete(id uint64) error {
	result := r.db.Delete(&model.DropPoint{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *dropPointRepository) FindAll() ([]model.DropPoint, error) {
	var dropPoints []model.DropPoint
	err := r.db.Order("id asc").Find(&dropPoints).Error
	if err != nil {
		return nil, err
	}
	return dropPoints, nil
}

func (r *dropPointRepository) FindByID(id uint64) (*model.DropPoint, error) {
	var dropPoint model.DropPoint
	err := r.db.First(&dropPoint, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dropPoint, nil
}
