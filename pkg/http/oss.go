package http

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (s *Server) OssGetTree(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "path参数不能为空",
		})
		return
	}

	exists, err := s.svc.Exists(c, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "检查路径是否存在失败",
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "路径不存在",
		})
		return
	}

	files, err := s.svc.GetDir(c, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取目录失败",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": files,
	})
}

func (s *Server) OssGetBlob(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "path参数不能为空",
		})
		return
	}

	// 获取文件内容
	data, err := s.svc.GetFile(c, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取文件失败",
		})
		return
	}

	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "文件不存在",
		})
		return
	}

	// 获取文件名
	filename := filepath.Base(path)

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))

	// 写入响应体
	c.Data(http.StatusOK, "application/octet-stream", data)
}
