package user

import (
	"goblog/models"
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/password"
	"goblog/pkg/route"
)

type User struct {
	models.BaseModel

	Name            string `gorm:"type:varchar(256);not null;unique" valid:"name"`
	Email           string `gorm:"type:varchar(256);unique" valid:"email"`
	Password        string `gorm:"type:varchar(256)" valid:"password"`
	PasswordConfirm string `gorm:"-" valid:"password_confirm"`
}

func (user *User) Create() error {
	if err := model.DB.Create(&user).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

func (user *User) ComparePassword(_password string) bool {
	return password.CheckHash(_password, user.Password)
}

func (user User) Link() string {
	return route.Name2URL("users.show", "id", user.GetStringID())
}
