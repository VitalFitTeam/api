package authdomain

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Password struct {
	text *string
	hash []byte
}

// Set hashs the given text and sets it to the Password struct
func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), 12)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash
	return nil
}

// compares users password with the given text
func (p *Password) Matches(text string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(text))
	if err == nil {
		return true, nil
	}
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return false, err
}

// Value implements driver.Valuer: indicates to GORM how to save the field in the DB.
func (p Password) Value() (driver.Value, error) {
	if len(p.hash) == 0 {
		return nil, nil
	}
	// GORM saves only basic types, so we return the hash as a byte slice
	return p.hash, nil
}

// Scan implements sql.Scanner: indicates to GORM how to read the field from the DB.
func (p *Password) Scan(value interface{}) error {
	if value == nil {
		p.hash = nil
		return nil
	}
	v, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unexpected scan type password: %T", value)
	}
	p.hash = v
	return nil
}

type UUIDModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Role corresponde a la tabla 'roles'
type Role struct {
	UUIDModel
	Name        string `gorm:"type:varchar(50);unique;not null" json:"name"`
	Level       int16  `gorm:"type:smallint;not null;default:0" json:"level"`
	Description string `gorm:"type:text" json:"description"`
}

// User corresponde a la tabla 'users'
type Users struct {
	UUIDModel
	FirstName        string `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName         string `gorm:"type:varchar(100);not null" json:"last_name"`
	Email            string `gorm:"type:citext;unique;not null" json:"email"`
	Phone            string `gorm:"type:varchar(50)" json:"phone"`
	IdentityDocument string `gorm:"type:varchar(50);unique" json:"identity_document"`

	// ⚡️ Usamos el struct Password personalizado aquí ⚡️
	PasswordHash Password `gorm:"column:password_hash;type:bytea;not null" json:"-"`

	ProfilePictureURL string `gorm:"type:varchar(255)" json:"profile_picture_url"`
	IsValidated       bool   `gorm:"default:false" json:"is_validated"`

	// Relación Belongs To: User pertenece a un Role
	RoleID uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`
	Role   Role      `gorm:"foreignKey:RoleID" json:"role"`
}
