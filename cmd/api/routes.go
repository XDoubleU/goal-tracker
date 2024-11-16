package main

import (
	"net/http"

	"github.com/XDoubleU/essentia/pkg/middleware"
	"github.com/getsentry/sentry-go"
	"github.com/justinas/alice"
)

func (app *Application) apiRoutes(mux *http.ServeMux) {
	apiPrefix := "/api"
	app.authRoutes(apiPrefix, mux)
}

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	app.templateRoutes(mux)
	app.apiRoutes(mux)
	//app.goalsRoutes(mux)

	var sentryClientOptions sentry.ClientOptions
	if len(app.config.SentryDsn) > 0 {
		//nolint:exhaustruct //other fields are optional
		sentryClientOptions = sentry.ClientOptions{
			Dsn:                app.config.SentryDsn,
			Environment:        app.config.Env,
			Release:            app.config.Release,
			EnableTracing:      true,
			TracesSampleRate:   app.config.SampleRate,
			ProfilesSampleRate: app.config.SampleRate,
		}
	}

	allowedOrigins := []string{app.config.WebURL}
	handlers, err := middleware.DefaultWithSentry(
		app.logger,
		allowedOrigins,
		app.config.Env,
		sentryClientOptions,
	)

	if err != nil {
		panic(err)
	}

	standard := alice.New(handlers...)
	return standard.Then(mux)
}
