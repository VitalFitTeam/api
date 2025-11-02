package billinghandlers

import (
	"encoding/json"
	"errors"

	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
)

type CreatePaymentMethodPayload struct {
	Name           string `json:"name" binding:"required"`
	Type           string `json:"type" binding:"required"`
	Description    string `json:"description"`
	ProcessingType string `json:"processing_type" binding:"required"`
}

func (p *CreatePaymentMethodPayload) toPaymentMethod() *billingdomain.PaymentMethods {
	return &billingdomain.PaymentMethods{
		Name:           p.Name,
		Type:           billingdomain.PaymentMethodTypeEnum(p.Type),
		Description:    p.Description,
		ProcessingType: billingdomain.PaymentProcessingTypeEnum(p.ProcessingType),
	}
}

type UpdatePaymentMethodPayload struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Description    string `json:"description"`
	ProcessingType string `json:"processing_type" binding:"required"`
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
	MethodID            string          `json:"method_id" binding:"required"`
	DisplayName         string          `json:"display_name"`
	Configuration       json.RawMessage `json:"configuration,omitempty" swaggertype:"object"`
	Visibility          string          `json:"visibility"`
	SurchargeFixed      int64           `json:"surcharge_fixed"`
	SurchargePercentage float64         `json:"surcharge_percentage"`
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
	IsActive            string          `json:"is_active,omitempty"`
	DisplayName         string          `json:"display_name,omitempty"`
	Configuration       json.RawMessage `json:"configuration,omitempty" swaggertype:"object"`
	Visibility          string          `json:"visibility,omitempty" binding:"omitempty,oneof=Client Staff All"`
	SurchargeFixed      int64           `json:"surcharge_fixed,omitempty" binding:"omitempty,min=0"`
	SurchargePercentage float64         `json:"surcharge_percentage,omitempty" binding:"omitempty,min=0"`
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
