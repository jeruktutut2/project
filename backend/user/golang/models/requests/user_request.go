package requests

type RegisterUserRequest struct {
	Username        string `json:"username" validate:"required,usernamevalidator"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,passwordvalidator"`
	Confirmpassword string `json:"confirmpassword" validate:"required,passwordvalidator"`
	Utc             string `json:"utc" validate:"required,gte=5"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwordvalidator"`
}
