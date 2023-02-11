package handler

import (
	"context"
	"net/http"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/services/apollo"
	"github.com/gin-gonic/gin"
)

type ListRepository interface {
	ByID(context.Context, string) (*apollo.List, error)
}

type list struct {
	repo ListRepository
}

func List(repo ListRepository) *list {
	return &list{repo: repo}
}

func (l *list) Register(r *gin.RouterGroup) {
	r.GET("/api/v1/lists", l.get)
}

type getListResponse struct {
	Lists []*apollo.List `json:"lists"`
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
			Lists: []*apollo.List{data},
		},
	})
}
