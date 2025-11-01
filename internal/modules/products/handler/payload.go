package productshandler

import (
	"time"

	"github.com/google/uuid"
	marketinghandlers "github.com/vitalfit/api/internal/modules/marketing/handlers"
	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
)

type CreateServicePayload struct {
	Name          string                  `json:"name" binding:"required"`
	CategoryID    string                  `json:"category_id" binding:"required"`
	Description   string                  `json:"description" binding:"required"`
	Duration      int64                   `json:"duration" binding:"required,gte=1,lte=60"`
	Priority      int64                   `json:"priority" binding:"required,gte=1,lte=100"`
	IsFeatured    bool                    `json:"is_featured"`
	BannerID      string                  `json:"banner_id" binding:"required"`
	ServiceImages []ServiceImagesPayloads `json:"service_images"`
}

func (s *CreateServicePayload) ToService() (*productsdomain.Service, error) {
	categoryID, err := uuid.Parse(s.CategoryID)
	if err != nil {
		return nil, err
	}
	serviceImages := make([]productsdomain.ServiceImage, 0, len(s.ServiceImages))
	for _, image := range s.ServiceImages {
		image, err := image.ToServiceImage()
		if err != nil {
			return nil, err
		}
		serviceImages = append(serviceImages, *image)
	}

	service := &productsdomain.Service{
		Name:            s.Name,
		CategoryID:      categoryID,
		Description:     s.Description,
		DurationMinutes: s.Duration,
		PriorityScore:   s.Priority,
		IsFeatured:      s.IsFeatured,
		Images:          serviceImages,
	}
	return service, nil

}

type UpdateServicePayload struct {
	CreateServicePayload
}

type ServiceImagesPayloads struct {
	ImageURL     string `json:"image_url" binding:"required"`
	AltText      string `json:"alt_text"`
	DisplayOrder int    `json:"display_order"`
	IsPrimary    bool   `json:"is_primary"`
}

func (s *ServiceImagesPayloads) ToServiceImage() (*productsdomain.ServiceImage, error) {
	serviceImage := &productsdomain.ServiceImage{
		ImageURL:     s.ImageURL,
		AltText:      s.AltText,
		DisplayOrder: s.DisplayOrder,
		IsPrimary:    s.IsPrimary,
	}
	return serviceImage, nil
}

type ServiceResponse struct {
	ServiceID       uuid.UUID                          `json:"service_id"`
	CategoryID      uuid.UUID                          `json:"category_id"`
	Name            string                             `json:"name"`
	Description     string                             `json:"description"`
	DurationMinutes int64                              `json:"duration_minutes"`
	PriorityScore   int64                              `json:"priority_score"`
	IsFeatured      bool                               `json:"is_featured"`
	CreatedAt       time.Time                          `json:"created_at"`
	UpdatedAt       time.Time                          `json:"updated_at"`
	ServiceCategory ServiceCategoryResponse            `json:"service_category"`
	Images          []ImagesRensponse                  `json:"images"`
	Banners         []marketinghandlers.BannerResponse `json:"banners"`
}

type ImagesRensponse struct {
	ImageID      uuid.UUID `json:"image_id"`
	ImageURL     string    `json:"image_url"`
	AltText      string    `json:"alt_text"`
	DisplayOrder int       `json:"display_order"`
	IsPrimary    bool      `json:"is_primary"`
}

type ServiceCategoryResponse struct {
	CategoryID uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
}
