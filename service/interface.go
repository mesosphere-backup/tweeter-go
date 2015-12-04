package service

type Service interface {
	CheckReady() error
}
