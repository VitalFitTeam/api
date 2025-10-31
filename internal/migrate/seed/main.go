package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"strings"

	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	instructordomain "github.com/vitalfit/api/internal/modules/instructor/domain"
	"github.com/vitalfit/api/internal/store"
	env "github.com/vitalfit/api/pkg/Env"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

type SeedStruct struct{}

func NewSeedStruct() *SeedStruct {
	return &SeedStruct{}
}

func (s *SeedStruct) Seed(store store.Storage, db *gorm.DB) {
	ctx := context.Background()
	s.CreateSuperAdmin(store, db, ctx)
	s.SeedPermissions(store, db, ctx)
	s.SeedInstructors(store, db, ctx)
	s.SeedUsers(store, db, ctx)
}

func (s *SeedStruct) CreateSuperAdmin(store store.Storage, db *gorm.DB, ctx context.Context) {
	email := env.GetString("ADMIN_EMAIL", "")
	password := env.GetString("ADMIN_PASSWORD", "")
	user := &authdomain.Users{
		FirstName:        "Super",
		LastName:         "Admin",
		Email:            email,
		Phone:            "+581235467890",
		IdentityDocument: "V-1234567891",
		Gender:           "male",
		IsValidated:      true,
	}
	date, err := time.Parse(time.RFC3339, "1990-01-01T00:00:00Z")
	if err != nil {
		log.Println("Error parsing date", err)
		return
	}
	user.BirthDate = date
	role, err := store.Roles.GetByName(ctx, "super_admin")
	if err != nil {
		log.Println("Error getting the role", err)
		return
	}
	user.RoleID = role.RoleID
	user.PasswordHash.Set(password)
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Println("Error creating the user", err)
		return
	}

	log.Println("User created successfully")
}

type permissionJSON struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *SeedStruct) SeedPermissions(store store.Storage, db *gorm.DB, ctx context.Context) {

	jsonFile, err := os.ReadFile("./internal/migrate/seed/data/permission.json")
	if err != nil {
		log.Fatalf("Fatal error: could not read permissions.json file: %v", err)
		return
	}

	var permissionsFromJSON []permissionJSON
	if err = json.Unmarshal(jsonFile, &permissionsFromJSON); err != nil {
		log.Fatalf("Fatal error: could not decode JSON: %v", err)
		return
	}
	log.Printf("Found %d permissions in permissions.json. Starting seeder...", len(permissionsFromJSON))

	err = db.Transaction(func(tx *gorm.DB) error {

		for _, p := range permissionsFromJSON {

			permissionToCreate := &authdomain.Permission{
				Name:        p.Name,
				Description: p.Description,
			}

			if err := store.Roles.CreatePermission(ctx, tx, permissionToCreate); err != nil {
				log.Printf("Error creating permission '%s': %v", p.Name, err)
				return err //rollback
			}
		}

		// Commit
		return nil
	})

	if err != nil {
		log.Println("Error in permissions seeder, transaction was rolled back:", err)
		return
	}

	log.Println("Permissions seeder completed successfully.")
}

type instructorJSON struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	IdentityDocument  string `json:"identity_document"`
	BirthDate         string `json:"birth_date"`
	Gender            string `json:"gender"`
	ProfilePictureURL string `json:"profile_picture_url"`
	Speciality        string `json:"speciality"`
	Biography         string `json:"biography"`
}

type UserJSON struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	IdentityDocument  string `json:"identity_document"`
	BirthDate         string `json:"birth_date"`
	Gender            string `json:"gender"`
	ProfilePictureURL string `json:"profile_picture_url"`
	RoleName          string `json:"role_name"`
}

