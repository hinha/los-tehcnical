package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/hinha/los-technical/docs"
	loanHandler "github.com/hinha/los-technical/internal/api/handler/loan"
	"github.com/hinha/los-technical/internal/infrastructure/email"
	loanRepo "github.com/hinha/los-technical/internal/infrastructure/repository/loan"
	"github.com/hinha/los-technical/internal/usecase/loan"
)

// @title Loan Service API
// @version 1.0
// @description API for managing loans
// @termsOfService http://swagger.io/terms/
// @contact.name Martinus Dawan
// @contact.email martinuz.dawan9@gmail.com
// @BasePath /
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

	// Serve Swagger UI
	e.Static("/", "web")
	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	e.File("/swagger/doc.json", "docs/swagger.json")
	e.File("/swagger/doc.yaml", "docs/swagger.yaml")
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	log.Info("Starting server on :7002")
	if err := e.Start(":7002"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
