// Package server src/server/server.go
package server

import (
	"msmc/src/config"
	"msmc/src/server/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

type Server struct {
	app    *fiber.App
	db     *gorm.DB
	config *config.Config
}

func NewServer(db *gorm.DB, cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		AppName: "msmc",
	})

	server := &Server{
		app:    app,
		db:     db,
		config: cfg,
	}

	server.setupMiddlewares()
	server.setupRoutes()

	return server
}

func (s *Server) setupMiddlewares() {
	s.app.Use(logger.New())
	s.app.Use(recover.New())
}

func (s *Server) setupRoutes() {
	api := s.app.Group("/api")

	routes.SetupArcaliveRoutes(api, s.db)
}

func (s *Server) Listen(port string) error {
	return s.app.Listen(port)
}
