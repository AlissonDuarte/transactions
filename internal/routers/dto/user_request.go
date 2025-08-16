package dto

type CreateUserRequest struct {
	Name     string `validate:"required,min=5,max=100"`
	Cpf      string `validate:"required,cpf"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6,max=100"`
}
