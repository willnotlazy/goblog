package controllers

import (
	"bytes"
	"fmt"
	"goblog/app/requests"
	"goblog/models/passwordReset"
	"goblog/models/user"
	"goblog/pkg/mail"
	"goblog/pkg/view"
	"net/http"
)

type ForgotPasswordController struct {

}

func NewPasswordForgotController() *ForgotPasswordController {
	return new(ForgotPasswordController)
}

func (*ForgotPasswordController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "password.forgotpassword")
}

func (*ForgotPasswordController) ForgotPasswordMail(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")

	errs := requests.EmailValid(email)

	if len(errs) > 0 {
		view.RenderSimple(w, view.D{"Errors":errs, "Email":email}, "password.forgotpassword")
	} else {
		_passwordReset := passwordReset.PasswordReset{
			Email: email,
		}

		err := _passwordReset.FirstOrCreateByEmail()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器错误，请联系管理员")
		} else {
			_passwordReset.GenerateSalt()
			rowsAffected, _ := _passwordReset.Save()
			if rowsAffected > 0 {
				w.WriteHeader(http.StatusOK)
				body := new(bytes.Buffer)
				view.RenderSimple(body, view.D{"Salt": _passwordReset.Salt, "Email": email}, "mail.reset_password_mail")
				err = mail.SendMail("goblog 密码重置", body.String(), _passwordReset.Email)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, "服务器内部错误，请联系管理员")
				} else {
					fmt.Fprint(w,"已发送重置密码链接至您的邮箱")
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "服务器错误，请联系管理员")
			}
		}

	}
}

func (*ForgotPasswordController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	var (
		salt = queries.Get("salt")
		email = queries.Get("email")
	)

	if salt == "" || email == ""{
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "无操作权限")
	} else {
		_passwordReset, err := passwordReset.GetByEmail(email)
		if err != nil{
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "无操作权限")
		} else {
			if err := _passwordReset.CanReset(salt); err != nil {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprint(w, "无操作权限," + err.Error())
			} else {
				view.RenderSimple(w, view.D{"Email": email, "Salt": salt}, "password.reset")
			}
		}
	}
}

func (*ForgotPasswordController) DoReset(w http.ResponseWriter, r *http.Request) {
	var _passwordReset = passwordReset.PasswordReset{
		Email: r.PostFormValue("email"),
		Salt: r.PostFormValue("salt"),
		Password: r.PostFormValue("password"),
	}

	errs := requests.ValidPasswordResetForm(_passwordReset)
	if len(errs) > 0 {
		_, saltOk := errs["salt"]
		_, emailOk := errs["email"]

		if saltOk || emailOk {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "无操作权限")
		} else {
			w.WriteHeader(http.StatusFound)
			view.RenderSimple(w, view.D{"Salt": _passwordReset.Salt, "Email": _passwordReset.Email, "Password": _passwordReset.Password, "Error": errs}, "password.reset")
		}
	} else {
		_passwordReset_db, err := passwordReset.GetByEmail(_passwordReset.Email)
		fmt.Println(_passwordReset.Email)
		_user, uerr := user.GetByEmail(_passwordReset.Email)
		if err != nil || uerr != nil || _passwordReset_db.CanReset(_passwordReset.Salt) != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "无权限操作, 改账户不存在或本次操作已过期")
		} else {
			_user.Password = _passwordReset.Password
			err := _user.Save()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "修改失败，请联系管理员")
			} else {
				w.WriteHeader(http.StatusFound)
				view.RenderSimple(w, view.D{}, "auth.login")
			}
		}
	}
}
