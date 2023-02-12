package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type task struct {
}

func Task() *task {
	return &task{}
}

func (l *task) Register(r *gin.RouterGroup) {
	r.GET("/api/v1/lists/{id}/tasks", l.get)
}

type getTaskResponse struct {
	Messsage string `json:"message"`
}

func (l *task) get(c *gin.Context) {
	c.JSON(http.StatusOK, getTaskResponse{
		Messsage: "ok",
	})
}
