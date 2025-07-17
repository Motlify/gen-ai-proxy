package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RegisterRoutes(e *echo.Echo, s *Service) {
	// Docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// All gen-ai-proxy routes
	apiGroup := e.Group("/api")
	apiGroup.Use(middleware.Logger())
	apiGroup.Use(JWTAuthMiddleware([]byte(s.cfg.JWTSecret)))

	// User Authentication
	e.POST("/api/register", s.Register)
	e.POST("/api/login", s.Login)

	// Api Keys
	apiGroup.POST("/api-keys", s.CreateAPIKey)
	apiGroup.GET("/api-keys", s.ListAPIKeys)
	apiGroup.DELETE("/api-keys/:id", s.DeleteAPIKey)

	// Connections
	apiGroup.POST("/connections", s.CreateConnection)
	apiGroup.GET("/connections", s.ListConnections)
	apiGroup.DELETE("/connections/:id", s.DeleteConnection)

	// Providers
	apiGroup.POST("/providers", s.CreateProvider)
	apiGroup.GET("/providers", s.ListProviders)
	apiGroup.DELETE("/providers/:id", s.DeleteProvider)

	// Models
	apiGroup.POST("/models", s.CreateModel)
	apiGroup.GET("/models", s.ListModels)
	apiGroup.DELETE("/models/:id", s.SoftDeleteModel)

	// Logs
	apiGroup.GET("/conversation_logs", s.ListConversationLogs)

	// Proxies
	apiKeyGroup := e.Group("/api")
	apiKeyGroup.Use(middleware.Logger())
	apiKeyGroup.Use(APIKeyAuthMiddleware(s.db))

	apiKeyGroup.POST("/chat", s.ProxyOllamaChat)
	apiKeyGroup.POST("/v1/chat/completions", s.ProxyOpenAIChat)

}
