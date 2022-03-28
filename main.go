package main

import (
	"fmt"
	"net/http"
	"strings"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "404 not found, path: "+r.URL.Path+", query: "+r.URL.RawQuery)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func main() {

	route := http.NewServeMux()

	route.HandleFunc("/", defaultHandler)
	route.HandleFunc("/about", aboutHandler)

	route.HandleFunc("/articles/", func(writer http.ResponseWriter, request *http.Request) {
		id := strings.SplitN(request.URL.Path, "/", 3)[2]
		fmt.Fprint(writer, "article id "+id)
	})

	http.ListenAndServe(":8088", route)
}
