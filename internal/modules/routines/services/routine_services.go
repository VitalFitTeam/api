package routineservices

import "github.com/vitalfit/api/internal/store"

type RoutineService struct {
	store store.Storage
}

func NewRoutineService(store store.Storage) *RoutineService {
	return &RoutineService{
		store: store,
	}

}
