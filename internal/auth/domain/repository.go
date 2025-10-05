package authdomain

type UserRepository interface {
	GetUser() error
}
