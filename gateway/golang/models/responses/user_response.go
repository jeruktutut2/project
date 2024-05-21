package modelresponse

import pbuser "gateway/protofiles/pb/api/v1/user"

type RegisterUserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Utc      string `json:"utc"`
}

func ToRegisterUserResponse(user *pbuser.RegisterResponse) (registerUserResponse RegisterUserResponse) {
	registerUserResponse.Username = user.Username
	registerUserResponse.Email = user.Email
	registerUserResponse.Utc = user.Utc
	return
}
