package billingrepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
)

type Billingstore struct {
	db *gorm.DB
}

func NewBillingStore(db *gorm.DB) *Billingstore {
	return &Billingstore{
		db: db,
	}
}

func (bs *Billingstore) CreateInvoice(ctx context.Context, invoice *billingdomain.Invoice) error {
	return db.WithTX(bs.db, func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(invoice).Error
	})
}

func (bs *Billingstore) GetInvoiceByID(ctx context.Context, invoiceID uuid.UUID) (*billingdomain.Invoice, error) {
	var invoice billingdomain.Invoice
	err := bs.db.WithContext(ctx).Preload("InvoiceItems").Preload("Payments").First(&invoice, "invoice_id = ?", invoiceID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.ErrNotFound
		}
		return nil, err
	}
	return &invoice, nil
}

func (bs *Billingstore) AddPaymentToInvoice(ctx context.Context, payment *billingdomain.Payment) error {
	return db.WithTX(bs.db, func(tx *gorm.DB) error {
		return tx.WithContext(ctx).Create(payment).Error
	})

}
