package billinghandlers

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
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

// CreateFiscalDocumentTypePayload defines the payload for creating a new fiscal document type.
type CreateFiscalDocumentTypePayload struct {
	Name   string `json:"name" binding:"required"`
	Prefix string `json:"prefix" binding:"required"`
}

// ToFiscalDocumentType converts the payload to a domain model.
func (p *CreateFiscalDocumentTypePayload) ToFiscalDocumentType() *billingdomain.FiscalDocumentType {
	return &billingdomain.FiscalDocumentType{
		Name:   p.Name,
		Prefix: p.Prefix,
	}
}

// UpdateFiscalDocumentTypePayload defines the payload for updating a fiscal document type.
type UpdateFiscalDocumentTypePayload struct {
	Name   string `json:"name,omitempty"`
	Prefix string `json:"prefix,omitempty"`
}

// ToFiscalDocumentType converts the payload to a domain model.
func (p *UpdateFiscalDocumentTypePayload) ToFiscalDocumentType() *billingdomain.FiscalDocumentType {
	return &billingdomain.FiscalDocumentType{
		Name:   p.Name,
		Prefix: p.Prefix,
	}
}

// FiscalDocumentTypeResponse defines the response for a fiscal document type.
type FiscalDocumentTypeResponse struct {
	DocumentTypeID uuid.UUID `json:"document_type_id"`
	Name           string    `json:"name"`
	Prefix         string    `json:"prefix"`
	CreatedAt      time.Time `json:"created_at"`
}

// NewFiscalDocumentTypeResponse creates a new response from the domain model.
func NewFiscalDocumentTypeResponse(docType *billingdomain.FiscalDocumentType) *FiscalDocumentTypeResponse {
	return &FiscalDocumentTypeResponse{
		DocumentTypeID: docType.DocumentTypeID,
		Name:           docType.Name,
		Prefix:         docType.Prefix,
		CreatedAt:      docType.CreatedAt,
	}
}
