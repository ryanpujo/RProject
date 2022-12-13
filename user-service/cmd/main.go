package main

import (
	"user-service/infrastructure"
	"user-service/registry"
)

func main() {
	app := infrastructure.Application()
	db := app.ConnectToDb()
	defer db.Close()
	registry := registry.NewRegistry(db)
	app.StartGrpcServer(registry.NewUserServer())
}
