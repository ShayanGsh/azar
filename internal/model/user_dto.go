package model

type UpdateUserData struct {
	Username    string `json:"username" validate:"required_without=Email,omitempty,min=1,max=100"`
	NewUsername string `json:"new_username" validate:"omitempty,min=1,max=100"`
	Email       string `json:"email" validate:"required_without=Username,omitempty,email"`
	NewEmail    string `json:"new_email" validate:"omitempty,email"`
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"omitempty,min=8"`
}

type UserData struct {
	Username string `json:"username" validate:"required_without=Email,omitempty,min=1,max=100"`
	Email    string `json:"email" validate:"required_without=Username,omitempty,email"`
	Password string `json:"password" validate:"required,min=8"`
}