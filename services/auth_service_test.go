
package services

import (
	"errors"
	"hackathon/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	users                  map[string]*models.User
	err                    error
	revokeTokensBeforeTime int64
}

func (m *MockUserRepository) Create(user *models.User) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.users[user.Username]; exists {
		return errors.New("username already exists")
	}
	user.ID = uint(len(m.users) + 1)
	m.users[user.Username] = user
	return nil
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if user, exists := m.users[username]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) UpdateRevokeTokensBefore(user *models.User, timestamp int64) error {
	if m.err != nil {
		return m.err
	}
	m.revokeTokensBeforeTime = timestamp
	return nil
}

func newTestAuthService(repo *MockUserRepository) *AuthService {
	return NewAuthService(repo, []byte("secret"), 1)
}

func TestAuthService_Register(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		repo := &MockUserRepository{users: make(map[string]*models.User)}
		service := newTestAuthService(repo)

		err := service.Register("testuser", "password")
		assert.NoError(t, err)

		user, err := repo.FindByUsername("testuser")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Username)
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password"))
		assert.NoError(t, err)
	})

	t.Run("existing username", func(t *testing.T) {
		repo := &MockUserRepository{
			users: map[string]*models.User{
				"testuser": {Username: "testuser", Password: "hashedpassword"},
			},
		}
		service := newTestAuthService(repo)

		err := service.Register("testuser", "password")
		assert.Error(t, err)
		assert.Equal(t, ErrUserExists, err)
	})
}

func TestAuthService_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	t.Run("successful login", func(t *testing.T) {
		repo := &MockUserRepository{
			users: map[string]*models.User{
				"testuser": {ID: 1, Username: "testuser", Password: string(hashedPassword)},
			},
		}
		service := newTestAuthService(repo)

		tokenResponse, err := service.Login("testuser", "password")
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenResponse.Token)
		assert.True(t, tokenResponse.ExpiredTime > time.Now().Unix())
	})

	t.Run("user not found", func(t *testing.T) {
		repo := &MockUserRepository{users: make(map[string]*models.User)}
		service := newTestAuthService(repo)

		_, err := service.Login("testuser", "password")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCreds, err)
	})

	t.Run("invalid password", func(t *testing.T) {
		repo := &MockUserRepository{
			users: map[string]*models.User{
				"testuser": {ID: 1, Username: "testuser", Password: string(hashedPassword)},
			},
		}
		service := newTestAuthService(repo)

		_, err := service.Login("testuser", "wrongpassword")
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidCreds, err)
	})
}

func TestAuthService_RevokeToken(t *testing.T) {
	t.Run("successful token revocation", func(t *testing.T) {
		repo := &MockUserRepository{
			users: map[string]*models.User{
				"testuser": {ID: 1, Username: "testuser", Password: "hashedpassword"},
			},
		}
		service := newTestAuthService(repo)

		err := service.RevokeToken(1)
		assert.NoError(t, err)
		assert.True(t, repo.revokeTokensBeforeTime > 0)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := &MockUserRepository{users: make(map[string]*models.User)}
		service := newTestAuthService(repo)

		err := service.RevokeToken(1)
		assert.Error(t, err)
	})
}
