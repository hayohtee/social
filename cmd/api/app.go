package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/hayohtee/social/internal/repository"
)

type application struct {
	config     config
	repository repository.Repository
	validate   *validator.Validate
}
