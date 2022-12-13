package main

import (
	"broker-service/infrastructure/router"
	"broker-service/infrastructure/server"
	"log"
)

func main() {
	if err := server.Serve(router.NewRouter()); err != nil {
		log.Fatal(err)
	}
}