func (s *SeedStruct) SeedInstructors(store store.Storage, db *gorm.DB, ctx context.Context) {
	password := env.GetString("ADMIN_PASSWORD", "")
	jsonFile, err := os.ReadFile("./internal/migrate/seed/data/instructor.json")
	if err != nil {
		log.Fatalf("Fatal error: could not read instructors.json file: %v", err)
	}

	var instructorsFromJSON []instructorJSON
	if err = json.Unmarshal(jsonFile, &instructorsFromJSON); err != nil {
		log.Fatalf("Fatal error: could not decode JSON: %v", err)
	}
	log.Printf("Found %d instructors in instructors.json. Starting seeder...", len(instructorsFromJSON))

	err = db.Transaction(func(tx *gorm.DB) error {

		for _, i := range instructorsFromJSON {
			user := &authdomain.Users{
				FirstName:         i.FirstName,
				LastName:          i.LastName,
				Email:             i.Email,
				Phone:             i.Phone,
				IdentityDocument:  i.IdentityDocument,
				Gender:            authdomain.GenderEnum(strings.ToLower(i.Gender)),
				ProfilePictureURL: i.ProfilePictureURL,
				IsValidated:       true,
			}
			date, err := time.Parse("2006-01-02", i.BirthDate)
			if err != nil {
				date, err = time.Parse(time.RFC3339, i.BirthDate)
				if err != nil {
					log.Printf("Error parsing date for instructor '%s': %v", i.Email, err)
					return err
				}
			}
			user.BirthDate = date

			user.PasswordHash.Set(password)

			role, err := store.Roles.GetByName(ctx, "instructor")
			if err != nil {
				log.Println("Error getting the role", err)
				return err
			}
			user.RoleID = role.RoleID
			instructor := &instructordomain.Instructor{
				Speciality: i.Speciality,
				Biography:  i.Biography,
				User:       user,
			}
			if err := store.Instructor.Create(ctx, tx, instructor); err != nil {
				log.Printf("Error creating instructor '%s': %v", i.FirstName, err)
				return err //rollback

			}

		}
		// Commit
		return nil
	})

	if err != nil {
		log.Println("Error in instructors seeder, transaction was rolled back:", err)
		return
	}

	log.Println("Instructors seeder completed successfully.")
}

func (s *SeedStruct) SeedUsers(store store.Storage, dbg *gorm.DB, ctx context.Context) {
	password := env.GetString("ADMIN_PASSWORD", "")
	jsonFile, err := os.ReadFile("./internal/migrate/seed/data/user.json")
	if err != nil {
		log.Fatalf("Fatal error: could not read users.json file: %v", err)
		return
	}

	var usersFromJSON []UserJSON
	if err = json.Unmarshal(jsonFile, &usersFromJSON); err != nil {
		log.Fatalf("Fatal error: could not decode JSON: %v", err)
		return
	}
	log.Printf("Found %d users in users.json. Starting seeder...", len(usersFromJSON))

	err = db.WithTX(dbg, func(tx *gorm.DB) error {

		for _, u := range usersFromJSON {
			user := &authdomain.Users{
				FirstName:         u.FirstName,
				LastName:          u.LastName,
				Email:             u.Email,
				Phone:             u.Phone,
				IdentityDocument:  u.IdentityDocument,
				Gender:            authdomain.GenderEnum(strings.ToLower(u.Gender)),
				ProfilePictureURL: u.ProfilePictureURL,
				IsValidated:       true,
			}
			date, err := time.Parse("2006-01-02", u.BirthDate)
			if err != nil {
				date, err = time.Parse(time.RFC3339, u.BirthDate)
				if err != nil {
					log.Printf("Error parsing date for user '%s': %v", u.Email, err)
					return err
				}
			}
			user.BirthDate = date

			user.PasswordHash.Set(password)

			role, err := store.Roles.GetByName(ctx, u.RoleName)
			if err != nil {
				log.Println("Error getting the role", err)
				return err
			}
			user.RoleID = role.RoleID
			if err := store.Users.Create(ctx, tx, user); err != nil {
				log.Printf("Error creating user '%s': %v", u.FirstName, err)
				return err //rollback
			}
		}
		// Commit
		return nil
	})

	if err != nil {
		log.Println("Error in users seeder, transaction was rolled back:", err)
		return
	}
	log.Println("Users seeder completed successfully.")
}

func main() {
	addr := env.GetString("DB_ADDR", "")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	store := store.NewStorage(conn)
	s := NewSeedStruct()
	s.Seed(store, conn)

}
