package main

import (
	"github.com/mesosphere/tweeter-go/controller"
	"github.com/mesosphere/tweeter-go/service"

	"github.com/karlkfi/inject"
	"github.com/gocql/gocql"

	"net/http"
	"os"
	"time"
)

type Tweeter struct {
	CQLHosts []string
	CQLReplicationFactor int
	CQLReconnectInterval time.Duration
}

func (o *Tweeter) NewGraph() inject.Graph {
	graph := inject.NewGraph()

	var instanceName string
	graph.Define(&instanceName, inject.NewProvider(func() string {
		name := os.Getenv("TWEETER_INSTANCE_NAME")
		if name == "" {
			return "instance-unknown"
		}
		return name
	}))

	var server *http.ServeMux
	graph.Define(&server, inject.NewProvider(http.NewServeMux))

	var cqlCluster *gocql.ClusterConfig
	graph.Define(&cqlCluster, inject.NewProvider(service.NewCQLCluster, &o.CQLHosts, &o.CQLReconnectInterval))

	var cqlSession *service.CQLSession
	graph.Define(&cqlSession, inject.NewProvider(service.NewCQLSession, &cqlCluster))

	var tweetRepo service.TweetRepo
	if len(o.CQLHosts) > 0 {
		graph.Define(&tweetRepo, inject.NewProvider(service.NewCQLTweetRepo, &cqlSession, &o.CQLReplicationFactor))
	} else {
		graph.Define(&tweetRepo, inject.NewProvider(service.NewMockTweetRepo))
	}

	var readyController *controller.ReadyController
	graph.Define(&readyController, inject.NewProvider(controller.NewReadyController, &tweetRepo, &instanceName))

	var assetsController *controller.AssetsController
	graph.Define(&assetsController, inject.NewProvider(controller.NewAssetsController))

	var indexController *controller.IndexController
	graph.Define(&indexController, inject.NewProvider(controller.NewIndexController, &tweetRepo, &instanceName))

	var tweetController *controller.TweetController
	graph.Define(&tweetController, inject.NewProvider(controller.NewTweetController, &tweetRepo))

	return graph
}
