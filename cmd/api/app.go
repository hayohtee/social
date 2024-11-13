package main

import (
	"github.com/hayohtee/social/internal/repository"
	"go.uber.org/zap"
)

type application struct {
	config     config
	repository repository.Repository
	logger     *zap.SugaredLogger
}
