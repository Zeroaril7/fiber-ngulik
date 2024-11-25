package repositories

import (
	"context"
	"fiber-ngulik/modules/user/domain"
	"fiber-ngulik/modules/user/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// Add implements domain.UserRepository.
func (r *userRepository) Add(ctx context.Context, data models.User) (result models.User, err error) {
	err = r.db.WithContext(ctx).Create(&data).Error
	return data, err
}

// Delete implements domain.UserRepository.
func (r *userRepository) Delete(ctx context.Context, username string) error {
	return r.db.WithContext(ctx).Where("username = ?", username).Delete(&models.User{}).Error
}

// Get implements domain.UserRepository.
func (r *userRepository) Get(ctx context.Context, filter models.UserFilter) (result []models.User, total int64, err error) {
	db := r.db.WithContext(ctx)
	db = buildFilterQuery(db, filter)

	if err = db.Model(&models.User{}).Count(&total).Error; err != nil {
		return
	}

	if !filter.DisablePagination {
		db = db.Offset(int(filter.GetOffset())).Limit(int(filter.GetLimit()))
	}

	if err = db.Find(&result).Error; err != nil {
		return
	}

	return
}

// GetByUsername implements domain.UserRepository.
func (r *userRepository) GetByUsername(ctx context.Context, username string) (result models.User, err error) {
	err = r.db.WithContext(ctx).Where("username = ?", username).First(&result).Error
	return
}

// Update implements domain.UserRepository.
func (r *userRepository) Update(ctx context.Context, data models.User) (result models.User, err error) {
	err = r.db.WithContext(ctx).Save(&data).Error
	return data, err
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}
