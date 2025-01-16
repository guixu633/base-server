package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Coin(c *gin.Context) {
	coinID := c.Param("coin_id")
	coin, err := s.svc.GetCoinData(c.Request.Context(), coinID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, coin)
}
