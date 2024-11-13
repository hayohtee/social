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

	app.logger.Infow("server has started", "addr", srv.Addr, "env", app.config.env)

	return srv.ListenAndServe()
}
