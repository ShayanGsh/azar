package controllers

type UpdateUser struct {
	Username    string `json:"username" validate:"required_without=Email,omitempty,min=1,max=100,excludesall=0x20"`
	Email   string `json:"email" validate:"required_without=Username,omitempty,email"`
	NewEmail    string `json:"new_email" validate:"omitempty,email"`
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"min=8"`
}