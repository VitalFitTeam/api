package instructordomain

import (
	"time"

	"github.com/google/uuid"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
	"gorm.io/gorm"
)

type Instructor struct {
	InstructorID uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"instructor_id"`
	UserID       uuid.UUID      `gorm:"type:uuid;unique;not null" json:"user_id"`
	Speciality   string         `gorm:"type:varchar(255)" json:"speciality"`
	Biography    string         `gorm:"type:text" json:"biography"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	User     *authdomain.Users     `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
	Branches []branchdomain.Branch `gorm:"many2many:branch_instructors;" json:"branches,omitempty"`
}

func (Instructor) TableName() string {
	return "instructors"
}

type BranchInstructor struct {
	BranchID     uuid.UUID `gorm:"primaryKey"`
	InstructorID uuid.UUID `gorm:"primaryKey"`
}

func (BranchInstructor) TableName() string {
	return "branch_instructors"
}
