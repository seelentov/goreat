package http_rest

import (
	"goreat/internal/controllers/http/http_utils"
	"goreat/internal/models/queries"
	"goreat/internal/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EntityRestController struct {
	entityRepo *repo.EntityRepository
}

func NewEntityRestController(entityRepo *repo.EntityRepository) *EntityRestController {
	return &EntityRestController{
		entityRepo: entityRepo,
	}
}

type PostGetDataReq queries.Query

func (c *EntityRestController) PostGetData(ctx *gin.Context) {
	var q PostGetDataReq
	if errs := http_utils.ShouldBindJSON(&q, ctx); errs != nil {
		ctx.JSON(http.StatusBadRequest, errs)
		return
	}

	res := c.entityRepo.ByQuery(queries.Query(q))
	if res.Error != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, res.Entities)
}
