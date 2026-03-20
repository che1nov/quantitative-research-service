package http

import (
	"context"
	"net/http"
	"time"

	"github.com/che1nov/quantitative-research-service/pkg/logger"
)

// Server инкапсулирует HTTP сервер приложения.
type Server struct {
	httpServer *http.Server
	log        logger.Logger
}

func New(addr string, handler http.Handler, log logger.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		},
		log: log,
	}
}

// Start запускает HTTP сервер.
func (s *Server) Start() error {
	s.log.Info("HTTP сервер запускается", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown останавливает HTTP сервер.
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.InfoContext(ctx, "Остановка HTTP сервера")
	return s.httpServer.Shutdown(ctx)
}
