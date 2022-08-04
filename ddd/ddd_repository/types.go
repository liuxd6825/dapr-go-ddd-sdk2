package ddd_repository

type OnSuccess[T any] func(data T) error
type OnSuccessList[T any] func(list []T) error
type OnError func(err error) error
type OnIsFond func() error
