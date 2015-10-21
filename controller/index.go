package controller

import (
	"github.com/karlkfi/oinker-go/model"
	"github.com/karlkfi/oinker-go/view"

	"net/http"
	"html/template"
	"time"
	"log"
	"fmt"
)

type IndexController struct {
	repo *model.OinkRepo
}

func NewIndexController(repo *model.OinkRepo) *IndexController {
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
		"sanitizeHTML": c.SanitizeHTML,
		"avatarURL": c.AvatarURL,
	}).ParseFiles(
		"view/layout.html.tmpl",
		"view/index.html.tmpl",
	))

	oinks := c.repo.List()
	log.Printf("Oinks: %+v\n", oinks)

	err := t.Execute(w, view.Index{
		Oinks: oinks,
		IsEmpty: len(oinks) == 0,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *IndexController) TimeSince(input time.Time) string {
	return time.Now().Sub(input).String()
}

func (c *IndexController) SanitizeHTML(input string) string {
	return input
}

func (c *IndexController) AvatarURL(handle string) string {
	return fmt.Sprintf("//robohash.org/%s.png?size=144x144&amp;bgset=bg2", handle)
}
