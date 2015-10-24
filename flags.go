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
	cassandraRepl *int
}

func (c *flags) addSet(s *flag.FlagSet) {
	cassandraAddr := s.String("cassandra-addr", "", "Address to a single Cassandra node")
	c.cassandraAddr = cassandraAddr

	cassandraRepl := s.Int("cassandra-repl", 1, "Replication factor to use for the oinker keyspace in Cassandra")
	c.cassandraRepl = cassandraRepl

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
