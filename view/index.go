package view

import (
	"github.com/mesosphere/tweeter-go/model"
)

type Index struct {
	Page
	Tweets []model.Tweet
	IsEmpty bool
}
