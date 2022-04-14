package route

import (
	"github.com/gorilla/mux"
	"goblog/pkg/logger"
	"html/template"
	"net/http"
)

var route *mux.Router

func SetRoute(r *mux.Router) {
	route = r
}

func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)

	return vars[parameterName]
}

func Name2URL(routeName string, pair ...string) string {
	url, err := route.Get(routeName).URL(pair...)
	if err != nil {
		logger.LogError(err)
		return ""
	}

	return url.String()
}

func Name2FullURL(routeName string, r *http.Request, pair ...string) string {
	url := Name2URL(routeName, pair...)

	host := r.Host
	return host + url
}

func URL2String(url template.URL) string {
	return string(url)
}

func SafeURL(url string) template.URL {
	return template.URL(url)
}
