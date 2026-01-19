package routinehandlers

import (
	"time"

	"github.com/google/uuid"
)

type CreateRoutineRequest struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Level       string               `json:"level"`
	ServiceID   *uuid.UUID           `json:"service_id"`
	Exercises   []RoutineExerciseDTO `json:"exercises"`
}

type RoutineExerciseDTO struct {
	ExerciseID uuid.UUID `json:"exercise_id"`
	Sets       int       `json:"sets"`
	Reps       string    `json:"reps"`
	RestTime   int       `json:"rest_time"`
	Order      int       `json:"order"`
	Notes      string    `json:"notes"`
}

type AssignRoutineRequest struct {
	ClientID  uuid.UUID  `json:"client_id"`
	RoutineID uuid.UUID  `json:"routine_id"`
	DueDate   *time.Time `json:"due_date"` //optional
}

type CreateExerciseRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	VideoURL    string `json:"video_url"`
	MuscleGroup string `json:"muscle_group"`
}
