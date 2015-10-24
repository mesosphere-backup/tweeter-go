package main

import (
	"github.com/karlkfi/oinker-go/controller"
	"github.com/karlkfi/oinker-go/model"

	"github.com/karlkfi/inject"
	"github.com/gocql/gocql"

	"net/http"
	"github.com/karlkfi/oinker-go/service"
)

const defaultReplicationFactor = 3

type Oinker struct {
	CQLHosts []string
	cqlReplicationFactor int
}

func (o *Oinker) SetCQLReplicationFactor(repl int) {
	if repl > 0 {
		o.cqlReplicationFactor = repl
	} else {
		numHosts := len(o.CQLHosts)
		if numHosts > 0 {
			quorum := int(numHosts / 2.0) + 1
			if quorum > defaultReplicationFactor {
				o.cqlReplicationFactor = quorum
			} else {
				o.cqlReplicationFactor = defaultReplicationFactor
			}
		}
	}
}

func (o *Oinker) NewGraph() inject.Graph {
	graph := inject.NewGraph()

	var server *http.ServeMux
	graph.Define(&server, inject.NewProvider(http.NewServeMux))

	var cqlCluster *gocql.ClusterConfig
	graph.Define(&cqlCluster, inject.NewProvider(func() *gocql.ClusterConfig {
		return gocql.NewCluster(o.CQLHosts...)
	}))

	var cqlSession *service.CQLSession
	graph.Define(&cqlSession, inject.NewProvider(service.NewCQLSession, &cqlCluster))

	var oinkRepo model.OinkRepo
	if len(o.CQLHosts) > 0 {
		graph.Define(&oinkRepo, inject.NewProvider(service.NewCQLOinkRepo, &cqlSession, &o.cqlReplicationFactor))
	} else {
		graph.Define(&oinkRepo, inject.NewProvider(service.NewMockOinkRepo))
	}

	var assetsController *controller.AssetsController
	graph.Define(&assetsController, inject.NewProvider(controller.NewAssetsController))

	var indexController *controller.IndexController
	graph.Define(&indexController, inject.NewProvider(controller.NewIndexController, &oinkRepo))

	var oinkController *controller.OinkController
	graph.Define(&oinkController, inject.NewProvider(controller.NewOinkController, &oinkRepo))

	return graph
}
