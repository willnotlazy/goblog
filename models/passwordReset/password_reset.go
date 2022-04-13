package passwordReset

import (
	"errors"
	"goblog/models"
	"goblog/pkg/password"
	"time"
)

type PasswordReset struct {
	models.BaseModel
	Email string `gorm:"column:email;type:varchar(256);not null;unique;"`
	Salt string `gorm:"type:varchar(256);"`
	ExpireAt time.Time `gorm:"column:expire_at;"`
}

func (_passwordReset *PasswordReset) GenerateSalt() {
	_passwordReset.ExpireAt = time.Now()
	_passwordReset.Salt = password.Hash(_passwordReset.ExpireAt.String())
}

func (_passwordReset *PasswordReset) CanReset(salt string) error {
	if salt != _passwordReset.Salt {
		return errors.New("邮箱验证不正确")
	}

	expireAt := _passwordReset.ExpireAt
	now := time.Now()

	if now.Sub(expireAt) < 0 {
		return errors.New("邮箱验证已过期，请重新操作")
	}

	return nil
}