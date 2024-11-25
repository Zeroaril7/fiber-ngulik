package repositories

import (
	"fiber-ngulik/modules/user/models"

	"gorm.io/gorm"
)

func buildFilterQuery(db *gorm.DB, f models.UserFilter) *gorm.DB {
	if f.Role != "" {
		db = db.Where("role = ?", f.Role)
	}

	return db
}
