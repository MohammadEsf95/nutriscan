package users

import "gorm.io/gorm"

type UserService interface {
	GetUserById(id int64) (User, error)
	CreateUser(user User) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(id int64) error
}

type userService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetUserById(id int64) (User, error) {
	return s.userRepository.GetById(id)
}

func (s *userService) CreateUser(user User) (User, error) {
	exists, err := s.userRepository.GetById(user.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return User{}, err
	}

	if exists.ID != 0 {
		return User{}, nil
	}
	return s.userRepository.Create(user)
}

func (s *userService) UpdateUser(user User) (User, error) {
	return s.userRepository.Update(user)
}

func (s *userService) DeleteUser(id int64) error {
	return s.userRepository.Delete(id)
}
