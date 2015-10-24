package service

import (
	"github.com/gocql/gocql"
	log "github.com/Sirupsen/logrus"

	"sync"
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

	log.Info("Creating cql session")

	session, err := s.cluster.CreateSession()
	if err != nil {
		log.Errorf("Creating CQL Session: %s", err)
		return
	}
	s.Session = *session

	s.initialized = true
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
