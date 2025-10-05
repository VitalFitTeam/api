package store

import (
	"context"
	"database/sql"

	"github.com/vitalfit/api/config"
	authdomain "github.com/vitalfit/api/internal/auth/domain"
	authrepository "github.com/vitalfit/api/internal/auth/repository"
	"gorm.io/gorm"
)

type Storage struct {
	Users authdomain.UserRepository
	Roles authdomain.RolesRepository
	config.Config
}

func NewStorage(db *gorm.DB, cfg config.Config) Storage {
	return Storage{
		Users:  authrepository.NewUserRepositoryDAO(db),
		Roles:  authrepository.NewRoleStore(db),
		Config: cfg,
	}

}

func withTX(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
