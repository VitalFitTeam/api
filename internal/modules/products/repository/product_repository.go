package productsrepository

<<<<<<< HEAD
import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	marketingdomain "github.com/vitalfit/api/internal/modules/marketing/domain"
	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)
=======
import "gorm.io/gorm"
>>>>>>> dev

type ProductsStore struct {
	db *gorm.DB
}

func NewProductsStore(db *gorm.DB) *ProductsStore {
	return &ProductsStore{
		db: db,
	}
}
<<<<<<< HEAD

func (s *ProductsStore) CreateServiceCategory(ctx context.Context, tx *gorm.DB, serviceCategory *productsdomain.ServiceCategory) error {
	if err := tx.Create(serviceCategory).Error; err != nil {
		return err
	}
	return nil
}

func (s *ProductsStore) ListServiceCategories(ctx context.Context, tx *gorm.DB) ([]productsdomain.ServiceCategory, error) {
	var serviceCategories []productsdomain.ServiceCategory
	if err := tx.Find(&serviceCategories).Error; err != nil {
		return nil, err
	}
	return serviceCategories, nil
}

func (s *ProductsStore) GetServicesCategories(ctx context.Context) ([]productsdomain.ServiceCategory, error) {
	var serviceCategories []productsdomain.ServiceCategory

	err := db.WithTX(s.db, func(tx *gorm.DB) error {
		var listErr error

		serviceCategories, listErr = s.ListServiceCategories(ctx, tx)
		if listErr != nil {
			return listErr
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return serviceCategories, nil
}
func (s *ProductsStore) CreateService(ctx context.Context, service *productsdomain.Service, bannerID uuid.UUID) error {

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(service).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" {
					return shared_errors.ErrConflict
				}
			}
			return err
		}

		if bannerID == uuid.Nil {
			return nil
		}

		bannerlink := &marketingdomain.BannerService{
			ServiceID: service.ServiceID,
			BannerID:  bannerID,
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(bannerlink).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *ProductsStore) CreateServiceTX(ctx context.Context, tx *gorm.DB, service *productsdomain.Service, bannerID uuid.UUID) error {

	if err := tx.Create(service).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return shared_errors.ErrConflict
			}
		}
		return err
	}

	if bannerID == uuid.Nil {
		return nil
	}

	bannerlink := &marketingdomain.BannerService{
		ServiceID: service.ServiceID,
		BannerID:  bannerID,
	}

	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(bannerlink).Error; err != nil {
		return err
	}

	return nil

}

func (s *ProductsStore) GetServices(ctx context.Context) ([]productsdomain.Service, error) {
	var services []productsdomain.Service
	if err := s.db.Preload("Images").Preload("Category").Preload("Banners").Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil

}
func (s *ProductsStore) DeleteService(ctx context.Context, serviceID uuid.UUID) error {

	err := db.WithTX(s.db, func(tx *gorm.DB) error {
		if err := tx.Delete(&productsdomain.Service{}, serviceID).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				return err
			default:
				return err
			}
		}
		return nil
	})
	return err

}
func (s *ProductsStore) GetServiceByID(ctx context.Context, serviceID uuid.UUID) (*productsdomain.Service, error) {
	service := &productsdomain.Service{}
	err := s.db.Preload("Images").Preload("Category").Preload("Banners").First(service, serviceID).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, shared_errors.ErrNotFound
		default:
			return nil, err
		}
	}
	return service, nil

}
func (s *ProductsStore) UpdateService(ctx context.Context, service *productsdomain.Service, bannerID uuid.UUID) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("service_id = ?", service.ServiceID).Delete(&productsdomain.ServiceImage{}).Error; err != nil {
			return err
		}

		if err := tx.Where("service_id = ?", service.ServiceID).Delete(&marketingdomain.BannerService{}).Error; err != nil {
			return err
		}

		result := tx.Session(&gorm.Session{FullSaveAssociations: true}).Updates(service)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return shared_errors.ErrNotFound
		}
		if bannerID != uuid.Nil {
			bannerLink := &marketingdomain.BannerService{
				ServiceID: service.ServiceID,
				BannerID:  bannerID,
			}
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(bannerLink).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
=======
>>>>>>> dev
