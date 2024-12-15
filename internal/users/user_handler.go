package users

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(user User) (User, error) {
	return h.userService.CreateUser(user)
}

func (h *UserHandler) GetUserById(id int64) (User, error) {
	return h.userService.GetUserById(id)
}

func (h *UserHandler) UpdateUser(user User) (User, error) {
	return h.userService.UpdateUser(user)
}

func (h *UserHandler) DeleteUser(id int64) error {
	return h.userService.DeleteUser(id)
}
