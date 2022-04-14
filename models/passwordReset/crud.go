package passwordReset

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
)

func GetByEmail(email string) (PasswordReset, error) {
	var _passwordReset PasswordReset
	rs := model.DB.Where("email = ?", email).First(&_passwordReset)

	if err := rs.Error; err != nil {
		logger.LogError(err)
		return _passwordReset, err
	}

	return _passwordReset, nil
}

func (_passwordReset *PasswordReset) Save() (int64, error){
	rs := model.DB.Save(&_passwordReset)

	if err := rs.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}

	return rs.RowsAffected, nil
}

func (_passwordReset *PasswordReset) FirstOrCreateByEmail() error {
	if err := model.DB.Where("email = ?", _passwordReset.Email).FirstOrCreate(&_passwordReset).Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}