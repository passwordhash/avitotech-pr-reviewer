// TODO: docs
package httpapp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

func New(lgr *slog.Logger, opts ...Option) *App {
	app := &App{
		lgr: lgr,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (a *App) MustRun(ctx context.Context) {
	err := a.Run(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic("failed to run HTTP server: " + err.Error())
	}
}

func (a *App) Run(ctx context.Context) error {
	const op = "httpapp.Run"

	lgr := a.lgr.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	lgr.InfoContext(ctx, "starting HTTP http_server")

	// инициализация хендлеров ...

	app := gin.New()
	app.Use(gin.Recovery())

	// регистрация маршрутов ...

	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.port),
		Handler:      app,
		ReadTimeout:  a.readTimeout,
		WriteTimeout: a.writeTimeout,
	}

	a.server = srv

	return srv.ListenAndServe()
}
