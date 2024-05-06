package model

type Login struct {
	Email    string
	Password string
}

type UserResponseDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
