package controller

import (
	"net/http"
)

type Controller interface {
	Name() string
	RegisterHandlers(MuxServer)
}

type MuxServer interface {
	http.Handler
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}
