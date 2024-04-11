// Package main ...
package main

import (
	"fmt"
	"net/http"

	"github.com/trevorjamesmartin/goproject/routes"
)

func main() {
	router := routes.NewRouter()
	fmt.Println("Listening @ http://127.0.0.1:3000")
	http.ListenAndServe(":3000", router)
}
