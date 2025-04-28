package app

import (
	"context"
	"database/sql"
	"effective-mobile-test-task/internal/configs"
	"effective-mobile-test-task/internal/handler"
	"effective-mobile-test-task/internal/httpclient"
	psqlImpl "effective-mobile-test-task/internal/repository/postgres"
	"effective-mobile-test-task/internal/service"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type AppBuilder struct {
	logger zerolog.Logger
	router *chi.Mux
	server *http.Server
	db     *sql.DB
	err    error
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{}
}

func (b *AppBuilder) error(err error) *AppBuilder {
	b.err = err
	return b
}

func (b *AppBuilder) WithEnv() *AppBuilder {
	if os.Getenv("ENV") != "prod" {
		if err := godotenv.Load(); err != nil {
			b.err = err
		}
	}
	return b
}

func (b *AppBuilder) WithLogger() *AppBuilder {
	var zl zerolog.Level

	if os.Getenv("ENV") == "prod" {
		zl = zerolog.InfoLevel
	} else {
		zl = zerolog.DebugLevel
	}

	b.logger = zerolog.New(os.Stdout).Level(zl).With().Timestamp().Logger()
	return b
}

func (b *AppBuilder) WithRouter() *AppBuilder {
	b.router = chi.NewRouter()
	b.router.Use(middleware.Recoverer)
	b.router.Use(middleware.RequestID)
	b.router.Use(hlog.NewHandler(b.logger))
	b.router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, dur time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", status).
			Dur("duration", dur).
			Msg("request completed")
	}))
	b.router.Use(hlog.RequestIDHandler("request_id", "X-Request-ID"))
	return b
}

func (b *AppBuilder) WithDatabase() *AppBuilder {
	dsn, err := configs.GetPostgresDSN()
	if err != nil {
		return b.error(err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return b.error(err)
	}
	b.db = db
	return b
}

func (b *AppBuilder) WithMigrations() *AppBuilder {
	err := goose.Up(b.db, "./migrations")
	if err != nil {
		b.error(err)
	}

	return b
}

func (b *AppBuilder) WithUserRouter() *AppBuilder {
	userRepo, err := psqlImpl.NewUserRepo(b.db)
	if err != nil {
		return b.error(err)
	}

	agifyConfig, err := configs.GetAgifyConfig()
	if err != nil {
		return b.error(err)
	}
	genderizeConfig, err := configs.GetGenderizeConfig()
	if err != nil {
		return b.error(err)
	}
	nationalizeConfig, err := configs.GetNationalizeConfig()
	if err != nil {
		return b.error(err)
	}

	agifyClient, err := httpclient.NewPredictorClient[httpclient.AgifyResponse](*agifyConfig)
	if err != nil {
		return b.error(err)
	}
	genderizeClient, err := httpclient.NewPredictorClient[httpclient.GenderizeResponse](*genderizeConfig)
	if err != nil {
		return b.error(err)
	}
	nationalizeClient, err := httpclient.NewPredictorClient[httpclient.NationalizeResponse](*nationalizeConfig)
	if err != nil {
		return b.error(err)
	}

	userService, err := service.NewUserService(userRepo, agifyClient, genderizeClient, nationalizeClient)
	if err != nil {
		return b.error(err)
	}
	userHandler, err := handler.NewUserHandler(userService)
	if err != nil {
		return b.error(err)
	}

	b.router.Mount("/users", userHandler.Routes())
	return b
}

func (b *AppBuilder) WithServer() *AppBuilder {
	if b.err != nil {
		return b
	}
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	b.server = &http.Server{
		Addr:    addr,
		Handler: b.router,
	}
	return b
}

func (b *AppBuilder) Build() error {
	return b.err
}

func (b *AppBuilder) Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		b.logger.Info().Msg("server running on " + b.server.Addr)
		if err := b.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			b.logger.Error().Err(err).Msg("error while starting server")
		}
	}()

	<-stop
	b.logger.Info().Msg("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := b.server.Shutdown(ctx); err != nil {
		b.logger.Error().Err(err).Msg("server shutdown failed")
		return
	}
	if err := b.db.Close(); err != nil {
		b.logger.Error().Err(err).Msg("db connection close failed")
		return
	}

	b.logger.Info().Msg("server exited properly")
}
