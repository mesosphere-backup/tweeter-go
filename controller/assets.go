package controller

import (
	"net/http"
)

type AssetsController struct {
}

func NewAssetsController() *AssetsController {
	return &AssetsController{}
}

func (c *AssetsController) Name() string {
	return "AssetsController"
}

func (c *AssetsController) RegisterHandlers(server MuxServer) {
	server.HandleFunc("/assets/", http.FileServer(http.Dir("./")).ServeHTTP)
}
