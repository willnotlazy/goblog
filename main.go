package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"goblog/bootstrap"
	"goblog/config"
	"goblog/http/middlerwares"
	app_config "goblog/pkg/config"
	"goblog/pkg/database"
	"net/http"
	"net/url"
)

var router *mux.Router
var db *sql.DB

type articlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func init() {
	config.Initialize()
}

func main() {

	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	http.ListenAndServe(":"+app_config.GetString("app.port"), middlerwares.RemoveTrailingSlash(router))
}
