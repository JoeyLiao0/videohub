package service

import "videohub/internal/repository"

type UserCollection struct {
	collectionRepo *repository.Collection
}

func NewUserCollection(cr *repository.Collection) *UserCollection {
	return &UserCollection{collectionRepo: cr}
}
