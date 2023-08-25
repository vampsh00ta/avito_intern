package service

import (
	repository "avito/internal/db"
	"avito/internal/service/dto"
	"context"
)

type Service interface {
	User
}
type service struct {
	r repository.Repository
}

func New(r repository.Repository) service {
	return service{
		r,
	}
}

type User interface {
	CreateUser(ctx context.Context, user dto.User) error
	DeleteUser(ctx context.Context, userId int) error
	AddTagsToUser(ctx context.Context, userId int, slugs ...any) error
	GetUsersTags(ctx context.Context, userId int) ([]repository.Segment, error)
}
