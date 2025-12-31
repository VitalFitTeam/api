package wishlistdomain

import (
	"context"

	"github.com/google/uuid"
)

type WishlistRepository interface {
	// AddToWishlist adds a service to the user's wishlist
	AddToWishlist(ctx context.Context, wishlist *Wishlist) (uuid.UUID, error)

	// RemoveFromWishlist removes a service from the user's wishlist
	RemoveFromWishlist(ctx context.Context, wishlistID uuid.UUID) error

	// GetUserWishlist retrieves all wishlist items for a user with service details
	GetUserWishlist(ctx context.Context, userID uuid.UUID) ([]*WishlistWithService, error)

	// WishlistExists checks if a user already has a service in their wishlist
	WishlistExists(ctx context.Context, userID, serviceID uuid.UUID) (bool, error)
}
