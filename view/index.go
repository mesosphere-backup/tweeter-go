package view

import (
	"github.com/karlkfi/oinker-go/model"
)

type Index struct {
	Page
	Oinks []model.Oink
	IsEmpty bool
}
