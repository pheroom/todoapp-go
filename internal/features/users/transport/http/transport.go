package users_transport_http

type UserHTTPHandler struct {
	userService UserService
}

type UserService interface{}

func NewUserHTTPHandler(userService UserService) *UserHTTPHandler {
	return &UserHTTPHandler{userService: userService}
}
