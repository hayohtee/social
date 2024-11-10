package main

type config struct {
	addr   string
	db     dbConfig
	env    string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
