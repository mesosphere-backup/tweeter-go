package main

import (
	"github.com/karlkfi/oinker-go/controller"
	"github.com/karlkfi/oinker-go/model"

	"github.com/gocql/gocql"
	"github.com/karlkfi/inject"

	"net/http"
)

type Oinker struct {
	CQLHosts []string
}

func (o *Oinker) NewGraph() inject.Graph {
	graph := inject.NewGraph()

	var server *http.ServeMux
	graph.Define(&server, inject.NewProvider(http.NewServeMux))

	var cqlCluster *gocql.ClusterConfig
	graph.Define(&cqlCluster, inject.NewProvider(func() *gocql.ClusterConfig {
		//TODO: use DiscoverHosts?
		return gocql.NewCluster(o.CQLHosts...)
	}))

	var oinkRepo model.OinkRepo
	if len(o.CQLHosts) > 0 {
		graph.Define(&oinkRepo, inject.NewProvider(model.NewCQLOinkRepo, &cqlCluster))
	} else {
		graph.Define(&oinkRepo, inject.NewProvider(model.NewMockOinkRepo))
	}

	var assetsController *controller.AssetsController
	graph.Define(&assetsController, inject.NewProvider(controller.NewAssetsController))

	var indexController *controller.IndexController
	graph.Define(&indexController, inject.NewProvider(controller.NewIndexController, &oinkRepo))

	var oinkController *controller.OinkController
	graph.Define(&oinkController, inject.NewProvider(controller.NewOinkController, &oinkRepo))

	var analyticsController *controller.AnalyticsController
	graph.Define(&analyticsController, inject.NewProvider(controller.NewAnalyticsController, &oinkRepo))

	return graph
}
