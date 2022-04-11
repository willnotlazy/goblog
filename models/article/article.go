package article

import (
	"goblog/models"
	"goblog/pkg/route"
	"goblog/pkg/types"
)

type Article struct {
	models.BaseModel
	Title string
	Body  string
}

func (article Article) Link() string {
	return route.Name2URL("articles.show", "id", types.Uint64ToString(article.ID))
}
