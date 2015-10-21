package controller

import (
	"github.com/karlkfi/oinker-go/model"

	"net/http"
	"fmt"
	"encoding/json"
	"strings"
	"log"
)

type OinkController struct {
	repo *model.OinkRepo
}

func NewOinkController(repo *model.OinkRepo) *OinkController {
	return &OinkController{
		repo: repo,
	}
}

func (c *OinkController) Name() string {
	return "OinkController"
}

func (c *OinkController) RegisterHandlers(server MuxServer) {
	server.HandleFunc("/oink", c.Handle)
}

func (c *OinkController) Handle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Path:", r.URL.Path, "Method:", r.Method, "Form:", r.Form)
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

func (c *OinkController) Get(w http.ResponseWriter, r *http.Request) {
	subPath := strings.Replace(r.URL.Path, "/oink/", "", 1)
	if subPath == "" {
		bytes, err := json.Marshal(c.repo.List())
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
	oink, found := c.repo.FindByID(id)
	if !found {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid ID: %s", id)
		return
	}

	bytes, err := json.Marshal(oink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (c *OinkController) Post(w http.ResponseWriter, r *http.Request) {
	handle := r.Form.Get("handle")
	content := r.Form.Get("content")

	oink := c.repo.Add(model.Oink{
		Handle: handle,
		Content: content,
	})
	log.Printf("Added Oink: %+v\n", oink)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
