package user

import (
	userRepository "github.com/OCCASS/avito-intern/internal/domain/user/repository"
)

type UserServices struct {
	userRepository userRepository.UserRepository
}

func NewUserServices(
	ur userRepository.UserRepository,
) UserServices {
	return UserServices{
		userRepository: ur,
	}
}
