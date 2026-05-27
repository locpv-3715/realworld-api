package dto

type RegisterReq struct {
	User struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	} `json:"user"`
}

type LoginReq struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	} `json:"user"`
}

type UserRes struct {
	User struct {
		Email    string `json:"email"`
		Token    string `json:"token"`
		Username string `json:"username"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"user"`
}

type UpdateUserReq struct {
	User struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
		Image    string `json:"image"`
		Bio      string `json:"bio"`
	} `json:"user"`
}
