package app

import (
	"os"

	"go.uber.org/zap"

	"github.com/Dokhoyan/2025-11-12-test/internal/config"
	"github.com/Dokhoyan/2025-11-12-test/internal/handler"
	"github.com/Dokhoyan/2025-11-12-test/internal/logger"
	"github.com/Dokhoyan/2025-11-12-test/internal/repository"
	"github.com/Dokhoyan/2025-11-12-test/internal/server"
	"github.com/Dokhoyan/2025-11-12-test/internal/service"
)

type App struct {
	config     *config.Config
	logger     *logger.Logger
	repository repository.Repository
	service    *service.Service
	handler    *handler.Handler
	httpServer *server.Server
}

func (a *App) Logger() *logger.Logger {
	return a.logger
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log, err := logger.New(getEnv("LOG_LEVEL", "info"))
	if err != nil {
		return nil, err
	}

	repo, err := repository.NewFileRepository(cfg.DataFile)
	if err != nil {
		log.Fatal("Failed to create repository", zap.Error(err))
	}

	checker := service.NewHTTPLinkChecker(cfg.CheckTimeout)
	pdfGenerator := service.NewFPDFGenerator()

	svc := service.NewService(repo, checker, pdfGenerator, cfg.CheckTimeout)

	h := handler.NewHandler(svc)

	srv := server.NewServer(cfg.GetAddr(), h, log.Logger, func() error {
		return repo.Close()
	})

	return &App{
		config:     cfg,
		logger:     log,
		repository: repo,
		service:    svc,
		handler:    h,
		httpServer: srv,
	}, nil
}

func (a *App) Run() error {
	a.logger.Info("Starting server", zap.String("addr", a.config.GetAddr()))
	return a.httpServer.Start()
}

func (a *App) Shutdown() error {
	a.logger.Info("Shutting down application")
	if err := a.repository.Close(); err != nil {
		a.logger.Error("Error closing repository", zap.Error(err))
		return err
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
