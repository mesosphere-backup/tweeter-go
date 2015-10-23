package controller

import (
	"github.com/karlkfi/oinker-go/model"
	"github.com/karlkfi/oinker-go/view"

	"fmt"
	"html/template"
	"log"
	"net/http"
)

type AnalyticsController struct {
	repo model.OinkRepo
}

func NewAnalyticsController(repo model.OinkRepo) *AnalyticsController {
	return &AnalyticsController{
		repo: repo,
	}
}

func (c *AnalyticsController) Name() string {
	return "AnalyticsController"
}

func (c *AnalyticsController) RegisterHandlers(server MuxServer) {
	server.HandleFunc("/analytics", c.Handle)
}

func (c *AnalyticsController) Handle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Path:", r.URL.Path, "Method:", r.Method, "Form:", r.Form)
	switch r.Method {
	case "GET":
		c.Get(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Invalid method: %s", r.Method)
	}
}

func (c *AnalyticsController) Get(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New(
		"layout.html.tmpl",
	).ParseFiles(
		"templates/layout.html.tmpl",
		"templates/analytics.html.tmpl",
	))

	a, err := c.repo.Analytics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Analytics: %+v\n", a)

	err = t.Execute(w, view.Analytics{
		Page: view.Page{
			RelativeRootPath: ".",
		},
		Analytics: a,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
