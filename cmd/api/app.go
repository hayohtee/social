package main

import "github.com/hayohtee/social/internal/repository"

type application struct {
	config     config
	repository repository.Repository
}
