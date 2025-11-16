package billingdomain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	combosdomain "github.com/vitalfit/api/internal/modules/combos/domain"
	membershipsdomain "github.com/vitalfit/api/internal/modules/memberships/domain"
	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
	"gorm.io/gorm"
)

type InvoiceStatus string

const (
	InvoiceStatusPaid    InvoiceStatus = "Paid"
	InvoiceStatusUnpaid  InvoiceStatus = "Unpaid"
	InvoiceStatusVoid    InvoiceStatus = "Void"
	InvoiceStatusOverdue InvoiceStatus = "Overdue"
)

type PaymentStatus string

const (
	PaymentStatusCompleted PaymentStatus = "Completed"
	PaymentStatusFailed    PaymentStatus = "Failed"
	PaymentStatusRefunded  PaymentStatus = "Refunded"
	PaymentStatusPending   PaymentStatus = "Pending"
)

type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "Pending"
	RefundStatusProcessed RefundStatus = "Processed"
	RefundStatusFailed    RefundStatus = "Failed"
)

type Invoice struct {
	InvoiceID      uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"invoice_id"`
	UserID         uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	DocumentTypeID uuid.UUID       `gorm:"type:uuid;not null" json:"document_type_id"`
	InvoiceNumber  string          `gorm:"type:varchar(50);not null;unique" json:"invoice_number"`
	IssueDate      time.Time       `gorm:"type:date;not null" json:"issue_date"`
	DueDate        time.Time       `gorm:"type:date;not null" json:"due_date"`
	TotalAmount    decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Tax            decimal.Decimal `gorm:"type:decimal(10,2)" json:"tax"`
	Status         InvoiceStatus   `gorm:"type:invoice_status;not null;index" json:"status"`
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`

	User         authdomain.Users   `gorm:"foreignKey:UserID"`
	DocumentType FiscalDocumentType `gorm:"foreignKey:DocumentTypeID"`

	InvoiceItems []InvoiceItem `gorm:"foreignKey:InvoiceID" json:"invoice_items,omitempty"`
	Payments     []Payment     `gorm:"foreignKey:InvoiceID" json:"payments,omitempty"`
}

func (Invoice) TableName() string {
	return "invoices"
}

type InvoiceItem struct {
	InvoiceItemID    uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"invoice_item_id"`
	InvoiceID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"invoice_id"`
	Quantity         int             `gorm:"not null;default:1" json:"quantity"`
	UnitPrice        decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	PromotionID      uuid.NullUUID   `gorm:"type:uuid" json:"promotion_id"`
	DiscountApplied  decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"discount_applied"`
	MembershipTypeID uuid.NullUUID   `gorm:"type:uuid" json:"membership_type_id"`
	ServiceID        uuid.NullUUID   `gorm:"type:uuid" json:"service_id"`
	PackageID        uuid.NullUUID   `gorm:"type:uuid" json:"package_id"`

	Invoice Invoice `gorm:"foreignKey:InvoiceID"`
	//Promotion      marketingdomain.Promotion       `gorm:"foreignKey:PromotionID"`
	MembershipType membershipsdomain.MembershipType `gorm:"foreignKey:MembershipTypeID"`
	Service        productsdomain.Service           `gorm:"foreignKey:ServiceID"`
	Package        combosdomain.Package             `gorm:"foreignKey:PackageID"`
}

func (InvoiceItem) TableName() string {
	return "invoice_items"
}

type Payment struct {
	PaymentID   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"payment_id"`
	InvoiceID   uuid.UUID `gorm:"type:uuid;not null;index" json:"invoice_id"`
	PaymentDate time.Time `gorm:"type:timestamptz;not null" json:"payment_date"`

	AmountPaid   decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"amount_paid"`
	CurrencyPaid string          `gorm:"type:varchar(3);not null" json:"currency_paid"`

	AmountBase   decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"amount_base"`
	CurrencyBase string          `gorm:"type:varchar(3);not null;default:'USD'" json:"currency_base"`
	ExchangeRate decimal.Decimal `gorm:"type:decimal(18,8);not null" json:"exchange_rate"`

	PaymentMethodID uuid.UUID      `gorm:"type:uuid;not null" json:"payment_method_id"`
	TransactionID   sql.NullString `gorm:"type:varchar(255)" json:"transaction_id"`
	ReceiptURL      sql.NullString `gorm:"type:varchar(255)" json:"receipt_url"`
	Status          PaymentStatus  `gorm:"type:payment_status;not null;index" json:"status"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Invoice       Invoice        `gorm:"foreignKey:InvoiceID"`
	PaymentMethod PaymentMethods `gorm:"foreignKey:MethodID"`

	Refunds []Refund `gorm:"foreignKey:PaymentID" json:"refunds,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

type Refund struct {
	RefundID  uuid.UUID     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"refund_id"`
	PaymentID uuid.UUID     `gorm:"type:uuid;not null;index" json:"payment_id"`
	InvoiceID uuid.NullUUID `gorm:"type:uuid;index" json:"invoice_id"`
	Reason    string        `gorm:"type:varchar(255);not null" json:"reason"`

	AmountRefunded   decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"amount_refunded"`
	CurrencyRefunded string          `gorm:"type:varchar(3);not null" json:"currency_refunded"`

	AmountBase   decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"amount_base"`
	CurrencyBase string          `gorm:"type:varchar(3);not null;default:'USD'" json:"currency_base"`
	ExchangeRate decimal.Decimal `gorm:"type:decimal(18,8);not null" json:"exchange_rate"`

	RefundMethodID uuid.NullUUID  `gorm:"type:uuid" json:"refund_method_id"`
	TransactionID  sql.NullString `gorm:"type:varchar(255)" json:"transaction_id"`
	Status         RefundStatus   `gorm:"type:refund_status;not null;default:'Pending';index" json:"status"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	ProcessedAt    sql.NullTime   `json:"processed_at"`

	Payment      Payment        `gorm:"foreignKey:PaymentID"`
	Invoice      Invoice        `gorm:"foreignKey:InvoiceID"`
	RefundMethod PaymentMethods `gorm:"foreignKey:MethodID"`
}

func (Refund) TableName() string {
	return "refunds"
}

//REDIS

const (
	LatestRatesKey        = "latest_rates"
	RatePrefix            = "rate"
	RateExpiration        = 3 * time.Minute
	LatestRatesExpiration = 2 * time.Hour
)

type ExchangeRates struct {
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Timestamp  int64              `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}
