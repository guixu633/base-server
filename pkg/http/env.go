package http

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (s *Server) ParseEnv(c *gin.Context) {
	key := c.Param("env")
	ctx := c.Request.Context()
	logrus.WithField("key", key).Info(ctx, "read env")
	c.JSON(http.StatusOK, gin.H{
		key: os.Getenv(key),
	})
}

func (s *Server) ParseAllEnv(c *gin.Context) {
	envVars := os.Environ()

	result := make(map[string]string)
	for _, envVar := range envVars {
		pair := strings.SplitN(envVar, "=", 2)
		key := pair[0]
		value := pair[1]
		result[key] = value
	}
	c.JSON(http.StatusOK, result)
}
