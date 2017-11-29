package gobastion

import (
	"context"
	"log"
	"net/http"

	"os"

	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ifreddyrondon/gobastion/config"
	"github.com/markbates/sigtx"
)

// Bastion offers an "augmented" Router instance.
// It has the minimal necessary to create an API with default handlers and middleware.
// Allows to have commons handlers and middleware between projects with the need for each one to do so.
// Mounted Routers
// It use go-chi router to modularize the applications. Each instance of GogApp, will have the possibility
// of mounting an API router, it will define the routes and middleware of the application with the app logic.
type Bastion struct {
	r         *chi.Mux
	cfg       *config.Config
	APIRouter chi.Router
}

// NewRouter returns a new Bastion instance.
// if configPath is empty the configuration will be from defaults.
// 	Defaults:
//		api:
//			base_path: "/"
//		server:
//			address ":8080"
// Otherwise the configuration will be loaded from configPath.
// If the config file is missing or unable to unmarshal the will panic.
func NewBastion(configPath string) *Bastion {
	app := new(Bastion)
	app.cfg = config.New()
	if configPath != "" {
		if err := app.cfg.FromFile(configPath); err != nil {
			log.Panic(err)
		}
	}
	initialize(app)
	return app
}

// Serve the application at the specified address/port
func (app *Bastion) Serve() error {
	server := http.Server{Addr: app.cfg.Server.Addr, Handler: app.r}

	ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	// check for a closing signal
	go func() {
		// graceful shutdown
		<-ctx.Done()
		log.Printf("shutting down application")

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("unable to shutdown server: %v", err)
		} else {
			log.Printf("server stopped")
		}
	}()

	// start the web server
	log.Printf("Starting application at %s\n", app.cfg.Server.Addr)
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func initialize(app *Bastion) {
	/**
	 * internal router
	 */
	app.r = chi.NewRouter()

	/**
	 * Ping route
	 */
	app.r.Get("/ping", pingHandler)

	/**
	 * API Router
	 */
	app.APIRouter = chi.NewRouter()
	app.APIRouter.Use(middleware.RequestID)
	app.APIRouter.Use(middleware.Logger)
	app.r.Mount(app.cfg.API.BasePath, app.APIRouter)
}