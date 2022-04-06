package route

import (
	"github.com/gorilla/mux"
	"net/http"
)

var Router *mux.Router

func Initialize() {
	Router = mux.NewRouter().StrictSlash(true)
}

func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)

	return vars[parameterName]
}

func Name2URL(routeName string, pair ...string) string {
	url, err := Router.Get(routeName).URL(pair...)

	if err != nil {
		//checkError(err)
		return ""
	}

	return url.String()
}
