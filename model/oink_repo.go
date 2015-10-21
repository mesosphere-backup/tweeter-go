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

func (r *OinkRepo) List() []Oink {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return append([]Oink{}, r.oinks...)
}