package main

import (
	"context"
	"encoding/json"
	"math/rand"

	"log"
	"os"
	"time"

	"strings"

	"github.com/google/uuid"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
	instructordomain "github.com/vitalfit/api/internal/modules/instructor/domain"
	inventorydomain "github.com/vitalfit/api/internal/modules/inventory/domain"
	marketingdomain "github.com/vitalfit/api/internal/modules/marketing/domain"
	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
	"github.com/vitalfit/api/internal/store"
	env "github.com/vitalfit/api/pkg/Env"
	"github.com/vitalfit/api/pkg/db"
	"github.com/vitalfit/api/pkg/pagination"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

type SeedStruct struct{}

func NewSeedStruct() *SeedStruct {
	return &SeedStruct{}
}

func (s *SeedStruct) Seed(store store.Storage, db *gorm.DB) {
	ctx := context.Background()
	//s.CreateSuperAdmin(store, db, ctx)
	//s.SeedPermissions(store, db, ctx)
	//s.SeedUsers(store, db, ctx)
	//s.SeedServiceCategories(store, db, ctx)
	//s.SeedBanners(store, db, ctx)
	//s.SeedInstructors(store, db, ctx)
	//s.SeedBranches(store, db, ctx)
	s.SeedEquipment(store, db, ctx)
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
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		specialties, err := store.Products.ListServiceCategories(ctx, tx)
		if err != nil {
			return err
		}

		if len(specialties) == 0 {
			log.Println("No specialties found to assign to instructors.")
			return nil
		}

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

			// Asignar especialidades aleatorias
			numSpecialtiesToAssign := r.Intn(len(specialties)) + 1 // Asignar de 1 a len(specialties)
			r.Shuffle(len(specialties), func(i, j int) {
				specialties[i], specialties[j] = specialties[j], specialties[i]
			})

			assignedSpecialties := make([]*productsdomain.ServiceCategory, 0, numSpecialtiesToAssign)
			for k := 0; k < numSpecialtiesToAssign; k++ {
				assignedSpecialties = append(assignedSpecialties, &specialties[k])
			}

			instructor := &instructordomain.Instructor{
				Biography: i.Biography,
				User:      user,
			}
			if err := store.Instructor.Create(ctx, tx, instructor); err != nil {
				log.Printf("Error creating instructor '%s': %v", i.FirstName, err)
				return err //rollback

			}

			specialtyIDs := make([]uuid.UUID, 0, len(assignedSpecialties))
			for _, s := range assignedSpecialties {
				specialtyIDs = append(specialtyIDs, s.CategoryID)
			}

			if err := store.Instructor.AssignInstructorSpecialtyTx(ctx, tx, instructor.InstructorID, specialtyIDs); err != nil {
				log.Printf("Error assigning specialties to instructor '%s': %v", i.FirstName, err)
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

type ServiceCategoryJSON struct {
	Name string `json:"name"`
}

func (s *SeedStruct) SeedServiceCategories(store store.Storage, db *gorm.DB, ctx context.Context) {
	jsonFile, err := os.ReadFile("./internal/migrate/seed/data/service_categories.json")
	if err != nil {
		log.Fatalf("Fatal error: could not read service_categories.json file: %v", err)
		return
	}

	var serviceCategoriesFromJSON []ServiceCategoryJSON
	if err = json.Unmarshal(jsonFile, &serviceCategoriesFromJSON); err != nil {
		log.Fatalf("Fatal error: could not decode JSON: %v", err)
		return
	}
	log.Printf("Found %d service categories in service_categories.json. Starting seeder...", len(serviceCategoriesFromJSON))

	err = db.Transaction(func(tx *gorm.DB) error {

		for _, sc := range serviceCategoriesFromJSON {
			serviceCategory := &productsdomain.ServiceCategory{
				Name: sc.Name,
			}
			err := store.Products.CreateServiceCategory(ctx, tx, serviceCategory)
			if err != nil {
				log.Printf("Error creating service category '%s': %v", sc.Name, err)
				return err //rollback
			}
		}
		// Commit
		return nil
	})

	if err != nil {
		log.Println("Error in service categories seeder, transaction was rolled back:", err)
		return
	}

	log.Println("Service categories seeder completed successfully.")
}

type bannerJSON struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	LinkURL  string `json:"link_url"`
	IsActive bool   `json:"is_active"`
}

func (s *SeedStruct) SeedBanners(store store.Storage, db *gorm.DB, ctx context.Context) {
	jsonFile, err := os.ReadFile("./internal/migrate/seed/data/banners.json")
	if err != nil {
		log.Fatalf("Fatal error: could not read banners.json file: %v", err)
		return
	}

	var bannersFromJSON []bannerJSON
	if err = json.Unmarshal(jsonFile, &bannersFromJSON); err != nil {
		log.Fatalf("Fatal error: could not decode JSON: %v", err)
		return
	}
	log.Printf("Found %d banners in banners.json. Starting seeder...", len(bannersFromJSON))

	err = db.Transaction(func(tx *gorm.DB) error {

		for _, b := range bannersFromJSON {
			banner := &marketingdomain.Banner{
				Name:     b.Name,
				ImageURL: b.ImageURL,
				LinkURL:  b.LinkURL,
				IsActive: b.IsActive,
			}
			err := store.Marketing.CreateBannerTX(ctx, tx, banner)
			if err != nil {
				log.Printf("Error creating banner '%s': %v", b.Name, err)
				return err //rollback
			}
		}
		return nil
	})

	if err != nil {
		log.Println("Error in banners seeder, transaction was rolled back:", err)
		return
	}

	log.Println("Banners seeder completed successfully.")

}

type branchJSON struct {
	Address        string               `json:"address"`
	Country        string               `json:"country"`
	Latitude       float64              `json:"latitude"`
	Longitude      float64              `json:"longitude"`
	MaxCapacity    int                  `json:"max_capacity"`
	Name           string               `json:"name"`
	OperatingHours []operatingHoursJSON `json:"operating_hours"`
	Phone          string               `json:"phone"`
	State          string               `json:"state"`
	Status         string               `json:"status"`
	TaxID          string               `json:"tax_id"`
}

type operatingHoursJSON struct {
	DayOfWeek string `json:"day_of_week"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
	IsClosed  bool   `json:"is_closed"`
}

func (s *SeedStruct) SeedBranches(store store.Storage, db *gorm.DB, ctx context.Context) {
	jsonFile, err := os.ReadFile("./internal/migrate/seed/data/branches.json")
	if err != nil {
		log.Fatalf("Fatal error: could not read branches.json file: %v", err)
		return
	}

	var branchesFromJSON []branchJSON
	if err = json.Unmarshal(jsonFile, &branchesFromJSON); err != nil {
		log.Fatalf("Fatal error: could not decode JSON: %v", err)
		return
	}
	log.Printf("Found %d branches in branches.json. Starting seeder...", len(branchesFromJSON))

	branchAdminsFeed, err := store.Users.GetBranchAdmins(ctx, pagination.PaginatedFeedQuery{})
	if err != nil {
		log.Printf("Error getting branch admins: %v", err)
		return
	}

	if len(branchAdminsFeed) == 0 {
		log.Println("Warning: No branch admins found. Branches will be created without a manager.")
	}

	err = db.Transaction(func(tx *gorm.DB) error {

		for _, b := range branchesFromJSON {
			branch := &branchdomain.Branch{
				Address:     b.Address,
				MaxCapacity: b.MaxCapacity,
				Name:        b.Name,
				Phone:       b.Phone,
				Status:      branchdomain.BranchStatusEnum(b.Status),
				TaxID:       b.TaxID,
			}

			state, err := store.Locations.FindOrCreateStateByCountry(ctx, b.State, b.Country)
			if err != nil {
				log.Printf("Error creating state '%s': %v", b.State, err)
				return err //rollback
			}
			branch.StateID = state.StateID

			for _, oh := range b.OperatingHours {
				branch.OperatingHours = append(branch.OperatingHours, branchdomain.OperatingHours{
					DayOfWeek: branchdomain.DayOfWeekEnum(oh.DayOfWeek),
					OpenTime:  &oh.OpenTime,
					CloseTime: &oh.CloseTime,
					IsClosed:  oh.IsClosed,
				})
			}

			if len(branchAdminsFeed) > 0 {

				randomIndex := rand.Intn(len(branchAdminsFeed))
				randomAdmin := branchAdminsFeed[randomIndex]

				branch.ManagerID = randomAdmin.UserID
			}

			err = store.Branches.Create(ctx, tx, branch)
			if err != nil {
				log.Printf("Error creating branch '%s': %v", b.Name, err)
				return err //rollback
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Seeder transaction failed: %v", err)
	}

	log.Println("Branchs seeder completed successfully.")

}

type equipmentJSON struct {
	Brand       string
	Category    string
	Description string
	Model       string
	Name        string
}

func (s *SeedStruct) SeedEquipment(store store.Storage, db *gorm.DB, ctx context.Context) {
	jsonFile, err := os.ReadFile("./internal/migrate/seed/data/equipments.json")
	if err != nil {
		log.Fatalf("Fatal error: could not read equipment.json file: %v", err)
		return
	}

	var equipmentFromJSON []equipmentJSON
	if err = json.Unmarshal(jsonFile, &equipmentFromJSON); err != nil {
		log.Fatalf("Fatal error: could not decode JSON: %v", err)
		return
	}
	log.Printf("Found %d equipment in equipment.json. Starting seeder...", len(equipmentFromJSON))

	err = db.Transaction(func(tx *gorm.DB) error {

		for _, e := range equipmentFromJSON {
			equipment := &inventorydomain.Equipment{
				Brand:       e.Brand,
				Category:    inventorydomain.EquipmentCategoryEnum(e.Category),
				Description: e.Description,
				Model:       e.Model,
				Name:        e.Name,
			}
			err := tx.WithContext(ctx).Create(&equipment).Error
			if err != nil {
				log.Printf("Error creating equipment '%s': %v", e.Name, err)
				return err //rollback
			}
		}
		return nil
	})
	if err != nil {
		log.Println("Error in equipment seeder, transaction was rolled back:", err)
		return
	}
	log.Println("Equipment seeder completed successfully.")

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
