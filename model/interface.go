package model

type OinkRepo interface {
	Init() error
	Create(Oink) (Oink, error)
	FindByID(id string) (Oink, bool, error)
	All() ([]Oink, error)
	Analytics() ([]Analytics, error)
}
