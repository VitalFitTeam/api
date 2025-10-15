package authhandlers

import (
	"time"

	authdomain "github.com/vitalfit/api/internal/auth/domain"
)

type CreateUserClientPayload struct {
	FirstName         string `json:"first_name" binding:"required"`
	LastName          string `json:"last_name" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	Phone             string `json:"phone" binding:"required"`
	IdentityDocument  string `json:"identity_document" binding:"required"`
	Password          string `json:"password" binding:"required,min=8"`
	BirthDate         string `json:"birth_date" binding:"required"`
	Gender            string `json:"gender" binding:"required"`
	ProfilePictureURL string `json:"profile_picture_url"`
}

type CreateUserStaffPayload struct {
	*CreateUserClientPayload
	RoleName string `json:"role_name" binding:"omitempty"`
}

func (c *CreateUserClientPayload) createUser() (*authdomain.Users, error) {
	birthdate, err := time.Parse("2006-01-02", c.BirthDate)
	if err != nil {
		birthdate, err = time.Parse(time.RFC3339, c.BirthDate)
		if err != nil {
			return nil, err
		}
	}
	return &authdomain.Users{
		FirstName:         c.FirstName,
		LastName:          c.LastName,
		Email:             c.Email,
		Phone:             c.Phone,
		IdentityDocument:  c.IdentityDocument,
		BirthDate:         birthdate,
		Gender:            authdomain.GenderEnum(c.Gender),
		ProfilePictureURL: c.ProfilePictureURL,
	}, nil
}

type CodePayload struct {
	Code string `json:"code" binding:"required"`
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=3,max=72"`
}

type ForgotPasswordPayload struct {
	Email string `json:"email" binding:"required,email,max=255"`
}

type ResetPasswordPayload struct {
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"` // Valida en el backend
	Token           string `json:"token" binding:"required"`
}
