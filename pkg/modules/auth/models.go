package auth

type RegisterRequest struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=20"`
	LastName  string `json:"lastName" validate:"required,min=2,max=20"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=30"`
}
