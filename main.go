package main

import (
	"github.com/karlkfi/oinker-go/controller"
	"github.com/karlkfi/oinker-go/model"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/karlkfi/inject"

	"net/http"
	"log"
)

func main() {
	graph := objectGraph()

	var mux controller.MuxServer
	inject.ExtractAssignable(graph, &mux)

	var controllers []controller.Controller
	inject.FindAssignable(graph, &controllers)
	for _, c := range controllers {
		log.Println("Registering controller:", c.Name())
		c.RegisterHandlers(mux)
	}

	// serve and listen for shutdown signals
	gracehttp.Serve(
		&http.Server{Addr: "0.0.0.0:8080", Handler: mux},
	)
}

func objectGraph() inject.Graph {
	var (
		server *http.ServeMux
		oinkRepo *model.OinkRepo
		indexController *controller.IndexController
		oinkController *controller.OinkController
		assetsController *controller.AssetsController
	)

	graph := inject.NewGraph()

	graph.Define(&server, inject.NewProvider(http.NewServeMux))

	graph.Define(&oinkRepo, inject.NewProvider(model.NewOinkRepo))

	graph.Define(&assetsController, inject.NewProvider(controller.NewAssetsController))
	graph.Define(&indexController, inject.NewProvider(controller.NewIndexController, &oinkRepo))
	graph.Define(&oinkController, inject.NewProvider(controller.NewOinkController, &oinkRepo))

	return graph
}
