package repository

import (
	"avito/pkg/client"
	"go.uber.org/zap"
)

type Db struct {
	client postgresql.Client
	log    *zap.Logger
}

func New(client postgresql.Client, logger *zap.Logger) Repository {
	return &Db{
		client,
		logger,
	}
}
