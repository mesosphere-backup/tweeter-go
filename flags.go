package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)


type flags struct {
	address *string
	cassandraAddr *string
	cassandraSRV *string
	cassandraDNS *string
	cassandraRepl *int
}

func (c *flags) addSet(s *flag.FlagSet) {
	cassandraAddr := s.String("cassandra-addr", "", "Address to a single Cassandra node")
	c.cassandraAddr = cassandraAddr

	cassandraRepl := s.Int("cassandra-repl", -1, "Replication factor to use for the oinker keyspace in Cassandra (defaults to max(3, floor(len(SRV records)/2.0)+1)")
	c.cassandraRepl = cassandraRepl

	cassandraSRV := s.String("cassandra-srv", "", "DNS service name of the Cassandra cluster")
	c.cassandraSRV = cassandraSRV

	cassandraDNS := s.String("cassandra-dns", "", "Domain name server to lookup the Cassandra service")
	c.cassandraDNS = cassandraDNS

	address := s.String("address", "0.0.0.0:8080", "host:port on which to listen")
	c.address = address
}

func usage(s *flag.FlagSet) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n", filepath.Base(os.Args[0]))
		s.PrintDefaults()
	}
}

func parseFlags(s *flag.FlagSet) *flags {
	c := &flags{}
	c.addSet(s)
	s.Usage = usage(s)
	s.Parse(os.Args[1:])
	return c
}
