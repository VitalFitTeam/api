package routinerepository

import "gorm.io/gorm"

type RoutineStore struct {
	db *gorm.DB
}

func NewRoutineStore(db *gorm.DB) *RoutineStore {
	return &RoutineStore{
		db: db,
	}

}
