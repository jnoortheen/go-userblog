package main

import (
	"fmt"
	"log"
	"net/http"

	"muserblog/actions"

	"github.com/gobuffalo/envy"
)

func main() {
	port := envy.Get("PORT", "3000")
	log.Printf("Starting muserblog on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), actions.App()))
}
