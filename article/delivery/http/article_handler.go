package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	models "github.com/bxcodec/go-clean-arch/models"

	articleUcase "github.com/bxcodec/go-clean-arch/article"
	"github.com/labstack/echo"

	validator "gopkg.in/go-playground/validator.v9"
)

type ResponseError struct {
	Message string `json:"message"`
}
type HttpArticleHandler struct {
	AUsecase articleUcase.ArticleUsecase
}

func (a *HttpArticleHandler) FetchArticle(c echo.Context) error {

	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listAr, nextCursor, err := a.AUsecase.Fetch(ctx, cursor, int64(num))

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

func (a *HttpArticleHandler) GetByID(c echo.Context) error {

	idP, err := strconv.Atoi(c.Param("id"))
	id := int64(idP)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	art, err := a.AUsecase.GetByID(ctx, id)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)
}

func isRequestValid(m *models.Article) (bool, error) {

	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *HttpArticleHandler) Store(c echo.Context) error {
	var article models.Article
	err := c.Bind(&article)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&article); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	ar, err := a.AUsecase.Store(ctx, &article)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, ar)
}
func (a *HttpArticleHandler) Delete(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	id := int64(idP)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	_, err = a.AUsecase.Delete(ctx, id)

	if err != nil {

		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func getStatusCode(err error) int {

	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case models.INTERNAL_SERVER_ERROR:

		return http.StatusInternalServerError
	case models.NOT_FOUND_ERROR:
		return http.StatusNotFound
	case models.CONFLIT_ERROR:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func NewArticleHttpHandler(e *echo.Echo, us articleUcase.ArticleUsecase) {
	handler := &HttpArticleHandler{
		AUsecase: us,
	}
	e.GET("/article", handler.FetchArticle)
	e.POST("/article", handler.Store)
	e.GET("/article/:id", handler.GetByID)
	e.DELETE("/article/:id", handler.Delete)

}
