package billinghandlers

import billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"

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
