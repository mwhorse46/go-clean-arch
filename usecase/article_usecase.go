package usecase

import "github.com/bxcodec/go-clean-arch/models"

type ArticleUsecase interface {
	Fetch(cursor string, num int64) ([]*models.Article, string, error)
	// Create(title, content)
}
