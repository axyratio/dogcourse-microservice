package validators

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type EditProfileInput struct {
	Name     *string `json:"name,omitempty"`
	LastName *string `json:"last_name,omitempty"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
}
