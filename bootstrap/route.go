package bootstrap

import (
	"github.com/gorilla/mux"
	"goblog/http/middlerwares"
	"goblog/pkg/route"
	"goblog/routes"
)

func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)
	router.Use(middlerwares.ForceHTMLMiddleware)

	route.SetRoute(router)


	return router
}