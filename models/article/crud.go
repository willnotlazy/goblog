package article

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/pagination"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"net/http"
)

func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToUint64(idstr)
	if err := model.DB.Preload("User").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

func GetAll(r *http.Request, perpage int) ([]Article, pagination.ViewData, error) {
	var articles []Article
	db := model.DB.Model(Article{}).Order("created_at desc")
	_page := pagination.New(r, db, route.Name2URL("articles.index"), perpage)
	viewData := _page.Paging()
	_page.Results(&articles)
	return articles, viewData, nil
}

func (article *Article) Create() error {
	rs := model.DB.Create(&article)

	if err := rs.Error; err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

func (article *Article) Update() (int64, error) {
	rs := model.DB.Save(&article)
	if err := rs.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}

	return rs.RowsAffected, nil
}

func (article *Article) Delete() (int64, error) {
	rs := model.DB.Delete(&article)

	if err := rs.Error; err != nil {
		logger.LogError(err)
		return 0, err
	}

	return rs.RowsAffected, nil
}

func GetByUserID(userId string) ([]Article, error) {
	var articles []Article

	if err := model.DB.Preload("User").Where("user_id = ?", userId).Find(&articles).Error; err != nil {
		logger.LogError(err)
		return articles, err
	}

	return articles, nil

}