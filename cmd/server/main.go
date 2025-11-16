package main

import (
	"github.com/Dokhoyan/2025-11-12-test/internal/app"

	"go.uber.org/zap"
)

func main() {
	application, err := app.New()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := application.Shutdown(); err != nil {
			application.Logger().Error("Error during shutdown", zap.Error(err))
		}
	}()

	if err := application.Run(); err != nil {
		application.Logger().Fatal("Application error", zap.Error(err))
	}
}
