package main

import (
	"net/http"
	"time"
)

func (app *application) serve(mux http.Handler) error {
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infof("starting server at %s", srv.Addr)

	return srv.ListenAndServe()
}
