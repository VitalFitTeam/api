package main

import (
	"context"
	"log"
	"time"

	authdomain "github.com/vitalfit/api/internal/auth/domain"
	authservices "github.com/vitalfit/api/internal/auth/services"
	"github.com/vitalfit/api/internal/store"
	env "github.com/vitalfit/api/pkg/Env"
	dbg "github.com/vitalfit/api/pkg/db"
	"github.com/vitalfit/api/pkg/mailer"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"github.com/vitalfit/api/config"
)

func Seed(store store.Storage, db *gorm.DB) {
	ctx := context.Background()
	user := &authdomain.Users{
		FirstName:        "Super",
		LastName:         "Admin",
		Email:            env.GetString("ADMIN_EMAIL", ""),
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
	user.PasswordHash.Set(env.GetString("ADMIN_PASSWORD", ""))
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

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/vitalfit?sslmode=disable")
	conn, err := dbg.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	cfg := config.LoadConfig()

	mailer, err := mailer.NewResendClient(cfg.Mail.Resend.ApiKey, cfg.Mail.FromEmail)
	if err != nil {
		log.Fatal(err)
	}

	auth := authservices.NewJWTAuthenticator(cfg.Auth.Token.Secret, cfg.Auth.Token.Iss, cfg.Auth.Token.Iss)

	store := store.NewStorage(conn, *cfg, mailer, auth)

	Seed(store, conn)
}
