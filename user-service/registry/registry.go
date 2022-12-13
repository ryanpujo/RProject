package registry

import (
	"database/sql"
	"user-service/interface/controller"
	"user-service/user-proto/users"
)

type Registry interface {
	NewUserServer() users.UserServiceServer
}

type registry struct {
	Db *sql.DB
}

func NewRegistry(db *sql.DB) Registry {
	return &registry{Db: db}
}

func (r *registry) NewUserServer() users.UserServiceServer {
	return controller.NewUserServer(r.NewUserInteractor())
}
