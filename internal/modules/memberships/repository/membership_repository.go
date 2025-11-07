package membershipsrepository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	membershipsdomain "github.com/vitalfit/api/internal/modules/memberships/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/pagination"
	"gorm.io/gorm"
)

type MembershipStore struct {
	db *gorm.DB
}

func NewMembershipStore(db *gorm.DB) *MembershipStore {
	return &MembershipStore{
		db: db,
	}
}

// CreateMembershipType crea un nuevo tipo de membresía en la base de datos.
func (s *MembershipStore) CreateMembershipType(ctx context.Context, membership *membershipsdomain.MembershipType) error {
	err := s.db.WithContext(ctx).Create(membership).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return shared_errors.ErrConflict
			}
		}
		return err
	}
	return nil
}

// CreateMembershipTypeTX permite crear un tipo de membresía dentro de una transacción activa.
func (s *MembershipStore) CreateMembershipTypeTX(ctx context.Context, tx *gorm.DB, membership *membershipsdomain.MembershipType) error {
	err := tx.WithContext(ctx).Create(membership).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return shared_errors.ErrConflict
			}
		}
		return err
	}
	return nil
}

// UpdateMembershipType actualiza los datos de un tipo de membresía existente.
func (s *MembershipStore) UpdateMembershipType(ctx context.Context, membership *membershipsdomain.MembershipType) error {
	err := s.db.WithContext(ctx).Save(membership).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return shared_errors.ErrNotFound
		default:
			return err
		}
	}
	return nil
}

// DeleteMembershipType realiza una eliminación lógica del tipo de membresía (soft delete).
func (s *MembershipStore) DeleteMembershipType(ctx context.Context, id uuid.UUID) error {
	var membership membershipsdomain.MembershipType
	result := s.db.WithContext(ctx).First(&membership, "membership_type_id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return shared_errors.ErrNotFound
	}

	now := time.Now()
	membership.IsActive = false
	membership.DeletedAt = &now

	if err := s.db.WithContext(ctx).Save(&membership).Error; err != nil {
		return err
	}
	return nil
}

// GetMembershipTypeByID obtiene un tipo de membresía por su UUID.
func (s *MembershipStore) GetMembershipTypeByID(ctx context.Context, id uuid.UUID) (*membershipsdomain.MembershipType, error) {
	membership := &membershipsdomain.MembershipType{}
	err := s.db.WithContext(ctx).
		Where("membership_type_id = ? AND deleted_at IS NULL", id).
		First(membership).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, shared_errors.ErrNotFound
		default:
			return nil, err
		}
	}
	return membership, nil
}

// GetMembershipTypes devuelve todos los tipos de membresía activos.
func (s *MembershipStore) GetMembershipTypes(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*membershipsdomain.MembershipType, error) {
	var memberships []*membershipsdomain.MembershipType

	baseQuery := s.db.WithContext(ctx).Model(&membershipsdomain.MembershipType{})

	if fq.Search != "" {
		baseQuery = baseQuery.Where("name ILIKE ?", "%"+fq.Search+"%")
	}

	query := baseQuery.
		Limit(fq.Limit).
		Offset(fq.Page*fq.Limit - fq.Limit).
		Order("created_at " + fq.Sort)

	if err := query.Find(&memberships).Error; err != nil {
		return nil, err
	}

	return memberships, nil
}

func (s *MembershipStore) GetMembershipTypesFTotal(ctx context.Context, fq pagination.PaginatedFeedQuery) (int64, error) {
	var count int64
	baseQuery := s.db.WithContext(ctx).Model(&membershipsdomain.MembershipType{})

	if fq.Search != "" {
		baseQuery = baseQuery.Where("name ILIKE ?", "%"+fq.Search+"%")
	}

	query := baseQuery.
		Order("created_at " + fq.Sort)

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *MembershipStore) GetSummary(ctx context.Context) (*membershipsdomain.MembershipSummary, error) {
	var summary membershipsdomain.MembershipSummary

	query := s.db.WithContext(ctx).Model(&membershipsdomain.MembershipType{})

	err := query.Select(
		"COUNT(*) as total",
		"COUNT(CASE WHEN is_active = true THEN 1 END) as actives",
		"COUNT(CASE WHEN is_active = false THEN 1 END) as inactives",
	).Take(&summary).Error

	if err != nil {
		return nil, err
	}

	return &summary, nil
}
