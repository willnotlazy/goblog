package controllers

import (
	"goblog/pkg/view"
	"net/http"
)

type AuthController struct {

}

func NewAuthController()  *AuthController {
	return new(AuthController)
}

func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {

}
