package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	loanHandler "github.com/hinha/los-technical/internal/api/handler/loan"
	"github.com/hinha/los-technical/internal/infrastructure/email"
	loanRepo "github.com/hinha/los-technical/internal/infrastructure/repository/loan"
	"github.com/hinha/los-technical/internal/usecase/loan"
)

func main() {
	// Initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	// Create repository
	repository := loanRepo.NewInMemoryRepository(log)
	emailSender := email.NewConsoleEmailSender(log)
	loanService := loan.NewLoanService(repository, emailSender, log)
	handler := loanHandler.NewHandler(loanService)

	// Initialize Echo
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register routes
	handler.RegisterRoutes(e)

	// Start server
	log.Info("Starting server on :8080")
	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
