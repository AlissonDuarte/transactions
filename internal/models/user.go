package models

type User struct {
	Name string `validate:"required, min=5, max=100"`
	Cpf string `validate:"required,cpf"`
	Email string `validate:"required, email"`
	Password string `validate:"required, min=6, max=100"`
}

func (u *User) Validate() error {
	return validate.Struct(u)
}

func validateCpf(cpf validator.FieldLevel) bool {
	cpfRegex := regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
	return cpfRegex.MatchString(cpf)
}