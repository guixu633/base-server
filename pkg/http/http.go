package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/guixu633/base-server/pkg/service"
	"github.com/sirupsen/logrus"
)

type Server struct {
	http.Server
	svc  *service.Service
	port int
}

func NewServer(port int, svc *service.Service) *Server {
	server := &Server{
		Server: http.Server{
			Addr:           fmt.Sprintf(":%d", port),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   1200 * time.Second,
			IdleTimeout:    1200 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		svc:  svc,
		port: port,
	}

	r := gin.Default()

	server.RawHandler(r)

	server.Server.Handler = r
	return server
}

func (s *Server) ListenAndServe() error {
	logrus.Infof("http server listen on port: %d", s.port)
	return s.Server.ListenAndServe()
}
