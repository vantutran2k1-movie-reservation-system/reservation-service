package payloads

type CreateUserRequest struct {
	Email    string                   `json:"email" binding:"required,email"`
	Password string                   `json:"password" binding:"required,min=8,max=32"`
	Profile  CreateUserProfileRequest `json:"profile" binding:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type CreatePasswordResetTokenRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8,max=32"`
}
