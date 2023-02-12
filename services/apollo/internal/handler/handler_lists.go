package handler

import (
	"context"
	"net/http"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/services/apollo/gen/sqlc/lists"
	"github.com/gin-gonic/gin"
)

type list struct {
	repo lists.Querier
}

func List(repo lists.Querier) *list {
	return &list{repo: repo}
}

func (l *list) Register(r *gin.RouterGroup) {
	r.GET("/api/v1/lists", l.get)
}

type getListResponse struct {
	Lists []lists.List `json:"lists"`
}

func (l *list) get(c *gin.Context) {
	data, err := l.repo.ByID(context.Background(), "something")
	if err != nil {
		c.JSON(http.StatusInternalServerError, app.Envelope{
			Status:  "error",
			Message: "something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, app.Envelope{
		Status: "ok",
		Data: getListResponse{
			Lists: []lists.List{data},
		},
	})
}
