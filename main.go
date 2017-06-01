package main

import (
	"github.com/mesosphere/tweeter-go/controller"
	"github.com/mesosphere/tweeter-go/service"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/karlkfi/inject"
	log "github.com/Sirupsen/logrus"

	"net/http"
	"flag"
)

func main() {
	flagSet := flag.CommandLine
	flags := parseFlags(flagSet)
	log.Infof("Flags: %+v", flags)

	tweeter := &Tweeter{}

	if *flags.cassandraAddr != "" {
		tweeter.CQLHosts = []string{*flags.cassandraAddr}
		tweeter.CQLReplicationFactor = *flags.cassandraRepl
	}

	graph := tweeter.NewGraph()
	defer graph.Finalize()

	// initialize cassandra (connection, keyspace, tables)
	var tweetRepo service.TweetRepo
	inject.ExtractAssignable(graph, &tweetRepo)
	svc, ok := tweetRepo.(inject.Initializable)
	if ok {
		svc.Initialize()
	}

	var mux controller.MuxServer
	inject.ExtractAssignable(graph, &mux)

	var controllers []controller.Controller
	inject.FindAssignable(graph, &controllers)
	for _, c := range controllers {
		log.Infof("Registering controller:", c.Name())
		c.RegisterHandlers(mux)
	}

	// serve and listen for shutdown signals
	gracehttp.Serve(
		&http.Server{Addr: *flags.address, Handler: mux},
	)
}
