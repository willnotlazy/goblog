package routes

import (
	"github.com/gorilla/mux"
	"goblog/controllers"
	"net/http"
)

func RegisterWebRoutes(r *mux.Router) {
	pc := controllers.NewPageController()
	r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	// 404
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)

	ac := controllers.NewArticlesController()
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.Index")
}
