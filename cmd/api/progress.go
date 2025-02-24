package main

import (
	"fmt"
	"net/http"

	"github.com/XDoubleU/essentia/pkg/parse"
)

func (app *Application) progressRoutes(prefix string, mux *http.ServeMux) {
	mux.HandleFunc(
		fmt.Sprintf("GET %s/progress", prefix),
		app.services.WebSocket.Handler(),
	)
	mux.HandleFunc(
		fmt.Sprintf("GET %s/progress/{id}/refresh", prefix),
		app.authAccess(app.refreshProgressHandler),
	)
}
func (app *Application) refreshProgressHandler(_ http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam[string](r, "id", nil)
	if err != nil {
		panic(err)
	}

	_, lastRunTime := app.jobQueue.FetchState(id)
	app.services.WebSocket.UpdateState(id, true, lastRunTime)

	app.jobQueue.ForceRun(id)
}
