package controller

import (
	"github.com/karlkfi/oinker-go/model"
	"github.com/karlkfi/oinker-go/view"

	"github.com/dustin/go-humanize"

	"net/http"
	"html/template"
	"time"
	"log"
	"fmt"
)

type IndexController struct {
	repo model.OinkRepo
}

func NewIndexController(repo model.OinkRepo) *IndexController {
	return &IndexController{
		repo: repo,
	}
}

func (c *IndexController) Name() string {
	return "IndexController"
}

func (c *IndexController) RegisterHandlers(server MuxServer) {
	server.HandleFunc("/", c.Handle)
}

func (c *IndexController) Handle(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New(
		"layout.html.tmpl",
	).Funcs(template.FuncMap{
		"timeSince": c.TimeSince,
		"avatarURL": c.AvatarURL,
	}).ParseFiles(
		"templates/layout.html.tmpl",
		"templates/index.html.tmpl",
	))

	oinks, err := c.repo.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Oinks: %+v\n", oinks)

	err = t.Execute(w, view.Index{
		Page: view.Page{
			RelativeRootPath: ".",
		},
		Oinks: oinks,
		IsEmpty: len(oinks) == 0,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *IndexController) TimeSince(input time.Time) string {
	return humanize.RelTime(input, time.Now(), "ago", "from now")
}

func (c *IndexController) AvatarURL(handle string) string {
	return fmt.Sprintf("//robohash.org/%s.png?size=144x144&amp;bgset=bg2", handle)
}
