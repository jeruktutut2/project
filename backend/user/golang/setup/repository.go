package setup

import repository "project-user/repositories"

type RepositorySetup struct {
	UserRepository           repository.UserRepository
	UserPermissionRepository repository.UserPermissionRepository
}

func NewRepositorySetup() *RepositorySetup {
	return &RepositorySetup{
		UserRepository:           repository.NewUserRepository(),
		UserPermissionRepository: repository.NewUserPermissinoRepository(),
	}
}
