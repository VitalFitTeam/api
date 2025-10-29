package instructorhandler

import (
	"time"

	"github.com/google/uuid"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	instructordomain "github.com/vitalfit/api/internal/modules/instructor/domain"
)

type CreateInstructorPayload struct {
	FirstName         string `json:"first_name" binding:"required"`
	LastName          string `json:"last_name" binding:"required"`
	Email             string `json:"email" binding:"required,email"`
	Phone             string `json:"phone" binding:"required"`
	IdentityDocument  string `json:"identity_document" binding:"required"`
	Password          string `json:"password" binding:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789,containsany=!@#$%^&*"`
	BirthDate         string `json:"birth_date" binding:"required"`
	Gender            string `json:"gender" binding:"required"`
	ProfilePictureURL string `json:"profile_picture_url"`
	Speciality        string `json:"speciality"`
	Biography         string `json:"biography"`
}

func (c *CreateInstructorPayload) toInstructor() (*instructordomain.Instructor, error) {
	birthdate, err := time.Parse("2006-01-02", c.BirthDate)
	if err != nil {
		birthdate, err = time.Parse(time.RFC3339, c.BirthDate)
		if err != nil {
			return nil, err
		}
	}
	user := &authdomain.Users{
		FirstName:         c.FirstName,
		LastName:          c.LastName,
		Email:             c.Email,
		Phone:             c.Phone,
		IdentityDocument:  c.IdentityDocument,
		BirthDate:         birthdate,
		Gender:            authdomain.GenderEnum(c.Gender),
		ProfilePictureURL: c.ProfilePictureURL,
	}

	return &instructordomain.Instructor{
		Speciality: c.Speciality,
		Biography:  c.Biography,
		User:       *user,
	}, nil

}

type UpdateInstructorPayload struct {
	Speciality string `json:"speciality"`
	Biography  string `json:"biography"`
}

func (u *UpdateInstructorPayload) toInstructor(instructorID uuid.UUID) *instructordomain.Instructor {
	return &instructordomain.Instructor{
		InstructorID: instructorID,
		Speciality:   u.Speciality,
		Biography:    u.Biography,
	}
}
