package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
)

type Claims struct {
	Sub   uuid.UUID `json:"sub"`
	Email string    `json:"email"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewManager(secret string, accessTTLSeconds, refreshTTLSeconds int) *Manager {
	return &Manager{
		secret:     []byte(secret),
		accessTTL:  time.Duration(accessTTLSeconds) * time.Second,
		refreshTTL: time.Duration(refreshTTLSeconds) * time.Second,
	}
}

func (m *Manager) generate(userUUID uuid.UUID, email string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		Sub:   userUUID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			ID:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) GenerateAccessToken(userUUID uuid.UUID, email string) (string, error) {
	return m.generate(userUUID, email, m.accessTTL)
}

func (m *Manager) GenerateRefreshToken(userUUID uuid.UUID, email string) (string, error) {
	return m.generate(userUUID, email, m.refreshTTL)
}

func (m *Manager) AccessTTLSeconds() int {
	return int(m.accessTTL.Seconds())
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
