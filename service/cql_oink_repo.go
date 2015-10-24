package service

import (
	"github.com/karlkfi/oinker-go/model"

	"github.com/gocql/gocql"
	"github.com/satori/go.uuid"
	log "github.com/Sirupsen/logrus"

	"time"
	"fmt"
	"errors"
	"sync"
)

type CQLOinkRepo struct {
	initialized bool
	lock sync.RWMutex
	session *CQLSession
	replicationFactor int
}

func NewCQLOinkRepo(session *CQLSession, replicationFactor int) *CQLOinkRepo {
	return &CQLOinkRepo{
		session: session,
		replicationFactor: replicationFactor,
	}
}

func (r *CQLOinkRepo) Initialize() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.initialized {
		log.Error("Cannot initialize repo, already initialized")
		return
	}

	log.Infof("Creating keyspace oinker (replication_factor: %d)", r.replicationFactor)

	err := r.session.Query(
		fmt.Sprintf(
			"CREATE KEYSPACE IF NOT EXISTS oinker WITH replication = {'class': 'SimpleStrategy','replication_factor': %d}",
			r.replicationFactor,
		),
	).Exec()
	if err != nil {
		log.Errorf("Creating keyspace (oinker): %s", err)
		return
	}

	log.Info("Creating table oinker.oinks")

	err = r.session.Query(
		"CREATE TABLE IF NOT EXISTS oinker.oinks " +
		"( kind VARCHAR, id VARCHAR, content VARCHAR, created_at timeuuid, handle VARCHAR, PRIMARY KEY (kind, created_at) ) " +
		"WITH CLUSTERING ORDER BY (created_at DESC)",
	).Exec()
	if err != nil {
		log.Errorf("Creating table (oinker.oinks): %s", err)
		return
	}

	log.Info("Creating table oinker.analytics")

	err = r.session.Query(
		"CREATE TABLE IF NOT EXISTS oinker.analytics " +
		"( kind VARCHAR, key VARCHAR, frequency INT, PRIMARY KEY (kind, frequency) ) " +
		"WITH CLUSTERING ORDER BY (frequency DESC)",
	).Exec()
	if err != nil {
		log.Errorf("Creating table (oinker.analytics): %s", err)
		return
	}

	r.initialized = true
}

func (r *CQLOinkRepo) Finalize() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if !r.initialized {
		log.Error("Cannot finalize repo, not initialized")
		return
	}
	r.session.Close()
	r.initialized = false
}

func (r *CQLOinkRepo) Create(o model.Oink) (model.Oink, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if !r.initialized {
		return model.Oink{}, errors.New("Uninitialized repo")
	}

	o.ID = uuid.NewV4().String()
	o.CreationTime = time.Now()

	err := r.session.Query(
		"INSERT INTO oinker.oinks (kind, id, content, created_at, handle) " +
		"VALUES (?, ?, ?, ?, ?)",
		"oink", o.ID, o.Content, gocql.UUIDFromTime(o.CreationTime), o.Handle,
	).Exec()
	if err != nil {
		return model.Oink{}, fmt.Errorf("Inserting oink (%+v): %s", o, err)
	}

	return o, nil
}

func (r *CQLOinkRepo) FindByID(id string) (model.Oink, bool, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if !r.initialized {
		return model.Oink{}, false, errors.New("Uninitialized repo")
	}

	iter := r.session.Query("SELECT content, created_at, handle FROM oinker.oinks WHERE id = ?", id).Iter()
	var (
		content string
		creationTime gocql.UUID
		handle string
	)
	oink := model.Oink{
		ID: id,
	}
	if iter.Scan(&content, &creationTime, &handle) {
		oink.Content = content
		oink.CreationTime = creationTime.Time()
		oink.Handle = handle
	}
	if err := iter.Close(); err != nil {
		return model.Oink{}, false, fmt.Errorf("Selecting oink (id: %s): %s", id, err)
	}

	return oink, false, nil
}

// Returns all oinks, from newest to oldest
func (r *CQLOinkRepo) All() ([]model.Oink, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if !r.initialized {
		return nil, errors.New("Uninitialized repo")
	}

	iter := r.session.Query(
		"SELECT id, content, created_at, handle FROM oinker.oinks " +
		"WHERE kind = ? ORDER BY created_at DESC",
		"oink",
	).Iter()
	var (
		id, content, handle string
		creationTime gocql.UUID
	)
	oinks := make([]model.Oink, 0)
	for iter.Scan(&id, &content, &creationTime, &handle) {
		oink := model.Oink{
			ID: id,
			Content: content,
			CreationTime: creationTime.Time(),
			Handle: handle,
		}
		oinks = append(oinks, oink)
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("Selecting oink (id: %s): %s", id, err)
	}

	return oinks, nil
}
