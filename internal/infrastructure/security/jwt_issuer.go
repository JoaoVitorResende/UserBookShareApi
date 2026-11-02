package security

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTIssuer struct {
	key *rsa.PrivateKey
	ttl time.Duration
}

func NewJWTIssuer(privateKeyPEM []byte, ttl time.Duration) (*JWTIssuer, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, err
	}
	return &JWTIssuer{key: key, ttl: ttl}, nil
}

func (j *JWTIssuer) Issue(userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(j.key)
}