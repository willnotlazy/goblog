package user

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/types"
)

func Get(idstr string) (User, error) {
	var _user User

	if err := model.DB.First(&_user, types.StringToUint64(idstr)).Error; err != nil {
		logger.LogError(err)
		return _user, err
	}

	return _user, nil
}

func GetByEmail(email string) (User, error) {
	var _user User

	if err := model.DB.First(&_user, "email = ?", email).Error; err != nil {
		logger.LogError(err)
		return _user, err
	}

	return _user, nil
}
