package middlerwares

import (
	"goblog/pkg/session"
	"net/http"
)

func StartSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session.StartSession(writer, request)

		next.ServeHTTP(writer, request)
	})
}
