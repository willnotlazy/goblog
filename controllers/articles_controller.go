package controllers

import (
	"database/sql"
	"fmt"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"html/template"
	"net/http"
	"net/url"
)

type ArticlesController struct {}

func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		tmpl, err := template.New("show.gohtml").Funcs(
			template.FuncMap{
				"RouteName2URL": route.Name2URL,
				"Int64ToString": types.Int64ToString,
			},
		).ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)

		err = tmpl.Execute(w, article)
		logger.LogError(err)
	}
}

func getArticleByID(id string) (Article, error) {
	query := "SELECT * FROM articles where id = ?"
	article := Article{}
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)

	return article, err
}

type Article struct {
	Title string
	Body  string
	ID    int64
}

type articlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func NewArticlesController() *ArticlesController {
	return new(ArticlesController)
}