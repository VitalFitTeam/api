package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/vitalfit/api/docs"
	"github.com/vitalfit/api/pkg/cors"

	"github.com/vitalfit/api/config"
	apphandlers "github.com/vitalfit/api/internal/app/handlers"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/store"
	"go.uber.org/zap"
)

var (
	version = "0.0.1"
)

type application struct {
	Config   *config.Config
	Logger   *zap.SugaredLogger
	store    store.Storage
	services appservices.Services
	handlers apphandlers.Handlers
}

// Mount config and return router
func (app *application) Mount() http.Handler {
	r := gin.New()
	docs.SwaggerInfo.BasePath = "/v1"
	r.Use(gin.Logger(), gin.Recovery())
	cors.SetupCORS(r)
	{

		v1 := r.Group("/v1")

		v1.GET("/health", app.HealthCheckHandler)

		app.handlers.AuthHandlers.AuthRoutes(v1)

		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	}

	return r
}

// Run starts HTTP server
func (app *application) Run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.Config.Addrs,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.Logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.Logger.Infow("server has started", "addr", app.Config.Addrs, "env", app.Config.Env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.Logger.Infow("server has stopped", "addr", app.Config.Addrs, "env", app.Config.Env)
	return nil
}
