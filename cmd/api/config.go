package main

import "time"

type config struct {
	addr string
	db   dbConfig
	env  string
	mail mailConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type mailConfig struct {
	exp time.Duration
}
