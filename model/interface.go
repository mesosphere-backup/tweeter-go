package model

type OinkRepo interface {
	Create(Oink) (Oink, error)
	FindByID(id string) (Oink, bool, error)
	All() ([]Oink, error)
}
