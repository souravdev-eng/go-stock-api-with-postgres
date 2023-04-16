package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/souravdev-eng/go-postgres-stock/router"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on PORT 8080....")
	log.Fatal(http.ListenAndServe("8080", r))
}
