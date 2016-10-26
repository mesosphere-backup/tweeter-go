package service

import (
	"github.com/gocql/gocql"
	"time"
)

// NewCQLCluster generates a new GoSQL config for Cassandra 3, which uses Proto 3
func NewCQLCluster(hosts []string, reconnectInterval time.Duration) *gocql.ClusterConfig {
	cfg := &gocql.ClusterConfig{
		Hosts:                  hosts,
		CQLVersion:             "3.0.0",
		ProtoVersion:           3,
		Timeout:                1500 * time.Millisecond,
		Port:                   9042,
		NumConns:               2,
		Consistency:            gocql.Quorum,
		MaxPreparedStmts:       1000,
		MaxRoutingKeyInfo:      1000,
		PageSize:               5000,
		DefaultTimestamp:       true,
		MaxWaitSchemaAgreement: 60 * time.Second,
		ReconnectInterval:      reconnectInterval,
	}
	return cfg
}
