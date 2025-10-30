package inventoryrepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	inventorydomain "github.com/vitalfit/api/internal/modules/inventory/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
)

type EquipmentStore struct {
	db *gorm.DB
}

func NewEquipmentStore(db *gorm.DB) *EquipmentStore {
	return &EquipmentStore{db: db}
}

func (s *EquipmentStore) Create(ctx context.Context, equipment *inventorydomain.Equipment) (*inventorydomain.Equipment, error) {
	err := db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(&equipment).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				if pgErr.ConstraintName == "idx_equipment_name_brand_model" {
					return shared_errors.ErrConflict
				}
			}
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return equipment, nil
}

func (s *EquipmentStore) GetAll(ctx context.Context) (*inventorydomain.EquipmentQueryResults, error) {
	var equipments []*inventorydomain.Equipment
	err := s.db.WithContext(ctx).Order("created_at desc").Find(&equipments).Error
	if err != nil {
		return nil, err
	}
	return &inventorydomain.EquipmentQueryResults{Equipments: equipments}, nil
}

func (s *EquipmentStore) GetByID(ctx context.Context, equipmentID uuid.UUID) (*inventorydomain.Equipment, error) {
	var equipment inventorydomain.Equipment
	err := s.db.WithContext(ctx).First(&equipment, "equipment_id = ?", equipmentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.ErrNotFound
		}
		return nil, err
	}
	return &equipment, nil
}

func (s *EquipmentStore) Update(ctx context.Context, equipment *inventorydomain.Equipment) error {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	return db.WithTX(s.db, func(tx *gorm.DB) error {
		result := tx.WithContext(ctx).
			Model(&inventorydomain.Equipment{EquipmentID: equipment.EquipmentID}).
			Updates(equipment)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return shared_errors.ErrNotFound
		}
		return nil
	})
}

func (s *EquipmentStore) Delete(ctx context.Context, equipmentID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()

	return db.WithTX(s.db, func(tx *gorm.DB) error {
		result := tx.WithContext(ctx).
			Delete(&inventorydomain.Equipment{}, "equipment_id = ?", equipmentID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return shared_errors.ErrNotFound
		}
		return nil
	})
}
