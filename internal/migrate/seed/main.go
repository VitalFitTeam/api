package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	"github.com/vitalfit/api/internal/store"
	env "github.com/vitalfit/api/pkg/Env"
	dbg "github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func Seed(store store.Storage, db *gorm.DB) {
	ctx := context.Background()
	CreateSuperAdmin(store, db, ctx)
	SeedPermissions(store, db, ctx)
}

func CreateSuperAdmin(store store.Storage, db *gorm.DB, ctx context.Context) {
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
	err = dbg.WithTX(db, func(tx *gorm.DB) error {
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

func SeedPermissions(store store.Storage, db *gorm.DB, ctx context.Context) {

	jsonFile, err := os.ReadFile("./internal/migrate/seed/permission.json")
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

func main() {
	addr := env.GetString("DB_ADDR", "")
	conn, err := dbg.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	store := store.NewStorage(conn)

	Seed(store, conn)
}
