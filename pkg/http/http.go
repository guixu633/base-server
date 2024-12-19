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
	svc      *service.Service
	port     int
	certFile string
	keyFile  string
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

func NewHttpsServer(port int, svc *service.Service, certFile, keyFile string) *Server {
	server := &Server{
		Server: http.Server{
			Addr:           fmt.Sprintf(":%d", port),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   1200 * time.Second,
			IdleTimeout:    1200 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		svc:      svc,
		port:     port,
		certFile: certFile,
		keyFile:  keyFile,
	}

	r := gin.Default()

	server.RawHandler(r)

	server.Server.Handler = r
	return server
}

func (s *Server) ListenAndServe() error {
	logrus.Infof("http server listen on port: %d", s.port)
	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Errorf("listen and serve failed: %v", err)
		return err
	}
	return nil
}

func (s *Server) ListenAndServeHttps() error {
	logrus.Infof("https server listen on port: %d", s.port)
	if err := s.Server.ListenAndServeTLS(s.certFile, s.keyFile); err != nil && err != http.ErrServerClosed {
		logrus.Errorf("listen and serve failed: %v", err)
		return err
	}
	return nil
}
