package handler

import (
	"net/http"

	"github.com/brunoluiz/go-lab/core/app"
	"github.com/brunoluiz/go-lab/services/apollo/gen/sqlc/lists"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type list struct {
	repo lists.Querier
}

func List(repo lists.Querier) *list {
	return &list{repo: repo}
}

func (l *list) Register(r *gin.RouterGroup) {
	r.GET("/api/v1/lists/:id", l.byUID)
	r.POST("/api/v1/lists", l.create)
}

type getListResponse struct {
	List lists.List `json:"list"`
}

func (l *list) byUID(c *gin.Context) {
	out, err := l.repo.ByUID(c.Request.Context(), c.GetString("id"))
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
			List: out,
		},
	})
}

type createListResponse struct {
	List lists.List `json:"list"`
}

func (l *list) create(c *gin.Context) {
	data := lists.CreateParams{
		UID: ksuid.New().String(),
	}
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusInternalServerError, app.Envelope{
			Status:  "error",
			Message: "something went wrong",
		})
		return
	}
	spew.Dump(data)

	out, err := l.repo.Create(c.Request.Context(), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, app.Envelope{
			Status:  "error",
			Message: "something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, app.Envelope{
		Status: "ok",
		Data: createListResponse{
			List: out,
		},
	})
}
