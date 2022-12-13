package registry

import (
	"user-service/interface/repository"
	"user-service/usecases/interactor"
	repo "user-service/usecases/repository"
)

func (r *registry) NewUserInteractor() interactor.UserInteractor {
	return interactor.NewUserInteractor(r.NewUserRepository())
}

func (r *registry) NewUserRepository() repo.UserRepository {
	return repository.NewUserRepository(r.Db)
}
