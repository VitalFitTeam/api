package billingservice

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/vitalfit/api/config"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/internal/store/cache"
)

type BillingService struct {
	store store.Storage
	cache cache.Storage
	cfg   config.Config
	http  *http.Client
}

func NewBillingService(store store.Storage, cache cache.Storage, cfg config.Config) *BillingService {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{
		Transport: customTransport,
		Timeout:   10 * time.Second,
	}

	return &BillingService{
		store: store,
		cache: cache,
		cfg:   cfg,
		http:  httpClient,
	}
}
func (bs *BillingService) CreateInvoice(ctx context.Context, invoice *billingdomain.Invoice, items []billingdomain.InvoiceItem) error {
	docType, err := bs.store.FiscalDocuments.GetFiscalDocumentTypeByName(ctx, "Invoice")
	if err != nil {
		return err
	}

	branch, err := bs.store.Branches.GetByID(ctx, invoice.BranchID)
	if err != nil {
		return err
	}

	invoice.DocumentTypeID = docType.DocumentTypeID
	invoice.InvoiceID = uuid.New()
	invoice.InvoiceNumber = docType.Prefix + time.Now().Format("060102") + "-" + invoice.InvoiceID.String()[:8]

	taxRate := billingdomain.GetTaxRateByLocation(branch.State.Country.Name)

	for i := range items {
		items[i].TaxRate = taxRate
		item := &items[i]

		if item.MembershipTypeID.Valid {
			membershipType, err := bs.store.Membership.GetMembershipTypeByID(ctx, item.MembershipTypeID.UUID)
			if err != nil {
				return err
			}
			item.UnitPrice = decimal.NewFromFloat(membershipType.Price)
		} else if item.PackageID.Valid {
			pkg, err := bs.store.Combos.GetPackageByID(ctx, item.PackageID.UUID)
			if err != nil {
				return err
			}
			item.UnitPrice = decimal.NewFromFloat(pkg.Price)
		} else if item.ServiceID.Valid {
			serviceDetail, err := bs.store.Products.GetBranchServiceByID(ctx, invoice.BranchID, item.ServiceID.UUID)
			if err != nil {
				return err
			}
			item.UnitPrice = decimal.NewFromFloat(serviceDetail.PriceForNonMember)
		}
	}

	invoice.InvoiceItems = items

	invoice.CalculateTotals()

	err = bs.store.Billing.CreateInvoice(ctx, invoice)
	if err != nil {
		return err
	}

	return nil
}
