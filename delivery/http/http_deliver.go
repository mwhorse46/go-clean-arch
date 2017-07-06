package http

import (
	artHandler "github.com/bxcodec/go-clean-arch/delivery/http/article"
	artUcase "github.com/bxcodec/go-clean-arch/usecase"
	"github.com/labstack/echo"
)

func Init(e *echo.Echo, au artUcase.ArticleUsecase) {
	articleHandler := &artHandler.ArticleHandler{AUsecase: au}
	e.GET(`/article`, articleHandler.FetchArticle)
	e.GET(`/article/:id`, articleHandler.GetByID)
}
