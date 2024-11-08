package main

import (
	"log"
	"net/http"
	"time"

	"github.com/swaggo/swag/example/basic/docs"
)

func (app *application) serve(mux http.Handler) error {
	docs.SwaggerInfo.Version = version

	srv := http.Server {
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout: 10 * time.Second,
		IdleTimeout: time.Minute,
	}

	log.Printf("starting server at %s\n", srv.Addr)

	return srv.ListenAndServe()
}