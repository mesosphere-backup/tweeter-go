package model

import (
	"github.com/gocql/gocql"

	"github.com/satori/go.uuid"

	"fmt"
	"log"
	"time"
)

type CQLOinkRepo struct {
	cluster *gocql.ClusterConfig
}

func NewCQLOinkRepo(cluster *gocql.ClusterConfig) *CQLOinkRepo {
	return &CQLOinkRepo{
		cluster: cluster,
	}
}

func (r *CQLOinkRepo) Init() error {
	session, err := r.cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("Creating CQL Session: %s", err)
	}
	defer session.Close()

	replicas := 2
	if len(r.cluster.Hosts) < replicas {
		replicas = len(r.cluster.Hosts)
	}

	log.Printf("Creating keyspace oinker (replication_factor: %d)\n", replicas)

	err = session.Query(
		fmt.Sprintf(
			"CREATE KEYSPACE IF NOT EXISTS oinker WITH replication = {'class': 'SimpleStrategy','replication_factor': %d}",
			replicas,
		),
	).Exec()
	if err != nil {
		return fmt.Errorf("Creating keyspace (oinker): %s", err)
	}

	r.cluster.Keyspace = "oinker"

	log.Println("Creating table oinker.oinks")

	err = session.Query(
		"CREATE TABLE IF NOT EXISTS oinker.oinks " +
			"( kind VARCHAR, id VARCHAR, content VARCHAR, created_at timeuuid, handle VARCHAR, PRIMARY KEY (kind, created_at) ) " +
			"WITH CLUSTERING ORDER BY (created_at DESC)",
	).Exec()
	if err != nil {
		return fmt.Errorf("Creating table (oinker.oinks): %s", err)
	}

	log.Println("Creating table oinker.analytics")

	err = session.Query(
		"CREATE TABLE IF NOT EXISTS oinker.analytics " +
			"( kind VARCHAR, key VARCHAR, frequency INT, PRIMARY KEY (kind, frequency) ) " +
			"WITH CLUSTERING ORDER BY (frequency DESC)",
	).Exec()
	if err != nil {
		return fmt.Errorf("Creating table (oinker.analytics): %s", err)
	}

	return nil
}

func (r *CQLOinkRepo) Create(o Oink) (Oink, error) {
	o.ID = uuid.NewV4().String()
	o.CreationTime = time.Now()

	session, err := r.cluster.CreateSession()
	if err != nil {
		return Oink{}, fmt.Errorf("Creating CQL Session: %s", err)
	}
	defer session.Close()

	err = session.Query(
		"INSERT INTO oinks (kind, id, content, created_at, handle) "+
			"VALUES (?, ?, ?, ?, ?)",
		"oink", o.ID, o.Content, gocql.UUIDFromTime(o.CreationTime), o.Handle,
	).Exec()
	if err != nil {
		return Oink{}, fmt.Errorf("inserting oink (%+v): %s", o, err)
	}

	return o, nil
}

func (r *CQLOinkRepo) FindByID(id string) (Oink, bool, error) {
	session, err := r.cluster.CreateSession()
	if err != nil {
		return Oink{}, false, fmt.Errorf("Creating CQL Session: %s", err)
	}
	defer session.Close()

	iter := session.Query("SELECT content, created_at, handle FROM oinks WHERE id = ?", id).Iter()
	var (
		content      string
		creationTime gocql.UUID
		handle       string
	)
	oink := Oink{
		ID: id,
	}
	if iter.Scan(&content, &creationTime, &handle) {
		oink.Content = content
		oink.CreationTime = creationTime.Time()
		oink.Handle = handle
	}
	if err := iter.Close(); err != nil {
		return Oink{}, false, fmt.Errorf("Selecting oink (id: %s): %s", id, err)
	}

	return oink, false, nil
}

// Returns all oinks, from newest to oldest
func (r *CQLOinkRepo) All() ([]Oink, error) {
	session, err := r.cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("Creating CQL Session: %s", err)
	}
	defer session.Close()

	iter := session.Query(
		"SELECT id, content, created_at, handle FROM oinks "+
			"WHERE kind = ? ORDER BY created_at DESC",
		"oink",
	).Iter()
	var (
		id, content, handle string
		creationTime        gocql.UUID
	)
	oinks := make([]Oink, 0)
	for iter.Scan(&id, &content, &creationTime, &handle) {
		oink := Oink{
			ID:           id,
			Content:      content,
			CreationTime: creationTime.Time(),
			Handle:       handle,
		}
		oinks = append(oinks, oink)
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("Selecting oink (id: %s): %s", id, err)
	}

	return oinks, nil
}

func (r *CQLOinkRepo) Analytics() ([]Analytics, error) {
	session, err := r.cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("Creating CQL Session: %s", err)
	}
	defer session.Close()

	iter := session.Query(
		"SELECT key, frequency FROM analytics "+
			"WHERE kind = ? ORDER BY frequency DESC",
		"oink",
	).Iter()
	var (
		key  string
		freq int
	)
	result := make([]Analytics, 0)
	for iter.Scan(&key, &freq) {
		anal := Analytics{
			Key:  key,
			Freq: freq,
		}
		result = append(result, anal)
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("Selecting analtics: %s", err)
	}

	return result, nil
}
