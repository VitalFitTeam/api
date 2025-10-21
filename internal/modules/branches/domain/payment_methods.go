package branchdomain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethods struct {
	MethodID     uuid.UUID             `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"method_id"`
	Name         string                `gorm:"type:varchar(255);not null;unique" json:"name"`
	Type         PaymentMethodTypeEnum `gorm:"type:payment_method_type_enum;not null" json:"type"`
	Description  string                `gorm:"type:text" json:"description,omitempty"`
	GlobalStatus bool                  `gorm:"type:bool;default:true" json:"global_status"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`

	BranchLinks []PaymentMethodsBranch `gorm:"foreignKey:MethodID" json:"branch_links,omitempty" swaggerignore:"true"`
}

type PaymentMethodsBranch struct {
	BranchID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_branch_method" json:"branch_id"`
	MethodID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_branch_method" json:"method_id"`
	IsActive  bool      `gorm:"type:bool;not null;default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Branch *Branch         `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Method *PaymentMethods `gorm:"foreignKey:MethodID" json:"method,omitempty"`
}
