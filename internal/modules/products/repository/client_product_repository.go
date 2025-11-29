package productsrepository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
)

func (s *ProductsStore) ClientServiceBalance(ctx context.Context, clientBalance *productsdomain.ClientServiceBalance) error {
	return s.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "service_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"balance": gorm.Expr("client_service_balances.balance + ?", clientBalance.Balance)}),
		}).
		Create(clientBalance).Error
}
