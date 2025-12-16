package domain

type PasswordValidator interface {
	Validate(password string) error
}
