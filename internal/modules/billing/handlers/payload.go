package billinghandlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
)

type CreatePaymentMethodPayload struct {
	Name                string          `json:"name" binding:"required"`
	Type                string          `json:"type" binding:"required"`
	Description         string          `json:"description"`
	ProcessingType      string          `json:"processing_type" binding:"required"`
	DisplayName         string          `json:"display_name"`
	Configuration       json.RawMessage `json:"configuration,omitempty" swaggertype:"object"`
	Visibility          string          `json:"visibility"`
	SurchargeFixed      int64           `json:"surcharge_fixed"`
	SurchargePercentage float64         `json:"surcharge_percentage"`
}

func (p *CreatePaymentMethodPayload) toPaymentMethod() *billingdomain.PaymentMethods {
	return &billingdomain.PaymentMethods{
		Name:                p.Name,
		Type:                billingdomain.PaymentMethodTypeEnum(p.Type),
		Description:         p.Description,
		ProcessingType:      billingdomain.PaymentProcessingTypeEnum(p.ProcessingType),
		DisplayName:         p.DisplayName,
		Visibility:          billingdomain.BranchPaymentVisibilityEnum(p.Visibility),
		SurchargeFixed:      p.SurchargeFixed,
		SurchargePercentage: p.SurchargePercentage,
	}
}

type UpdatePaymentMethodPayload struct {
	Name                string          `json:"name"`
	Type                string          `json:"type"`
	Description         string          `json:"description"`
	ProcessingType      string          `json:"processing_type" binding:"required"`
	DisplayName         string          `json:"display_name"`
	Configuration       json.RawMessage `json:"configuration,omitempty" swaggertype:"object"`
	Visibility          string          `json:"visibility"`
	SurchargeFixed      int64           `json:"surcharge_fixed"`
	SurchargePercentage float64         `json:"surcharge_percentage"`
}

func (p *UpdatePaymentMethodPayload) toPaymentMethod() *billingdomain.PaymentMethods {
	return &billingdomain.PaymentMethods{
		Name:           p.Name,
		Type:           billingdomain.PaymentMethodTypeEnum(p.Type),
		Description:    p.Description,
		ProcessingType: billingdomain.PaymentProcessingTypeEnum(p.ProcessingType),
	}
}

type MethodBranchConfigPayload struct {
	MethodID []string `json:"method_id" binding:"required"`
}

type ZelleConfig struct {
	Email string `json:"email" binding:"required,email"`
}
type BankTransferConfig struct {
	BankName      string `json:"bank_name" binding:"required"`
	AccountNumber string `json:"account_number" binding:"required"`
	TaxID         string `json:"tax_id" binding:"required"`
}
type PagoMovilConfig struct {
	Phone  string `json:"phone" binding:"required"`
	BankID string `json:"bank_id" binding:"required"`
	TaxID  string `json:"tax_id" binding:"required"`
}

type UpdateBranchConfigPayload struct {
	IsActive string `json:"is_active,omitempty"`
}

type BranchPaymentMethodResponse struct {
	BranchID uuid.UUID `json:"branch_id" binding:"required"`
	MethodID uuid.UUID `json:"method_id" binding:"required"`
	Name     string    `json:"name" binding:"required"`
	Type     string    `json:"type" binding:"required"`
	IsActive bool      `json:"is_active" binding:"required"`
}

func (h *BillingHandlers) validateBranchConfig(payloadConfig json.RawMessage, method *billingdomain.PaymentMethods) error {

	switch method.Type {

	case billingdomain.PaymentMethodTransfer:

		switch method.Name {
		case "Zelle":
			var zelleConfig ZelleConfig
			if err := json.Unmarshal(payloadConfig, &zelleConfig); err != nil {
				return errors.New("invalid configuration json for Zelle")
			}
			if err := validate.Struct(zelleConfig); err != nil {
				return err
			}

		case "Bank Transfer":
			var bankTransferConfig BankTransferConfig
			if err := json.Unmarshal(payloadConfig, &bankTransferConfig); err != nil {
				return errors.New("invalid configuration json for Bank Transfer")
			}
			if err := validate.Struct(bankTransferConfig); err != nil {
				return err
			}

		case "Pago Movil":
			var pagoMovilConfig PagoMovilConfig
			if err := json.Unmarshal(payloadConfig, &pagoMovilConfig); err != nil {
				return errors.New("invalid configuration json for Pago Móvil")
			}
			if err := validate.Struct(pagoMovilConfig); err != nil {
				return err
			}

		default:
			return errors.New("configuration not supported for this transfer method")
		}

	case billingdomain.PaymentMethodCard:
		if string(payloadConfig) != "{}" && string(payloadConfig) != "" {
			return errors.New("configuration must be empty for Card methods")
		}

	case billingdomain.PaymentMethodCash:
		if string(payloadConfig) != "{}" && string(payloadConfig) != "" {
			return errors.New("configuration must be empty for Cash")
		}

	case billingdomain.PaymentMethodOther:
		if string(payloadConfig) != "{}" && string(payloadConfig) != "" {
			return errors.New("configuration not supported for 'Other' type")
		}

	default:
		return errors.New("internal validation error: unhandled payment type")
	}

	return nil
}

