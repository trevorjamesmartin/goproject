// Package main ...
package main

import (
	"fmt"
	"net/http"

	"github.com/trevorjamesmartin/goproject/middleware"
	"github.com/trevorjamesmartin/goproject/routes"
)

func main() {
	api := middleware.Logging(routes.NewRouter())

	server := http.Server{
		Addr:    ":3000",
		Handler: api,
	}

	fmt.Println("Listening @ http://127.0.0.1:3000")

	server.ListenAndServe()
}
