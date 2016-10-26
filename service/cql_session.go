package service

import (
	"github.com/gocql/gocql"
	log "github.com/Sirupsen/logrus"
	"github.com/cenkalti/backoff"

	"sync"
	"os"
	"fmt"
)

type CQLSession struct {
	gocql.Session
	initialized bool
	lock sync.Mutex
	cluster *gocql.ClusterConfig
}

func NewCQLSession(cluster *gocql.ClusterConfig) *CQLSession {
	return &CQLSession{
		cluster: cluster,
	}
}

func (s *CQLSession) Initialize() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.initialized {
		log.Error("Cannot initialize cql session, already initialized")
		return
	}

	err := backoff.Retry(s.initializeAttempt, backoff.NewExponentialBackOff())
	if err != nil {
		log.Errorf("Initializing CQLSession: %s", err)
		os.Exit(1)
	}

	s.initialized = true
}

func (s *CQLSession) initializeAttempt() error {
	log.Info("Attempting to create CQL session")
	session, err := s.cluster.CreateSession()
	// inverted error check
	if err != nil {
		log.Warnf("Attempt failed to create CQL Session: %s", err)
		return fmt.Errorf("Creating CQL Session: %s", err)
	}
	s.Session = *session
	return nil
}

func (s *CQLSession) Finalize() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if !s.initialized {
		log.Error("Cannot finalize cql session, not initialized")
		return
	}
	s.Session.Close()
	s.initialized = false
}
