package user

type userResponseDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
} // @name UserResponse

type updatePasswordRequestDTO struct {
	OldPassword string `json:"old_password" validate:"required,password"`
	NewPassword string `json:"new_password" validate:"required,password"`
} // @name UpdatePasswordRequest

func NewUserResponseDTO(id int64, username string, role string) userResponseDTO {
	return userResponseDTO{
		ID:       id,
		Username: username,
		Role:     role,
	}
}
