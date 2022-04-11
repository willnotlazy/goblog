package routes

import (
	"github.com/gorilla/mux"
	"goblog/controllers"
	"net/http"
)

func RegisterWebRoutes(r *mux.Router) {
	pc := controllers.NewPageController()

	// 前段资源
	r.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))

	r.HandleFunc("/", pc.Home).Methods("GET").Name("home")
	r.HandleFunc("/about", pc.About).Methods("GET").Name("about")
	// 404
	r.NotFoundHandler = http.HandlerFunc(pc.NotFound)

	ac := controllers.NewArticlesController()
	r.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
	r.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
	r.HandleFunc("/articles", ac.Store).Methods("POST").Name("articles.store")
	r.HandleFunc("/articles/create", ac.Create).Methods("GET").Name("articles.create")
	r.HandleFunc("/articles/{id:[0-9]+}/edit", ac.Edit).Methods("GET").Name("articles.edit")
	r.HandleFunc("/articles/{id:[0-9+]}", ac.Update).Methods("POST").Name("articles.update")
	r.HandleFunc("/articles/{id:[0-9+]}/delete", ac.Delete).Methods("POST").Name("articles.delete")


	//r.Use(middlerwares.ForceHTML)
}
