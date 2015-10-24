package main

import (
	"github.com/karlkfi/oinker-go/controller"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/karlkfi/inject"
	log "github.com/Sirupsen/logrus"

	"net/http"
	"flag"
	"net"
	"fmt"
	"strings"
	"github.com/karlkfi/oinker-go/model"
	"strconv"
)

func main() {
	flagSet := flag.CommandLine
	flags := parseFlags(flagSet)
	log.Infof("Flags: %+v", flags)

	oinker := &Oinker{}

	if *flags.cassandraAddr != "" {
		oinker.CQLHosts = []string{*flags.cassandraAddr}
	} else if *flags.cassandraSRV != "" && *flags.cassandraDNS == "" || *flags.cassandraSRV == "" && *flags.cassandraDNS != "" {
		log.Fatalf("Invalid input: cassandra-srv and cassandra-dns must both be specified to enable cassandra usage")
	} else if *flags.cassandraSRV != "" && *flags.cassandraDNS != "" {
		hosts, err := lookupCassandraHosts(*flags.cassandraSRV, *flags.cassandraDNS)
		if err != nil {
			log.Fatalf("Error looking up Cassandra SRV records: %s", err)
		}
		oinker.CQLHosts = hosts
	}

	oinker.SetCQLReplicationFactor(*flags.cassandraRepl)

	graph := oinker.NewGraph()
	defer graph.Finalize()

	// initialize cassandra (connection, keyspace, tables)
	var oinkRepo model.OinkRepo
	inject.ExtractAssignable(graph, &oinkRepo)
	svc, ok := oinkRepo.(inject.Initializable)
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

func lookupCassandraHosts(service, dns string) ([]string, error) {
	cname, srvs, err := net.LookupSRV(service, "tcp", dns)
	if err != nil {
		return nil, fmt.Errorf("Looking up SRV record (srv: %s, dns: %s): %s", service, dns, err)
	}
	log.Infof("CNAME: %s", cname)
	log.Infof("SRVs: %+v", srvs)

	if len(srvs) == 0 {
		return nil, fmt.Errorf("No SRV records found (srv: %s, dns: %s)", service, dns)
	}

	addrs := make([]string, len(srvs), len(srvs))
	for i, srv := range srvs {
		addr := srv.Target
		if strings.HasSuffix(srv.Target, ".") {
			addr = addr[:len(addr)-1]
		}
		if srv.Port > 0 {
			addr = addr+":"+strconv.FormatUint(uint64(srv.Port), 10)
		}
		addrs[i] = addr
	}
	return addrs, nil
}
