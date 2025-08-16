package repository

type UserRepository interface {
	FindByCpf(cpf string) (*User, error)
	FindByEmail(email string) (*User, error)
	Save(user *User) error
}