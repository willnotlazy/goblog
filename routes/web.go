package routes

import (
	"github.com/gorilla/mux"
	"goblog/controllers"
	"net/http"
)

func RegisterWebRoutes(r *mux.Router) {
	pageController := controllers.NewPageController()
	r.HandleFunc("/", pageController.Home).Methods("GET").Name("home")
	r.HandleFunc("/about", pageController.About).Methods("GET").Name("about")
	// 404
	r.NotFoundHandler = http.HandlerFunc(pageController.NotFound)

	articlesController := controllers.NewArticlesController()
	r.HandleFunc("/articles/{id:[0-9]+}", articlesController.Show).Methods("GET").Name("articles.show")
}
