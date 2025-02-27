package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   string   `json:"user_id"`
	Issuer   string   `json:"iss"`
	Subject  string   `json:"sub"`
	Audience []string `json:"aud"`
}

type JWKSHandler struct {
	privateKey *rsa.PrivateKey
	issuer     string
}

func NewJWKSHandler(privateKey *rsa.PrivateKey) JWKSHandler {
	return JWKSHandler{privateKey: privateKey, issuer: "backend"}
}

func (m *JWKSHandler) Generate(user_id string) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * 7 * 30 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   user_id,
			Issuer:    m.issuer,
		},
		Issuer:   m.issuer,
		Subject:  user_id,
		UserID:   user_id,
		Audience: []string{"api-audience"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "key1"

	return token.SignedString(m.privateKey)
}

func (j *JWKSHandler) Verify(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.privateKey.Public(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

func (j *JWKSHandler) Validate() map[string]interface{} {
	publicKey := j.privateKey.Public().(*rsa.PublicKey)

	jwks := map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"kid": "key1",
				"n":   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
				"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
				"alg": "RS256",
				"use": "sig",
			},
		},
	}
	return jwks
}
