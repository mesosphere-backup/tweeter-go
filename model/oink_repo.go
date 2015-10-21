package model

import (
	"sync"
	"strconv"
	"time"
)

type OinkRepo struct {
	oinks []Oink
	lock sync.RWMutex
}

func NewOinkRepo() *OinkRepo {
	return &OinkRepo{
		oinks: []Oink{},
	}
}

func (r *OinkRepo) Add(o Oink) Oink {
	r.lock.Lock()
	defer r.lock.Unlock()
	o.ID = strconv.Itoa(len(r.oinks))
	o.CreationTime = time.Now()
	r.oinks = append(r.oinks, o)
	return o
}

func (r *OinkRepo) FindByID(id string) (Oink, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, oink := range r.oinks {
		if oink.ID == id {
			return oink, true
		}
	}
	return Oink{}, false
}

// Returns all oinks, from newest to oldest
func (r *OinkRepo) List() []Oink {
	r.lock.RLock()
	defer r.lock.RUnlock()

	result := make([]Oink, 0, len(r.oinks))
	for i := len(r.oinks)-1; i >= 0; i-- {
		result = append(result, r.oinks[i])
	}
	return result
}

// Returns the last n oinks, from newest to oldest
func (r *OinkRepo) Last(n int) []Oink {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if len(r.oinks) == 0 {
		return []Oink{}
	}

	top := len(r.oinks)-1  // inclusive
	bottom := top-n+1 // inclusive
	if bottom < 0 {
		bottom = 0
	}
	result := make([]Oink, 0, top-bottom)
	for i := top; i >= bottom; i-- {
		result = append(result, r.oinks[i])
	}

	return result
}
