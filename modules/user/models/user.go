package models

type User struct {
	ID       int64  `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (User) TableName() string {
	return "user"
}
