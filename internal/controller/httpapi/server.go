// Package httpapi contains required primitives to start Server
package httpapi

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/osapers/mch-back/internal/service/event"
	"github.com/osapers/mch-back/internal/service/tag"
	"github.com/osapers/mch-back/internal/service/user"

	"github.com/gin-gonic/gin"
	corsWrapper "github.com/rs/cors/wrapper/gin"

	"go.uber.org/zap"
)

// Server is model for handling http requests
type Server struct {
	srv    *http.Server
	es     *event.Service
	us     *user.Service
	ts     *tag.Service
	cfg    *config
	logger *zap.Logger
}

// NewServer returns instance of real Server
func NewServer(es *event.Service, us *user.Service, ts *tag.Service) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	s := &Server{
		cfg:    newConfig(),
		logger: newLogger(),
		es:     es,
		us:     us,
		ts:     ts,
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.Use(gin.Recovery())
	r.Use(loggerMiddleware(s.logger))
	r.Use(corsWrapper.AllowAll())

	s.addHandlers(r)

	s.srv = &http.Server{
		Addr:    net.JoinHostPort(s.cfg.host, s.cfg.port),
		Handler: r,
	}

	return s
}

// Start the Server
func (s *Server) Start() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("listen and serve", zap.Error(err))
		}
	}()

	s.logger.Info("server started", zap.String("url", fmt.Sprintf("http://%s", s.srv.Addr)))

	<-done
}

// Shutdown the Server
func (s *Server) Shutdown() error {
	s.logger.Info("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		_ = s.logger.Sync()
		cancel()
	}()

	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("Server shutdown failed - %s", err)
	}

	s.logger.Info("server exited properly")

	return nil
}
