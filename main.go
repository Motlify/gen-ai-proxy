package main

import (
	"context"
	"fmt"
	"gen-ai-proxy/src/api"
	"gen-ai-proxy/src/config"
	"gen-ai-proxy/src/database"
	"gen-ai-proxy/src/metrics"
	"io"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "gen-ai-proxy/docs"
	"html/template"
	"net/http"
)

// @title GenAI Proxy API
// @version 1.0
// @description Proxy server for Generative AI models.

// @license.name Apache 2.0

// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.oauth2.password BearerAuth
// @tokenUrl /api/login
// @description Enter your username and password to get a token.
// @security      BearerAuth


func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	log.Printf("Viper settings: %+v", viper.AllSettings())

	conn, err := database.Connect(&cfg)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	log.Println("Database connection successful")

	// Run database migrations
	databaseURL := fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	m, err := migrate.New(
		"file://db/migration",
		databaseURL,
	)
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("could not run migrations: %v", err)
	}
	log.Println("Database migrations applied successfully!")

	e := echo.New()

	db := database.New(conn)
	s, err := api.NewService(db, &cfg)
	if err != nil {
		log.Fatalf("could not create API service: %v", err)
	}
	api.RegisterRoutes(e, s)

	// Register Prometheus metrics collector
	collector := metrics.NewMetricsCollector(db)
	prometheus.MustRegister(collector)

	// Setup template renderer
	funcMap := template.FuncMap{
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
	}

	t := &Template{
		templates: template.Must(template.New("main").Funcs(funcMap).ParseFiles(
			"src/templates/base.html",
			"src/templates/index.html",
			"src/templates/endpoint.html",
			"src/templates/dashboard.html",
			"src/templates/ui/index.html",
		)),
	}
	log.Println("Templates parsed successfully.")
	e.Renderer = t

	// Register static file serving routes
	registerStaticRoutes(e)

	// Register UI routes
	registerUIRoutes(e)

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "pong",
		})
	})

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/swagger.yaml", func(c echo.Context) error {
		return c.File("swagger.yaml")
	})

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := e.Start(":" + cfg.ServerPort); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}

// registerStaticRoutes sets up routes for serving static files.
func registerStaticRoutes(e *echo.Echo) {
	e.Static("/js", "src/templates/js")
	e.Static("/ui/css", "src/templates/ui/css")

	// Apply custom middleware for /ui/js to ensure correct Content-Type
	e.Static("/ui/js", "src/templates/ui/js")
}

// registerUIRoutes sets up routes for the main UI pages.
func registerUIRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.GET("/dashboard.html", func(c echo.Context) error {
		return c.Render(http.StatusOK, "dashboard.html", nil)
	})

	e.GET("/ui", func(c echo.Context) error {
		return c.File("src/templates/ui/index.html")
	})
}

// Template is a custom html/template renderer for Echo framework
type Template struct {
	templates *template.Template
}

// Render renders a template document
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	log.Printf("Attempting to render template: %s", name)
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", name, err)
	}
	return err
}
