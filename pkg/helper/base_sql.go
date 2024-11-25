package helper

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepository interface {
	Add(ctx context.Context, data interface{}) (result interface{}, err error)
	Update(ctx context.Context, data interface{}) (result interface{}, err error)
}

type baseRepository struct {
	db *gorm.DB
}

func (r *baseRepository) Add(ctx context.Context, data interface{}) (result interface{}, err error) {
	err = r.db.WithContext(ctx).Create(&data).Error
	return data, err
}

func (r *baseRepository) Update(ctx context.Context, data interface{}) (result interface{}, err error) {
	err = r.db.WithContext(ctx).Save(&data).Error
	return data, err
}

func NewBaseRepository(db *gorm.DB) BaseRepository {
	return &baseRepository{db: db}
}
