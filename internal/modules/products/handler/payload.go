package productshandler

import (
	"time"

	"github.com/google/uuid"
)

type InstructorResponse struct {
	InstructorID      uuid.UUID `json:"instructor_id"`
	UserID            uuid.UUID `json:"user_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	IdentityDocument  string    `json:"identity_document"`
	BirthDate         time.Time `json:"birth_date"`
	Gender            string    `json:"gender"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	Biography         string    `json:"biography"`
}
type CreateServicePayload struct {
	Name        string `json:"name" binding:"required"`
	CategoryID  string `json:"category_id" binding:"required"`
	Description string `json:"description" binding:"required"`
	Duration    int64  `json:"duration" binding:"required"`
}
