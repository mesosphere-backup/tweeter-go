package controller

import (
	"github.com/mesosphere/tweeter-go/service"

	log "github.com/Sirupsen/logrus"

	"net/http"
	"fmt"
)

type ReadyController struct {
	svc service.Service
	instanceName string
}

func NewReadyController(svc service.Service, instanceName string) *ReadyController {
	return &ReadyController{
		svc: svc,
		instanceName: instanceName,
	}
}

func (c *ReadyController) Name() string {
	return "ReadyController"
}

func (c *ReadyController) RegisterHandlers(server MuxServer) {
	server.HandleFunc("/ready", c.Handle)
}

func (c *ReadyController) Handle(w http.ResponseWriter, r *http.Request) {
	err := c.checkReady()
	if err != nil {
		log.Errorf("Ready check failed (%d) handling request (%s, %s, %s): %s", err.Code(), c.Name(), r.Method, r.URL.Path, err)
		http.Error(w, err.Error(), err.Code())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	fmt.Fprintf(w, "Ready\nInstance: %s", c.instanceName)
}

func (c *ReadyController) checkReady() HTTPError {
	err := c.svc.CheckReady()
	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Checking repo readiness: %s", err))
	}
	return nil
}
