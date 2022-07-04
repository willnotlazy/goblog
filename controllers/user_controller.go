package controllers

import (
	"fmt"
	"goblog/models/article"
	"goblog/models/user"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

type UserController struct {
	BaseController
}

func NewUserController() *UserController {
	return new(UserController)
}

func (uc *UserController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_user, err := user.Get(id)

	if err != nil {
		uc.ResponseForSQLError(w, err)
	} else {
		_articles, err := article.GetByUserID(_user.GetStringID())
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器内部错误")
		} else {
			view.Render(w, view.D{
				"Articles": _articles,
			}, "articles.index", "articles._article_meta")
		}
	}
}
