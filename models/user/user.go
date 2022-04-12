package user

import (
	"goblog/models"
	"goblog/pkg/logger"
	"goblog/pkg/model"
)

type User struct {
	models.BaseModel

	Name      string `gorm:"type:varchar(256);not null;unique" valid:"name"`
	Email     string `gorm:"type:varchar(256);unique" valid:"email"`
	Password  string `gorm:"type:varchar(256)" valid:"password"`
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

func (user *User) Create() error {
	if err := model.DB.Create(&user).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}


