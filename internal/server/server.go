package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Dokhoyan/2025-11-12-test/internal/handler"
)

// Server представляет HTTP сервер
type Server struct {
	httpServer *http.Server
	handler    *handler.Handler
	logger     *zap.Logger
	shutdownCh chan struct{}
	onShutdown func() error
}

// NewServer создает новый сервер
func NewServer(addr string, h *handler.Handler, logger *zap.Logger, onShutdown func() error) *Server {
	mux := http.NewServeMux()

	// Применяем middleware к обработчикам
	mux.HandleFunc("/links", h.LoggingMiddleware(logger, h.MethodMiddleware(http.MethodPost, h.AddLinks)))
	mux.HandleFunc("/report", h.LoggingMiddleware(logger, h.MethodMiddleware(http.MethodPost, h.GenerateReport)))
	mux.HandleFunc("/links/get", h.LoggingMiddleware(logger, h.MethodMiddleware(http.MethodGet, h.GetLinkSet)))

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		handler:    h,
		logger:     logger,
		shutdownCh: make(chan struct{}),
		onShutdown: onShutdown,
	}
}

// Start запускает сервер
func (s *Server) Start() error {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		s.logger.Info("Received shutdown signal")
		s.Shutdown()
	}()

	s.logger.Info("Server starting", zap.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	<-s.shutdownCh
	return nil
}

// Shutdown корректно останавливает сервер
func (s *Server) Shutdown() {
	s.logger.Info("Stopping accepting new connections...")

	// Останавливаем прием новых соединений, но даем время завершиться текущим запросам
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("Server forced to shutdown", zap.Error(err))
	}

	// Выполняем дополнительные действия при shutdown (например, сохранение данных)
	if s.onShutdown != nil {
		s.logger.Info("Saving data...")
		if err := s.onShutdown(); err != nil {
			s.logger.Error("Error during shutdown", zap.Error(err))
		}
	}

	s.logger.Info("Server stopped")
	close(s.shutdownCh)
}
