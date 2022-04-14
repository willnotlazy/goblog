package requests

import (
	"github.com/thedevsaddam/govalidator"
	"goblog/models/passwordReset"
)

func ValidPasswordResetForm(data passwordReset.PasswordReset) map[string][]string {
	rules := govalidator.MapData{
		"password": []string{"required", "min:6"},
		"email": []string{"required", "exists:password_resets,email"},
		"salt": []string{"required"},
	}

	messages := govalidator.MapData{
		"password": []string{
			"required:密码为必填项",
			"min:长度需大于 6",
		},
		"email": []string{
			"required:Email 为必填项",
		},
		"salt": []string{
			"required:salt 为必填项",
		},
	}

	opts := govalidator.Options{
		Data: &data,
		Rules: rules,
		Messages: messages,
		TagIdentifier: "valid",
	}

	return govalidator.New(opts).ValidateStruct()
}
