package productshandler

type CreateServicePayload struct {
	Name        string `json:"name" binding:"required"`
	CategoryID  string `json:"category_id" binding:"required"`
	Description string `json:"description" binding:"required"`
	Duration    int64  `json:"duration" binding:"required"`
	Priority    int64  `json:"priority" binding:"required"`
	IsFeatured  bool   `json:"is_featured" binding:"required"`
	BannerID    string `json:"banner_id" binding:"required"`
}
