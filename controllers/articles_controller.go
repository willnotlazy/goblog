package controllers

import (
	"fmt"
	"goblog/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"unicode/utf8"
)

type ArticlesController struct {}

func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {

	_articles, err := article.GetAll()
	logger.LogError(err)

	tmpl, err := template.ParseFiles("resources/views/articles/index.gohtml")
	logger.LogError(err)

	err = tmpl.Execute(w, _articles)
	logger.LogError(err)
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)

	_article, err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
				"Uint64ToString": types.Uint64ToString,
			},
		).ParseFiles("resources/views/articles/show.gohtml")
		logger.LogError(err)

		err = tmpl.Execute(w, _article)
		logger.LogError(err)
	}
}

func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	// 检查是否有错误
	if len(errors) == 0 {
		var _article = article.Article{
			Title: title,
			Body: body,
		}
		_article.Create()
		if _article.ID > 0 {
			fmt.Fprint(w, "插入成功，ID为"+types.Uint64ToString(_article.ID))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}
	} else {
		storeURL := route.Name2URL("articles.store")

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}

		tml, err := template.ParseFiles("resources/views/articles/create.gohtml")

		if err != nil {
			panic(err)
		}

		err = tml.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	storeURL := route.Name2URL("articles.store")
	data := ArticlesFormData{
		URL: storeURL,
	}

	tml, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tml.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

type ArticlesFormData struct {
	Title, Body string
	URL         string
	Errors      map[string]string
}

func NewArticlesController() *ArticlesController {
	return &ArticlesController{}
}

func validateArticleFormData(title, body string) map[string]string {
	errors := make(map[string]string)

	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}