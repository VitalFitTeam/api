package authhandlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/vitalfit/api/config"
	"github.com/vitalfit/api/internal/app"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	authmocks "github.com/vitalfit/api/internal/modules/auth/mocks"
	authservices "github.com/vitalfit/api/internal/modules/auth/services"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	mailermocks "github.com/vitalfit/api/pkg/mailer/mocks"
)

// newTestUser is a helper function to create a test user with a hashed password.
func newTestUser(email, password string) *authdomain.Users {
	user := &authdomain.Users{
		UserID:    uuid.New(),
		Email:     email,
		FirstName: "Test",
		LastName:  "User",
	}
	// It's important to set the password for the login test to work
	err := user.PasswordHash.Set(password)
	if err != nil {
		panic(err)
	}
	return user
}

func newLoginPayload(email, password, context string) *bytes.Buffer {
	payload := map[string]string{
		"email":    email,
		"password": password,
		"context":  context,
	}
	body, _ := json.Marshal(payload)
	return bytes.NewBuffer(body)
}

func TestWhoAmI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testApp := app.NewTestApplication(t, config.LoadConfig())
	mux := testApp.Mount()

	mockUser := &authdomain.Users{
		UserID:    uuid.New(),
		Email:     "testuser@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	testToken, err := testApp.Services.AuthServices.GenerateToken(mockUser)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/user/whoami", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := app.ExecuteRequest(req, mux)
		app.CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {

		userStoreMock, ok := testApp.Store.Users.(*authmocks.UserStoreMock)
		if !ok {
			t.Fatalf("Failed to cast store.Users to *authmocks.UserStoreMock")
		}

		userStoreMock.On("GetByID", mock.Anything, mockUser.UserID).Return(mockUser, nil).Once()

		req, err := http.NewRequest(http.MethodGet, "/v1/user/whoami", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := app.ExecuteRequest(req, mux)

		app.CheckResponseCode(t, http.StatusOK, rr.Code)

		userStoreMock.AssertExpectations(t)
	})
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testApp := app.NewTestApplication(t, config.LoadConfig())
	mux := testApp.Mount()

	// Get the UserStore mock to set expectations
	userStoreMock, ok := testApp.Store.Users.(*authmocks.UserStoreMock)
	if !ok {
		t.Fatalf("Failed to cast store.Users to *authmocks.UserStoreMock")
	}

	testPassword := "password123"
	mockUser := newTestUser("login@example.com", testPassword)

	t.Run("should fail with invalid credentials (wrong password)", func(t *testing.T) {
		// We configure the mock to return the user when searching by email
		userStoreMock.On("GetByEmail", mock.Anything, mockUser.Email).Return(mockUser, nil).Once()

		body := newLoginPayload(mockUser.Email, "wrongpassword", "web")
		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/login", body)
		req.Header.Set("Content-Type", "application/json")

		rr := app.ExecuteRequest(req, mux)

		app.CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
		userStoreMock.AssertExpectations(t)
	})

	t.Run("should succeed with valid credentials", func(t *testing.T) {
		// We configure the mock to return the user
		userStoreMock.On("GetByEmail", mock.Anything, mockUser.Email).Return(mockUser, nil).Once()

		body := newLoginPayload(mockUser.Email, testPassword, "web")
		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/login", body)
		req.Header.Set("Content-Type", "application/json")

		rr := app.ExecuteRequest(req, mux)

		app.CheckResponseCode(t, http.StatusOK, rr.Code)

		// We verify that the response contains a token
		var responseBody map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}
		if _, exists := responseBody["token"]; !exists {
			t.Errorf("Expected response to contain a token, but it didn't")
		}

		userStoreMock.AssertExpectations(t)
	})

	t.Run("should fail when a client tries to log in to the dashboard", func(t *testing.T) {
		clientUser := newTestUser("client-login@example.com", testPassword)
		clientUser.Role = authdomain.Roles{Name: "client"}

		// Mock for GetByEmail, it should return a user with his role
		userStoreMock.On("GetByEmail", mock.Anything, clientUser.Email).Return(clientUser, nil).Once()

		body := newLoginPayload(clientUser.Email, testPassword, "dashboard")
		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/login", body)
		req.Header.Set("Content-Type", "application/json")

		rr := app.ExecuteRequest(req, mux)

		// We wait for a 403 forbidden
		app.CheckResponseCode(t, http.StatusForbidden, rr.Code)
		userStoreMock.AssertExpectations(t)

	})
}

func TestRegisterUserClientHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testApp := app.NewTestApplication(t, config.LoadConfig())
	mux := testApp.Mount()

	userStoreMock := testApp.Store.Users.(*authmocks.UserStoreMock)
	roleStoreMock := testApp.Store.Roles.(*authmocks.RoleStoreMock)
	mailerMock := testApp.Services.AuthServices.(*authservices.AuthService).Mailer.(*mailermocks.MockMailer)

	t.Run("should register a new client successfully", func(t *testing.T) {

		mockRole := &authdomain.Roles{RoleID: uuid.New(), Name: "client"}
		roleStoreMock.On("GetByName", mock.Anything, "client").Return(mockRole, nil).Once()

		userStoreMock.On("CreateAndInvitate", mock.Anything, mock.AnythingOfType("*authdomain.Users"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()

		mailerMock.On("Send",
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.Anything,
			mock.AnythingOfType("bool"),
		).Return(http.StatusOK, nil).Once()

		payload := map[string]string{
			"first_name":        "New",
			"last_name":         "User",
			"email":             "newuser@example.com",
			"phone":             "123456789",
			"identity_document": "1234567890",
			"password":          "Password123!",
			"birth_date":        "2000-01-01",
			"gender":            "male",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := app.ExecuteRequest(req, mux)

		app.CheckResponseCode(t, http.StatusCreated, rr.Code)

		// We verify that the mock expectations were met
		userStoreMock.AssertExpectations(t)
		roleStoreMock.AssertExpectations(t)
		mailerMock.AssertExpectations(t)
	})

	t.Run("should fail with invalid payload", func(t *testing.T) {
		// Payload without required fields (e.g., email)
		payload := map[string]string{"first_name": "Missing"}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := app.ExecuteRequest(req, mux)

		// Gin returns 400 for binding validation errors
		app.CheckResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should fail with weak password", func(t *testing.T) {
		payload := map[string]string{
			"first_name":        "Weak",
			"last_name":         "Password",
			"email":             "weakpass@example.com",
			"phone":             "123456789",
			"identity_document": "1234567890",
			"password":          "password", // Contraseña débil
			"birth_date":        "2000-01-01",
			"gender":            "male",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := app.ExecuteRequest(req, mux)
		app.CheckResponseCode(t, http.StatusBadRequest, rr.Code)
	})
}

func TestActivateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testApp := app.NewTestApplication(t, config.LoadConfig())
	mux := testApp.Mount()
	userStoreMock := testApp.Store.Users.(*authmocks.UserStoreMock)

	t.Run("should activate user with a valid code", func(t *testing.T) {
		activationCode := "valid-code"
		userStoreMock.On("Activate", mock.Anything, activationCode).Return(nil).Once()

		payload := map[string]string{"code": activationCode}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, "/v1/auth/activate", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := app.ExecuteRequest(req, mux)

		app.CheckResponseCode(t, http.StatusNoContent, rr.Code)
		userStoreMock.AssertExpectations(t)
	})
}

func TestRegisterUserStaffHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testApp := app.NewTestApplication(t, config.LoadConfig())
	mux := testApp.Mount()

	userStoreMock := testApp.Store.Users.(*authmocks.UserStoreMock)
	roleStoreMock := testApp.Store.Roles.(*authmocks.RoleStoreMock)
	mailerMock := testApp.Services.AuthServices.(*authservices.AuthService).Mailer.(*mailermocks.MockMailer)

	// Admin user who will perform the action
	adminUser := &authdomain.Users{
		UserID: uuid.New(),
		Email:  "admin@example.com",
		Role:   authdomain.Roles{Name: "super_admin", Level: 99}, // Sufficient level
	}
	adminToken, _ := testApp.Services.AuthServices.GenerateToken(adminUser)

	// Client user (non-admin) to test denied access
	clientUser := &authdomain.Users{
		UserID: uuid.New(),
		Email:  "client@example.com",
		Role:   authdomain.Roles{Name: "client", Level: 1}, // Insufficient level
	}
	clientToken, _ := testApp.Services.AuthServices.GenerateToken(clientUser)

	t.Run("should fail when unauthenticated", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/register-staff", nil)
		req.Header.Set("Content-Type", "application/json")
		rr := app.ExecuteRequest(req, mux)
		app.CheckResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should fail when user is not an admin", func(t *testing.T) {
		payload := map[string]string{"role_name": "instructor"}
		body, _ := json.Marshal(payload)

		// Mock for the middleware: GetByID for the client user
		userStoreMock.On("GetByID", mock.Anything, clientUser.UserID).Return(clientUser, nil).Once()
		// Mock for the middleware: GetByName for the required 'branch_admin' role
		roleStoreMock.On("GetByName", mock.Anything, "branch_admin").Return(&authdomain.Roles{Level: 10}, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/register-staff", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+clientToken)

		rr := app.ExecuteRequest(req, mux)
		app.CheckResponseCode(t, http.StatusForbidden, rr.Code)

		userStoreMock.AssertExpectations(t)
		roleStoreMock.AssertExpectations(t)
	})

	t.Run("should fail when trying to register a client", func(t *testing.T) {
		// The handler has a guard to prevent this. It shouldn't reach the services.
		userStoreMock.On("GetByID", mock.Anything, adminUser.UserID).Return(adminUser, nil).Once()
		roleStoreMock.On("GetByName", mock.Anything, "branch_admin").Return(&authdomain.Roles{Level: 50}, nil).Once()

		payload := map[string]string{
			"first_name": "Client", "last_name": "User", "email": "shouldfail@example.com",
			"phone": "111", "identity_document": "111", "password": "Password123!",
			"birth_date": "2000-01-01", "gender": "male",
			"role_name": "client", // Trying to register a client
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/register-staff", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)

		rr := app.ExecuteRequest(req, mux)
		app.CheckResponseCode(t, http.StatusBadRequest, rr.Code)

		userStoreMock.AssertExpectations(t)
		roleStoreMock.AssertExpectations(t)
	})

	staffRoles := []string{"instructor", "accountant", "recepcionist", "branch_admin", "super_admin"}

	for _, role := range staffRoles {
		t.Run("should succeed registering a "+role, func(t *testing.T) {
			userStoreMock.On("GetByID", mock.Anything, adminUser.UserID).Return(adminUser, nil).Once()
			roleStoreMock.On("GetByName", mock.Anything, "branch_admin").Return(&authdomain.Roles{Level: 50}, nil).Once()

			mockRole := &authdomain.Roles{RoleID: uuid.New(), Name: role}
			roleStoreMock.On("GetByName", mock.Anything, role).Return(mockRole, nil).Once()
			userStoreMock.On("CreateAndInvitate", mock.Anything, mock.AnythingOfType("*authdomain.Users"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()
			mailerMock.On("Send", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("bool")).Return(http.StatusOK, nil).Once()

			email := "new" + role + "@example.com"
			identityDoc := "doc-" + role

			payload := map[string]string{
				"first_name":        "New",
				"last_name":         role,
				"email":             email,
				"phone":             "987654321",
				"identity_document": identityDoc,
				"password":          "StaffPassword123!",
				"birth_date":        "1995-01-01",
				"gender":            "female",
				"role_name":         role,
			}
			body, _ := json.Marshal(payload)

			req, _ := http.NewRequest(http.MethodPost, "/v1/auth/register-staff", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+adminToken)

			rr := app.ExecuteRequest(req, mux)
			app.CheckResponseCode(t, http.StatusCreated, rr.Code)

			userStoreMock.AssertExpectations(t)
			roleStoreMock.AssertExpectations(t)
			mailerMock.AssertExpectations(t)
		})
	}
}

func TestResetPasswordHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testApp := app.NewTestApplication(t, config.LoadConfig())
	mux := testApp.Mount()

	userStoreMock := testApp.Store.Users.(*authmocks.UserStoreMock)

	t.Run("should fail with mismatched passwords", func(t *testing.T) {
		payload := map[string]string{
			"token":            "valid-token",
			"password":         "NewPassword123!",
			"confirm_password": "mismatchedpassword",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/password/reset", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := app.ExecuteRequest(req, mux)
		app.CheckResponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should fail if token is invalid", func(t *testing.T) {
		userStoreMock.On("ResetUserPassword", mock.Anything, "invalid-token", mock.AnythingOfType("*authdomain.Users")).Return(shared_errors.ErrNotFound).Once()

		payload := map[string]string{
			"token":            "invalid-token",
			"password":         "NewPassword123!",
			"confirm_password": "NewPassword123!",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/password/reset", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := app.ExecuteRequest(req, mux)
		app.CheckResponseCode(t, http.StatusBadRequest, rr.Code)
		userStoreMock.AssertExpectations(t)
	})

	t.Run("should succeed with a valid token and matching passwords", func(t *testing.T) {
		userStoreMock.On("ResetUserPassword", mock.Anything, "valid-token", mock.AnythingOfType("*authdomain.Users")).Return(nil).Once()

		payload := map[string]string{
			"token":            "valid-token",
			"password":         "NewPassword123!",
			"confirm_password": "NewPassword123!",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/password/reset", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := app.ExecuteRequest(req, mux)
		app.CheckResponseCode(t, http.StatusOK, rr.Code)
		userStoreMock.AssertExpectations(t)
	})
}

func TestForgotPasswordHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testApp := app.NewTestApplication(t, config.LoadConfig())
	mux := testApp.Mount()

	userStoreMock := testApp.Store.Users.(*authmocks.UserStoreMock)
	mailerMock := testApp.Services.AuthServices.(*authservices.AuthService).Mailer.(*mailermocks.MockMailer)
	mockUser := newTestUser("forgot@example.com", "password123")

	t.Run("should fail if user does not exist", func(t *testing.T) {
		nonExistentEmail := "notfound@example.com"
		userStoreMock.On("GetByEmail", mock.Anything, nonExistentEmail).Return(nil, shared_errors.ErrNotFound).Once()

		payload := map[string]string{"email": nonExistentEmail}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/password/forgot", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := app.ExecuteRequest(req, mux)

		app.CheckResponseCode(t, http.StatusNotFound, rr.Code)
		userStoreMock.AssertExpectations(t)
	})

	t.Run("should succeed if user exists", func(t *testing.T) {
		// Mock the sequence of calls for a successful request
		userStoreMock.On("GetByEmail", mock.Anything, mockUser.Email).Return(mockUser, nil).Times(2)
		userStoreMock.On("CreatePasswordResetToken", mock.Anything, mockUser.UserID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()
		mailerMock.On("Send", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("bool")).Return(http.StatusOK, nil).Once()

		payload := map[string]string{"email": mockUser.Email}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/v1/auth/password/forgot", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rr := app.ExecuteRequest(req, mux)

		app.CheckResponseCode(t, http.StatusOK, rr.Code)
		userStoreMock.AssertExpectations(t)
		mailerMock.AssertExpectations(t)
	})
}
