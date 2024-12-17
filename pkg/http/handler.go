package http

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) RawHandler(r *gin.Engine) {
	r.GET("/env/all", s.ParseAllEnv)
	r.GET("/env/:env", s.ParseEnv)
}
