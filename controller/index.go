package controller

import (
	"github.com/mesosphere/tweeter-go/service"
	"github.com/mesosphere/tweeter-go/view"

	"github.com/dustin/go-humanize"
	log "github.com/Sirupsen/logrus"

	"net/http"
	"html/template"
	"time"
	"fmt"
)

type IndexController struct {
	repo service.TweetRepo
	instanceName string
}

func NewIndexController(repo service.TweetRepo, instanceName string) *IndexController {
	return &IndexController{
		repo: repo,
		instanceName: instanceName,
	}
}

func (c *IndexController) Name() string {
	return "IndexController"
}

func (c *IndexController) RegisterHandlers(server MuxServer) {
	server.HandleFunc("/", c.Handle)
}

func (c *IndexController) Handle(w http.ResponseWriter, r *http.Request) {
	err := c.handleInner(w, r)
	if err != nil {
		log.Errorf("Error (%d) handling request (%s, %s, %s): %s", err.Code(), c.Name(), r.Method, r.URL.Path, err)
		http.Error(w, err.Error(), err.Code())
	}
}

func (c *IndexController) handleInner(w http.ResponseWriter, r *http.Request) HTTPError {
	t := template.Must(template.New(
		"page.html.tmpl",
	).Funcs(template.FuncMap{
		"timeSince": c.TimeSince,
		"avatarURL": c.AvatarURL,
	}).ParseFiles(
		"templates/page.html.tmpl",
		"templates/index.html.tmpl",
	))

	tweets, err := c.repo.All()
	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Retrieving all tweets: %s", err))
	}
	log.Debugf("Tweets: %+v\n", tweets)

	err = t.Execute(w, view.Index{
		Page: view.Page{
			RelativeRootPath: ".",
			InstanceName: c.instanceName,
		},
		Tweets: tweets,
		IsEmpty: len(tweets) == 0,
	})
	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Rendering templates: %s", err))
	}

	return nil
}

func (c *IndexController) TimeSince(input time.Time) string {
	return humanize.RelTime(input, time.Now(), "ago", "from now")
}

func (c *IndexController) AvatarURL(handle string) string {
	return fmt.Sprintf("//robohash.org/%s.png?size=144x144&amp;bgset=bg2", handle)
}
