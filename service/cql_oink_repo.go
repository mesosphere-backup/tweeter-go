package service

import (
	"github.com/mesosphere/tweeter-go/model"

	"github.com/gocql/gocql"
	"github.com/satori/go.uuid"
	log "github.com/Sirupsen/logrus"
	"github.com/cenkalti/backoff"

	"time"
	"fmt"
	"errors"
	"sync"
	"os"
)

type CQLTweetRepo struct {
	initialized bool
	lock sync.RWMutex
	session *CQLSession
	replicationFactor int
}

func NewCQLTweetRepo(session *CQLSession, replicationFactor int) *CQLTweetRepo {
	return &CQLTweetRepo{
		session: session,
		replicationFactor: replicationFactor,
	}
}

func (r *CQLTweetRepo) CheckReady() error {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if !r.initialized {
		return errors.New("Uninitialized repo")
	}

	iter := r.session.Query("SELECT now() FROM tweeter.tweets").Iter()
	if err := iter.Close(); err != nil {
		return fmt.Errorf("Selecting now(): %s", err)
	}

	return nil
}

func (r *CQLTweetRepo) Initialize() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.initialized {
		log.Error("Cannot initialize repo, already initialized")
		return
	}

	err := backoff.Retry(r.initializeAttempt, backoff.NewExponentialBackOff())
	if err != nil {
		log.Errorf("Initializing CQLTweetRepo: %s", err)
		os.Exit(1)
	}

	r.initialized = true
}

func (r *CQLTweetRepo) initializeAttempt() error {
	log.Infof("Attempting to create keyspace tweeter (replication_factor: %d)", r.replicationFactor)

	err := r.session.Query(
		fmt.Sprintf(
			"CREATE KEYSPACE IF NOT EXISTS tweeter WITH replication = {'class': 'SimpleStrategy','replication_factor': %d}",
			r.replicationFactor,
		),
	).Exec()
	if err != nil {
		log.Warnf("Attempt failed to create keyspace (tweeter): %s", err)
		return err
	}

	log.Info("Attempting to create table tweeter.tweets")

	err = r.session.Query(
		"CREATE TABLE IF NOT EXISTS tweeter.tweets " +
		"( kind VARCHAR, id VARCHAR, content VARCHAR, created_at timeuuid, handle VARCHAR, PRIMARY KEY (kind, created_at) ) " +
		"WITH CLUSTERING ORDER BY (created_at DESC)",
	).Exec()
	if err != nil {
		log.Warnf("Attempt failed to create table (tweeter.tweets): %s", err)
		return err
	}

	log.Info("Attempting to create table tweeter.analytics")

	err = r.session.Query(
		"CREATE TABLE IF NOT EXISTS tweeter.analytics " +
		"( kind VARCHAR, key VARCHAR, frequency INT, PRIMARY KEY (kind, frequency) ) " +
		"WITH CLUSTERING ORDER BY (frequency DESC)",
	).Exec()
	if err != nil {
		log.Warnf("Attempt failed to create table (tweeter.analytics): %s", err)
		return err
	}

	return nil
}

func (r *CQLTweetRepo) Finalize() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if !r.initialized {
		log.Error("Cannot finalize repo, not initialized")
		return
	}
	r.session.Close()
	r.initialized = false
}

func (r *CQLTweetRepo) Create(o model.Tweet) (model.Tweet, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if !r.initialized {
		return model.Tweet{}, errors.New("Uninitialized repo")
	}

	o.ID = uuid.NewV4().String()
	o.CreationTime = time.Now()

	err := r.session.Query(
		"INSERT INTO tweeter.tweets (kind, id, content, created_at, handle) " +
		"VALUES (?, ?, ?, ?, ?)",
		"tweet", o.ID, o.Content, gocql.UUIDFromTime(o.CreationTime), o.Handle,
	).Exec()
	if err != nil {
		return model.Tweet{}, fmt.Errorf("Inserting tweet (%+v): %s", o, err)
	}

	return o, nil
}

func (r *CQLTweetRepo) FindByID(id string) (model.Tweet, bool, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if !r.initialized {
		return model.Tweet{}, false, errors.New("Uninitialized repo")
	}

	iter := r.session.Query("SELECT content, created_at, handle FROM tweeter.tweets WHERE id = ?", id).Iter()
	var (
		content string
		creationTime gocql.UUID
		handle string
	)
	tweet := model.Tweet{
		ID: id,
	}
	if iter.Scan(&content, &creationTime, &handle) {
		tweet.Content = content
		tweet.CreationTime = creationTime.Time()
		tweet.Handle = handle
	}
	if err := iter.Close(); err != nil {
		return model.Tweet{}, false, fmt.Errorf("Selecting tweet (id: %s): %s", id, err)
	}

	return tweet, false, nil
}

// Returns all tweets, from newest to oldest
func (r *CQLTweetRepo) All() ([]model.Tweet, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if !r.initialized {
		return nil, errors.New("Uninitialized repo")
	}

	iter := r.session.Query(
		"SELECT id, content, created_at, handle FROM tweeter.tweets " +
		"WHERE kind = ? ORDER BY created_at DESC",
		"tweet",
	).Iter()
	var (
		id, content, handle string
		creationTime gocql.UUID
	)
	tweets := make([]model.Tweet, 0)
	for iter.Scan(&id, &content, &creationTime, &handle) {
		tweet := model.Tweet{
			ID: id,
			Content: content,
			CreationTime: creationTime.Time(),
			Handle: handle,
		}
		tweets = append(tweets, tweet)
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("Selecting tweet (id: %s): %s", id, err)
	}

	return tweets, nil
}
