package route

import (
	"github.com/gorilla/mux"
	"goblog/pkg/logger"
	"net/http"
)

func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)

	return vars[parameterName]
}

func Name2URL(routeName string, pair ...string) string {
	var router *mux.Router
	url, err := router.Get(routeName).URL(pair...)

	if err != nil {
		logger.LogError(err)
		return ""
	}

	return url.String()
}