type CreateFiscalDocumentTypePayload struct {
	Name   string `json:"name" binding:"required"`
	Prefix string `json:"prefix" binding:"required"`
}

func (p *CreateFiscalDocumentTypePayload) ToFiscalDocumentType() *billingdomain.FiscalDocumentType {
	return &billingdomain.FiscalDocumentType{
		Name:   p.Name,
		Prefix: p.Prefix,
	}
}

type UpdateFiscalDocumentTypePayload struct {
	Name   string `json:"name,omitempty"`
	Prefix string `json:"prefix,omitempty"`
}

func (p *UpdateFiscalDocumentTypePayload) ToFiscalDocumentType() *billingdomain.FiscalDocumentType {
	return &billingdomain.FiscalDocumentType{
		Name:   p.Name,
		Prefix: p.Prefix,
	}
}

type FiscalDocumentTypeResponse struct {
	DocumentTypeID uuid.UUID `json:"document_type_id"`
	Name           string    `json:"name"`
	Prefix         string    `json:"prefix"`
	CreatedAt      time.Time `json:"created_at"`
}

func NewFiscalDocumentTypeResponse(docType *billingdomain.FiscalDocumentType) *FiscalDocumentTypeResponse {
	return &FiscalDocumentTypeResponse{
		DocumentTypeID: docType.DocumentTypeID,
		Name:           docType.Name,
		Prefix:         docType.Prefix,
		CreatedAt:      docType.CreatedAt,
	}
}

//orders payload

type CreateInvoicePayload struct {
	UserID   *uuid.UUID `json:"user_id"`
	BranchID *uuid.UUID `json:"branch_id" binding:"required"`

	Items []ItemDTO `json:"items" binding:"required,min=1,dive"`
}

type ItemDTO struct {
	ItemID   uuid.UUID `json:"item_id" binding:"required"`
	ItemType string    `json:"item_type" binding:"required,oneof=membership package service"`
	Quantity int       `json:"quantity" binding:"required,min=1"`
}

func (i *CreateInvoicePayload) ToInvoice() *billingdomain.Invoice {
	invoice := &billingdomain.Invoice{
		BranchID:       *i.BranchID,
		IssueDate:      time.Now(),
		DueDate:        time.Now().AddDate(0, 0, 30),
		TotalAmount:    decimal.Zero,
		Tax:            decimal.Zero,
		Status:         billingdomain.InvoiceStatusUnpaid,
		DocumentTypeID: uuid.Nil,
	}
	return invoice
}

func (i *CreateInvoicePayload) ToInvoiceItems() []billingdomain.InvoiceItem {
	invoiceItems := make([]billingdomain.InvoiceItem, len(i.Items))
	for idx, item := range i.Items {
		invoiceItem := billingdomain.InvoiceItem{
			Quantity: item.Quantity,
		}
		switch item.ItemType {
		case "membership":
			invoiceItem.MembershipTypeID = uuid.NullUUID{UUID: item.ItemID, Valid: true}
		case "package":
			invoiceItem.PackageID = uuid.NullUUID{UUID: item.ItemID, Valid: true}
		case "service":
			invoiceItem.ServiceID = uuid.NullUUID{UUID: item.ItemID, Valid: true}
		}
		invoiceItems[idx] = invoiceItem
	}
	return invoiceItems
}

type CreatePaymentPayload struct {
	InvoiceID       uuid.UUID       `json:"invoice_id" binding:"required"`
	PaymentMethodID uuid.UUID       `json:"payment_method_id" binding:"required"`
	AmountPaid      decimal.Decimal `json:"amount_paid" binding:"required"`
	CurrencyPaid    string          `json:"currency_paid" binding:"required,len=3"`
	TransactionID   string          `json:"transaction_id"`
	ReceiptURL      string          `json:"receipt_url"`
}

func (p *CreatePaymentPayload) ToPayment() *billingdomain.Payment {
	return &billingdomain.Payment{
		InvoiceID:       p.InvoiceID,
		PaymentMethodID: p.PaymentMethodID,
		AmountPaid:      p.AmountPaid,
		CurrencyPaid:    p.CurrencyPaid,
		TransactionID:   sql.NullString{String: p.TransactionID, Valid: p.TransactionID != ""},
		ReceiptURL:      sql.NullString{String: p.ReceiptURL, Valid: p.ReceiptURL != ""},
		Status:          billingdomain.PaymentStatusPending,
	}
}
