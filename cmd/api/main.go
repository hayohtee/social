package main

import "log"

func main() {
	cfg := config{
		addr: ":8080",
	}

	app := &application{
		config: cfg,
	}

	routes := app.routes()

	if err := app.serve(routes); err != nil {
		log.Fatal(err)
	}
}
