package article

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/types"
)

func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToUint64(idstr)
	if err := model.DB.Preload("User").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}

func GetAll() ([]Article, error) {
	var articles []Article

	if err := model.DB.Preload("User").Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
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
