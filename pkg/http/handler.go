package http

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) RawHandler(r *gin.Engine) {
	r.GET("/env/all", s.ParseAllEnv)
	r.GET("/env/:env", s.ParseEnv)
	r.GET("/config", s.GetConfig)

	r.GET("/oss/tree/*path", s.OssGetTree)
	r.GET("/oss/blob/*path", s.OssGetBlob)
}
