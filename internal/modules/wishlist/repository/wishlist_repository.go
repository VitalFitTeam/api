package wishlistrepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	wishlistdomain "github.com/vitalfit/api/internal/modules/wishlist/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"gorm.io/gorm"
)

type WishlistStore struct {
	db *gorm.DB
}

func NewWishlistStore(db *gorm.DB) *WishlistStore {
	return &WishlistStore{db: db}
}

func (s *WishlistStore) AddToWishlist(ctx context.Context, wishlist *wishlistdomain.Wishlist) (uuid.UUID, error) {
	if err := s.db.WithContext(ctx).Create(wishlist).Error; err != nil {
		// Check for unique constraint violation
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return uuid.Nil, shared_errors.ErrConflict
		}
		return uuid.Nil, err
	}

	return wishlist.WishlistID, nil
}

func (s *WishlistStore) RemoveFromWishlist(ctx context.Context, wishlistID uuid.UUID) error {
	result := s.db.WithContext(ctx).
		Where("wishlist_id = ?", wishlistID).
		Delete(&wishlistdomain.Wishlist{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return shared_errors.ErrNotFound
	}

	return nil
}

func (s *WishlistStore) GetUserWishlist(ctx context.Context, userID uuid.UUID) ([]*wishlistdomain.WishlistWithService, error) {
	var wishlist []*wishlistdomain.WishlistWithService

	err := s.db.WithContext(ctx).
		Table("wishlist").
		Select("wishlist.wishlist_id, wishlist.user_id, wishlist.service_id, services.name as service_name, services.description, services.image_url, wishlist.created_at").
		Joins("JOIN services ON wishlist.service_id = services.service_id").
		Where("wishlist.user_id = ?", userID).
		Order("wishlist.created_at DESC").
		Scan(&wishlist).Error

	if err != nil {
		return nil, err
	}

	return wishlist, nil
}

func (s *WishlistStore) WishlistExists(ctx context.Context, userID, serviceID uuid.UUID) (bool, error) {
	var count int64

	err := s.db.WithContext(ctx).
		Model(&wishlistdomain.Wishlist{}).
		Where("user_id = ? AND service_id = ?", userID, serviceID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
