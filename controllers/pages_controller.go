package controllers

import (
	"fmt"
	"goblog/pkg/view"
	"net/http"
)

type PageController struct{}

func (page *PageController) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

func (page *PageController) About(w http.ResponseWriter, r *http.Request) {
	view.Render(w, nil, "pages.about")
}

func (page *PageController) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑,请联系我们.</p>")
}

func NewPageController() *PageController {
	return &PageController{}
}
