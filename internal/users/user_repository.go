package users

import "gorm.io/gorm"

type UserRepository interface {
	GetById(id int64) (User, error)
	Create(user User) (User, error)
	Update(user User) (User, error)
	Delete(id int64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) GetById(id int64) (User, error) {
	var user User
	err := u.db.First(&user, id).Error
	return user, err
}

func (u *userRepository) Create(user User) (User, error) {
	err := u.db.Create(&user).Error
	return user, err
}

func (u *userRepository) Update(user User) (User, error) {
	err := u.db.Save(&user).Error
	return user, err
}

func (u *userRepository) Delete(id int64) error {
	return u.db.Delete(&User{}, id).Error
}
