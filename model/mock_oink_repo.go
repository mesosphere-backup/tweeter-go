package model

import (
	"strconv"
	"sync"
	"time"
)

// MockOinkRepo is an in-memory OinkRepo
type MockOinkRepo struct {
	oinks     []Oink
	analytics []Analytics
	lock      sync.RWMutex
}

func NewMockOinkRepo() *MockOinkRepo {
	return &MockOinkRepo{
		oinks: []Oink{},
	}
}

func (r *MockOinkRepo) Init() error {
	return nil
}

func (r *MockOinkRepo) Create(o Oink) (Oink, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	o.ID = strconv.Itoa(len(r.oinks))
	o.CreationTime = time.Now()
	r.oinks = append(r.oinks, o)
	return o, nil
}

func (r *MockOinkRepo) FindByID(id string) (Oink, bool, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, oink := range r.oinks {
		if oink.ID == id {
			return oink, true, nil
		}
	}
	return Oink{}, false, nil
}

// Returns all oinks, from newest to oldest
func (r *MockOinkRepo) All() ([]Oink, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	result := make([]Oink, 0, len(r.oinks))
	for i := len(r.oinks) - 1; i >= 0; i-- {
		result = append(result, r.oinks[i])
	}
	return result, nil
}

func (r *MockOinkRepo) Analytics() ([]Analytics, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	//TODO(jdef) key = user:<username>, freq = <total-oinks>
	result := make([]Analytics, len(r.analytics))
	copy(result, r.analytics)
	return result, nil
}
