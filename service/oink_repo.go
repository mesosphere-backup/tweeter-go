package service

import "github.com/mesosphere/tweeter-go/model"

type TweetRepo interface {
	Service
	Create(model.Tweet) (model.Tweet, error)
	FindByID(id string) (model.Tweet, bool, error)
	All() ([]model.Tweet, error)
}
