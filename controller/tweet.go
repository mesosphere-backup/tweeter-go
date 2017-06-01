package controller

import (
	"github.com/mesosphere/tweeter-go/model"
	"github.com/mesosphere/tweeter-go/service"

	log "github.com/Sirupsen/logrus"

	"net/http"
	"fmt"
	"encoding/json"
	"strings"
)

type TweetController struct {
	repo service.TweetRepo
}

func NewTweetController(repo service.TweetRepo) *TweetController {
	return &TweetController{
		repo: repo,
	}
}

func (c *TweetController) Name() string {
	return "TweetController"
}

func (c *TweetController) RegisterHandlers(server MuxServer) {
	server.HandleFunc("/tweet", c.Handle)
}

func (c *TweetController) Handle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Debug("Path:", r.URL.Path, "Method:", r.Method, "Form:", r.Form)
	switch r.Method {
	case "GET":
		c.Get(w, r)
	case "POST":
		c.Post(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Invalid method: %s", r.Method)
	}
}

func (c *TweetController) Get(w http.ResponseWriter, r *http.Request) {
	subPath := strings.Replace(r.URL.Path, "/tweet/", "", 1)
	if subPath == "" {
		tweets, err := c.repo.All()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bytes, err := json.Marshal(tweets)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(bytes)
		return
	}

	segments := strings.Split(subPath, "/")
	if len(segments) < 1 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid path: %s", r.URL.Path)
		return
	}

	id := segments[0]
	tweet, found, err := c.repo.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !found {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid ID: %s", id)
		return
	}

	bytes, err := json.Marshal(tweet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (c *TweetController) Post(w http.ResponseWriter, r *http.Request) {
	handle := r.Form.Get("handle")
	content := r.Form.Get("content")

	tweet, err := c.repo.Create(model.Tweet{
		Handle: handle,
		Content: content,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return //TODO: redirect to index with error popup?
	}
	log.Debugf("Added Tweet: %+v\n", tweet)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
