package service

import (
	"github.com/mesosphere/tweeter-go/model"

	"sync"
	"strconv"
	"time"
)

// MockTweetRepo is an in-memory TweetRepo
type MockTweetRepo struct {
	lock sync.RWMutex
	tweets []model.Tweet
}

func NewMockTweetRepo() *MockTweetRepo {
	return &MockTweetRepo{
		tweets: []model.Tweet{},
	}
}

func (r *MockTweetRepo) CheckReady() error {
	// in-memory repo is always ready
	return nil
}

func (r *MockTweetRepo) Create(o model.Tweet) (model.Tweet, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	o.ID = strconv.Itoa(len(r.tweets))
	o.CreationTime = time.Now()
	r.tweets = append(r.tweets, o)
	return o, nil
}

func (r *MockTweetRepo) FindByID(id string) (model.Tweet, bool, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, tweet := range r.tweets {
		if tweet.ID == id {
			return tweet, true, nil
		}
	}
	return model.Tweet{}, false, nil
}

// Returns all tweets, from newest to oldest
func (r *MockTweetRepo) All() ([]model.Tweet, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	result := make([]model.Tweet, 0, len(r.tweets))
	for i := len(r.tweets)-1; i >= 0; i-- {
		result = append(result, r.tweets[i])
	}
	return result, nil
}
