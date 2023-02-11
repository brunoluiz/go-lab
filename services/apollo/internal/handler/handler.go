package handler

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine, hs ...interface {
	Register(*gin.RouterGroup)
}) {
	g := r.Group("/")
	for _, h := range hs {
		h.Register(g)
	}
}
