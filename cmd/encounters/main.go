package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/emiliopalmerini/due-draghi-combattimenti/internal/application/encounter"
	monsterApp "github.com/emiliopalmerini/due-draghi-combattimenti/internal/application/monster"
	"github.com/emiliopalmerini/due-draghi-combattimenti/internal/infrastructure/config"
	"github.com/emiliopalmerini/due-draghi-combattimenti/internal/infrastructure/persistence/memory"
	"github.com/emiliopalmerini/due-draghi-combattimenti/internal/infrastructure/web/handlers"
	"github.com/emiliopalmerini/due-draghi-combattimenti/internal/infrastructure/web/templates"
)

type ctxKey string

// App represents the main application with all dependencies
type App struct {
	router           chi.Router
	config           *config.Config
	logger           *slog.Logger
	encounterService *encounter.Service
	encounterHandler *handlers.EncounterHandler
	monsterHandler   *handlers.MonsterHandler
	queryHandler     *encounter.QueryHandler
}

// NewApp creates a new application instance with all dependencies
func NewApp(cfg *config.Config, logger *slog.Logger) (*App, error) {
	// Initialize repositories
	repo := memory.NewEncounterRepository()
	monsterRepo := memory.NewMonsterRepository()

	// Initialize application services
	encounterService := encounter.NewService(logger, repo)
	queryHandler := encounter.NewQueryHandler(logger, repo)
	monsterService := monsterApp.NewService(monsterRepo)

	// Initialize HTTP handlers
	encounterHandler := handlers.NewEncounterHandler(encounterService, queryHandler, monsterService, logger)
	monsterHandler := handlers.NewMonsterHandler(monsterService, logger)

	app := &App{
		config:           cfg,
		logger:           logger,
		encounterService: encounterService,
		encounterHandler: encounterHandler,
		monsterHandler:   monsterHandler,
		queryHandler:     queryHandler,
	}

	app.setupRouter()

	return app, nil
}

// setupRouter configures the Chi router with middleware and routes
func (app *App) setupRouter() {
	r := chi.NewRouter()

	// Core middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(app.loggingMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// HTMX-specific middleware
	r.Use(app.htmxMiddleware)

	// CORS middleware for HTMX requests
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"HX-*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Compress responses (good for HTMX responses)
	if app.config.IsProduction() {
		r.Use(middleware.Compress(5))
	}

	// Health check endpoints
	r.Get("/health", app.healthHandler)
	r.Get("/ready", app.readinessHandler)

	// Serve static files with no-cache headers
	r.Handle("/static/*", noCacheMiddleware(http.StripPrefix("/static/", http.FileServer(http.Dir("internal/infrastructure/static")))))

	// Application routes
	r.Route("/", func(r chi.Router) {
		r.Get("/", app.indexHandler)
		r.Post("/calculate", app.encounterHandler.CalculateHandler)
		r.Get("/party-input", app.encounterHandler.PartyInputHandler)
		r.Get("/api/difficulties", app.encounterHandler.GetDifficultiesHandler)
		r.Get("/api/monsters", app.monsterHandler.SearchHandler)
	})

	app.router = r
}

// indexHandler renders the encounters calculator page
func (app *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Render the encounters home template
	component := templates.Home()
	if err := component.Render(r.Context(), w); err != nil {
		app.logger.Error("Failed to render encounters home template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// noCacheMiddleware sets headers to prevent browser caching
func noCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

// htmxMiddleware adds HTMX-specific functionality to requests
func (app *App) htmxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add HTMX request information to context
		ctx := r.Context()

		// Check if this is an HTMX request
		if r.Header.Get("HX-Request") == "true" {
			ctx = context.WithValue(ctx, ctxKey("is_htmx_request"), true)
			ctx = context.WithValue(ctx, ctxKey("hx_target"), r.Header.Get("HX-Target"))
			ctx = context.WithValue(ctx, ctxKey("hx_trigger"), r.Header.Get("HX-Trigger"))
			ctx = context.WithValue(ctx, ctxKey("hx_current_url"), r.Header.Get("HX-Current-URL"))
		}

		// Set HTMX-compatible headers
		w.Header().Set("Vary", "HX-Request")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loggingMiddleware provides structured logging for requests
func (app *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			duration := time.Since(start)

			app.logger.Info("HTTP request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration_ms", duration.Milliseconds(),
				"bytes", ww.BytesWritten(),
				"user_agent", r.UserAgent(),
				"remote_addr", r.RemoteAddr,
				"request_id", middleware.GetReqID(r.Context()),
				"is_htmx", r.Header.Get("HX-Request") == "true",
			)
		}()

		next.ServeHTTP(ww, r)
	})
}

// Health check handlers
func (app *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok","timestamp":"`)
	fmt.Fprint(w, time.Now().Format(time.RFC3339))
	fmt.Fprint(w, `"}`)
}

func (app *App) readinessHandler(w http.ResponseWriter, r *http.Request) {
	// Check dependencies
	checks := []struct {
		name string
		ok   bool
	}{
		{"encounter_service", app.encounterService != nil},
		{"repository", true}, // Memory repo is always available
	}

	allReady := true
	for _, check := range checks {
		if !check.ok {
			allReady = false
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if !allReady {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, `{"status":"not ready","checks":[`)
		for i, check := range checks {
			if i > 0 {
				fmt.Fprint(w, ",")
			}
			status := "ok"
			if !check.ok {
				status = "fail"
			}
			fmt.Fprintf(w, `{"name":"%s","status":"%s"}`, check.name, status)
		}
		fmt.Fprint(w, `]}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ready","timestamp":"`)
	fmt.Fprint(w, time.Now().Format(time.RFC3339))
	fmt.Fprint(w, `"}`)
}

// Server creates and configures the HTTP server
func (app *App) Server() *http.Server {
	return &http.Server{
		Addr:         app.config.ServerAddr(),
		Handler:      app.router,
		ReadTimeout:  app.config.ReadTimeout,
		WriteTimeout: app.config.WriteTimeout,
		IdleTimeout:  app.config.IdleTimeout,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
}

// Run starts the application server with graceful shutdown
func (app *App) Run() error {
	server := app.Server()

	// Channel to capture server errors
	serverErrors := make(chan error, 1)

	// Channel to capture shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		app.logger.Info("Server starting",
			"address", server.Addr,
			"environment", app.config.Environment,
		)
		serverErrors <- server.ListenAndServe()
	}()

	// Wait for either an error or shutdown signal
	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server error: %w", err)
		}
		app.logger.Info("Server stopped due to closure")
		return nil

	case sig := <-shutdown:
		app.logger.Info("Shutdown signal received", "signal", sig)

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(
			context.Background(),
			app.config.ShutdownTimeout,
		)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			app.logger.Error("Graceful shutdown failed", "error", err)

			// Force close if graceful shutdown fails
			if closeErr := server.Close(); closeErr != nil {
				app.logger.Error("Forced server close failed", "error", closeErr)
			}
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		app.logger.Info("Server gracefully stopped")
		return nil
	}
}

// setupLogger creates and configures the application logger
func setupLogger(cfg *config.Config) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if cfg.IsProduction() {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Setup logger
	logger := setupLogger(cfg)
	slog.SetDefault(logger)

	// Create application
	app, err := NewApp(cfg, logger)
	if err != nil {
		logger.Error("Failed to create application", "error", err)
		os.Exit(1)
	}

	// Run application
	if err := app.Run(); err != nil {
		logger.Error("Application error", "error", err)
		os.Exit(1)
	}
}
