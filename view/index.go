package view

import (
	"github.com/mesosphere/oinker-go/model"
)

type Index struct {
	Page
	Oinks []model.Oink
	IsEmpty bool
}
