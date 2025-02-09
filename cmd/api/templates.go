package main

import (
	"errors"
	"net/http"

	"github.com/XDoubleU/essentia/pkg/context"
	"github.com/XDoubleU/essentia/pkg/parse"

	"goal-tracker/api/internal/constants"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/tplhelper"
)

func (app *Application) templateRoutes(mux *http.ServeMux) {
	mux.Handle(
		"GET /images/",
		http.FileServerFS(app.images),
	)
	mux.HandleFunc(
		"GET /{$}",
		app.authTemplateAccess(app.rootHandler),
	)
	mux.HandleFunc(
		"GET /link/{id}",
		app.authTemplateAccess(app.linkHandler),
	)
	mux.HandleFunc(
		"GET /goals/{id}",
		app.authTemplateAccess(app.goalProgressHandler),
	)
}

func (app *Application) rootHandler(w http.ResponseWriter, r *http.Request) {
	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	goals, err := app.services.Goals.GetAllGoalsGroupedByStateAndParentGoal(
		r.Context(),
		user.ID,
	)
	if err != nil {
		panic(err)
	}

	tplhelper.RenderWithPanic(app.tpl, w, "root.html", goals)
}

type LinkTemplateData struct {
	Goal    models.Goal
	Sources []models.Source
	Tags    []string
}

func (app *Application) linkHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam[string](r, "id", nil)
	if err != nil {
		panic(err)
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	goal, err := app.services.Goals.GetGoalByID(r.Context(), id, user.ID)
	if err != nil {
		panic(err)
	}

	tags, err := app.services.Goodreads.GetAllTags(r.Context(), user.ID)
	if err != nil {
		panic(err)
	}

	goalAndSources := LinkTemplateData{
		Goal:    *goal,
		Sources: models.Sources,
		Tags:    tags,
	}
	tplhelper.RenderWithPanic(app.tpl, w, "link.html", goalAndSources)
}

func (app *Application) goalProgressHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parse.URLParam[string](r, "id", nil)
	if err != nil {
		panic(err)
	}

	user := context.GetValue[models.User](r.Context(), constants.UserContextKey)
	if user == nil {
		panic(errors.New("not signed in"))
	}

	goal, err := app.services.Goals.GetGoalByID(r.Context(), id, user.ID)
	if err != nil {
		panic(err)
	}

	viewType := models.Types[*goal.TypeID].ViewType
	switch viewType {
	case models.Graph:
		app.graphViewProgress(w, r, goal, user.ID)
	case models.List:
		app.listViewProgress(w, r, goal, user.ID)
	}
}

type GraphData struct {
	Goal           models.Goal
	ProgressLabels []string
	ProgressValues []string
}

func (app *Application) graphViewProgress(
	w http.ResponseWriter,
	r *http.Request,
	goal *models.Goal,
	userID string,
) {
	progressLabels, progressValues, err := app.services.Goals.GetProgressByTypeIDAndDates(
		r.Context(),
		*goal.TypeID,
		userID,
		goal.PeriodStart(true),
		goal.PeriodEnd(),
	)
	if err != nil {
		panic(err)
	}

	graphData := GraphData{
		Goal:           *goal,
		ProgressLabels: progressLabels,
		ProgressValues: progressValues,
	}

	tplhelper.RenderWithPanic(app.tpl, w, "graph.html", graphData)
}

type ListData struct {
	Goal      models.Goal
	ListItems []models.ListItem
}

func (app *Application) listViewProgress(
	w http.ResponseWriter,
	r *http.Request,
	goal *models.Goal,
	userID string,
) {
	listItems, err := app.services.Goals.GetListItemsByGoalID(
		r.Context(),
		goal.ID,
		userID,
	)
	if err != nil {
		panic(err)
	}

	listData := ListData{
		Goal:      *goal,
		ListItems: listItems,
	}
	tplhelper.RenderWithPanic(app.tpl, w, "list.html", listData)
}
