package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/guixu633/base-server/module/qdrant"
)

func (s *Server) GetCryptoArticle(c *gin.Context) {
	query := c.Query("query")
	ctx := c.Request.Context()
	article, err := s.svc.GetCryptoArticle(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"article": article})
}

type UpsertCryptoArticleRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Url         string `json:"url"`
	Src         string `json:"src"`
	PublishTime string `json:"publish_time"`
}

func (s *Server) UpsertCryptoArticle(c *gin.Context) {
	request := &UpsertCryptoArticleRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article := &qdrant.CryptoArticle{
		Title:   request.Title,
		Content: request.Content,
		Url:     request.Url,
		Src:     request.Src,
	}

	ctx := c.Request.Context()
	if request.PublishTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "publish time is required"})
		return
	}
	if _, err := time.Parse(time.RFC3339, request.PublishTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid publish time"})
		return
	}
	article.PublishTime = request.PublishTime

	err := s.svc.UpsertCryptoArticle(ctx, article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
