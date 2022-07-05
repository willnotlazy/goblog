package controllers

import (
	"fmt"
	"goblog/app/policies"
	"goblog/app/requests"
	"goblog/models/article"
	"goblog/pkg/auth"
	"goblog/pkg/config"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"goblog/pkg/view"
	"net/http"
)

type ArticlesController struct{
	BaseController
}

func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {

	perpage := config.GetInt("pagination.perpage")
	_articles, viewData, err := article.GetAll(r, perpage)
	logger.LogError(err)

	view.Render(w, view.D{"Articles": _articles, "Pagination": viewData}, "articles.index", "articles._article_meta")
}

func (ac *ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		view.Render(w, view.D{"Article": _article, "CanModifyArticle": policies.CanModifyArticle(_article)}, "articles.show", "articles._article_meta")
	}
}

func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	_article := article.Article{
		Title: r.PostFormValue("title"),
		Body: r.PostFormValue("body"),
		UserID: types.StringToUint64(auth.User().GetStringID()),
	}

	errors := requests.ValidateArticleForm(_article)

	// 检查是否有错误
	if len(errors) == 0 {
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功，ID为"+types.Uint64ToString(_article.ID))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}
	} else {

		data := view.D{
			"Article": _article,
			"Errors": errors,
		}
		view.Render(w, data, []string{"articles.create", "articles._form_field"}...)
	}
}

func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, []string{"articles.create", "articles._form_field"}...)
}

func (ac *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			data := view.D{
				"Article": _article,
				"Errors":  nil,
			}
			view.Render(w, data, "articles.edit", "articles._form_field")
		}
	}
}

func (ac *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			_article.Title = r.PostFormValue("title")
			_article.Body = r.PostFormValue("body")

			errors := requests.ValidateArticleForm(_article)

			if len(errors) == 0 {
				rowsAffected, err := _article.Update()

				if err != nil {
					logger.LogError(err)
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, "500 服务器内部错误")
				}

				if rowsAffected > 0 {
					showURL := route.Name2URL("articles.show", "id", id)

					http.Redirect(w, r, showURL, http.StatusFound)
				} else {
					fmt.Fprint(w, "你没有做任何修改")
				}
			} else {

				data := view.D{
					"Article": _article,
					"Errors":  errors,
				}
				view.Render(w, data, "articles.edit", "articles._form_field")
			}
		}
	}
}

func (ac *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			rowsAffected, err := _article.Delete()

			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			} else {
				if rowsAffected > 0 {
					indexURL := route.Name2URL("articles.index")

					http.Redirect(w, r, indexURL, http.StatusFound)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprint(w, "404 文章未找到")
				}
			}
		}
	}
}

type ArticlesFormData struct {
	Title, Body string
	Article     article.Article
	Errors      map[string]string
}

func NewArticlesController() *ArticlesController {
	return &ArticlesController{}
}
