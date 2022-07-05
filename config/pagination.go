package config

import (
	"goblog/pkg/config"
)

func init() {
	config.Add("pagination", config.StrMap{
		"url_query" : config.Env("PAGINATION_URL_QUERY", "page"),
		"perpage" : config.Env("PAGINATION_PERPAGE", 5),
	})
}