// Package httpapp реализует HTTP сервер с возможностью настройки параметров
// через func options. Сервер использует фреймворк Gin для обработки HTTP запросов.
package httpapp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	teamHandler "avitotech-pr-reviewer/internal/api/v1/team"
	usersHandler "avitotech-pr-reviewer/internal/api/v1/users"
	teamService "avitotech-pr-reviewer/internal/service/team"

	"github.com/gin-gonic/gin"
)

const (
	srvPortDefault           = 8080
	srvReadTimeoutDefault    = 10 * time.Second
	srvWriteTimeoutDefault   = 10 * time.Second
	srvGatewayTimeoutDefault = 10 * time.Second
)

type App struct {
	lgr *slog.Logger

	port           int
	readTimeout    time.Duration
	writeTimeout   time.Duration
	requestTimeout time.Duration

	mu     sync.Mutex
	server *http.Server
}

type Option func(*App)

func WithPort(port int) Option {
	return func(a *App) {
		if port <= 0 || port > 65535 {
			a.lgr.Error("invalid port number, using default port", slog.Int("default_port", srvPortDefault))
			a.port = srvPortDefault
		} else {
			a.port = port
		}
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(a *App) {
		a.readTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(a *App) {
		a.writeTimeout = timeout
	}
}

func WithRequestTimeout(timeout time.Duration) Option {
	return func(a *App) {
		a.requestTimeout = timeout
	}
}

// New создает новый экземпляр HTTP сервера с заданными опциями.
func New(lgr *slog.Logger, opts ...Option) *App {
	app := &App{
		lgr:            lgr,
		readTimeout:    srvReadTimeoutDefault,
		writeTimeout:   srvWriteTimeoutDefault,
		requestTimeout: srvGatewayTimeoutDefault,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// MustRun запускает HTTP сервер и паникует в случае ошибки.
func (a *App) MustRun(ctx context.Context) {
	err := a.Run(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic("failed to run HTTP server: " + err.Error())
	}
}

// Run запускает HTTP сервер.
func (a *App) Run(ctx context.Context) error {
	const op = "httpapp.Run"

	lgr := a.lgr.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	lgr.InfoContext(ctx, "starting HTTP http_server")

	teamSvc := teamService.New(lgr)

	teamHlr := teamHandler.New(teamSvc)
	usersHlr := usersHandler.New()

	app := gin.New()
	app.Use(gin.Recovery())

	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	teamHlr.RegisterRoutes(v1)
	usersHlr.RegisterRoutes(v1)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.port),
		Handler:      app,
		ReadTimeout:  a.readTimeout,
		WriteTimeout: a.writeTimeout,
	}

	a.server = srv

	return srv.ListenAndServe()
}

// Stop останавливает HTTP сервер.
// Нужно дожидаться завершения работы этого метода.
// Контекст должен быть с тайм-аутом, чтобы избежать
// зависания в случае проблем с остановкой сервера.
func (a *App) Stop(ctx context.Context) error {
	const op = "httpapp.Stop"

	lgr := a.lgr.With("op", op)

	lgr.Info("stopping HTTP http_server")

	a.mu.Lock()
	server := a.server
	a.mu.Unlock()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
