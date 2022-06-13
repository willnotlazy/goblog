package article

import (
	"goblog/models"
	"goblog/models/user"
	"goblog/pkg/route"
	"goblog/pkg/types"
)

type Article struct {
	models.BaseModel
	Title string `gorm:"type:varchar(255);not null;" valid:"title"`
	Body  string `gorm:"type:longtext;not null;" valid:"body"`

	UserID uint64 `gorm:"not null;index"`
	User user.User
}

func (article Article) Link() string {
	return route.Name2URL("articles.show", "id", types.Uint64ToString(article.ID))
}

func (article Article) CreatedAtDate() string {
	return article.CreatedAt.Format("2006-01-02")
}
