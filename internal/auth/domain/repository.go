package authdomain

type UserRepository interface {
	Create()
	GetUser() error
}
