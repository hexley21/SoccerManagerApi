package auth

type (
	registerRequestDTO struct {
		Username string `json:"username" validate:"required,username"`
		Password string `json:"password" validate:"required,password"`
		Role     string `json:"role"     validate:"required,userrole"`
	} // @name RegisterRequest

	registerResponseDTO struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} // @name RegisterResponse
)

type (
	loginRequestDTO struct {
		Username string `json:"username" validate:"required,username"`
		Password string `json:"password" validate:"required,password"`
	} // @name LoginRequest

	loginResponseDTO struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} // @name LoginResponse
)

type (
	RefreshResponseDTO struct {
		AccessToken string `json:"access_token"`
	} // @name RefreshResponse
)
