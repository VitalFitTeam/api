package billingrepository

import (
	"context"

	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
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
