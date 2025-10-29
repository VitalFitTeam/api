package instructorservices

import "github.com/vitalfit/api/internal/store"

type InstructorServices struct {
	store store.Storage
}

func NewInstructorServices(store store.Storage) *InstructorServices {
	return &InstructorServices{
		store: store,
	}
}
