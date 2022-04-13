package requests

import (
	"github.com/thedevsaddam/govalidator"
)

type EmailValidator struct {
	Email string `valid:"email"`
}

func EmailValid(email string) map[string][]string {
	var rules = govalidator.MapData{
		"email": []string{"required", "email", "exists:users,email"},
	}

	var messages = govalidator.MapData{
		"email": []string{"required:邮箱必填", "email:输入必须为正确的邮箱格式"},
	}

	opts := govalidator.Options{
		Data: &EmailValidator{Email: email},
		Rules: rules,
		Messages: messages,
		TagIdentifier: "valid",
	}

	errs := govalidator.New(opts).ValidateStruct()

	return errs
}
