package service

import (
	"github.com/karlkfi/oinker-go/model"

	"sync"
	"strconv"
	"time"
)

// MockOinkRepo is an in-memory OinkRepo
type MockOinkRepo struct {
	lock sync.RWMutex
	oinks []model.Oink
}

func NewMockOinkRepo() *MockOinkRepo {
	return &MockOinkRepo{
		oinks: []model.Oink{},
	}
}

func (r *MockOinkRepo) Create(o model.Oink) (model.Oink, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	o.ID = strconv.Itoa(len(r.oinks))
	o.CreationTime = time.Now()
	r.oinks = append(r.oinks, o)
	return o, nil
}

func (r *MockOinkRepo) FindByID(id string) (model.Oink, bool, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, oink := range r.oinks {
		if oink.ID == id {
			return oink, true, nil
		}
	}
	return model.Oink{}, false, nil
}

// Returns all oinks, from newest to oldest
func (r *MockOinkRepo) All() ([]model.Oink, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	result := make([]model.Oink, 0, len(r.oinks))
	for i := len(r.oinks)-1; i >= 0; i-- {
		result = append(result, r.oinks[i])
	}
	return result, nil
}
