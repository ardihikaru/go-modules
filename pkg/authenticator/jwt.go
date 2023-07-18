// Package authenticator provides functions to authenticate incoming requests
package authenticator

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwa"
	"golang.org/x/crypto/bcrypt"
)

// Session defines the user raw document to be inserted
type Session struct {
	AccountId      string  `json:"account_id"`
	UserId         string  `json:"user_id"`
	GymId          string  `json:"gym_id"`
	Username       string  `json:"username"`
	Role           string  `json:"role"`
	TrainingInHour float32 `json:"training_in_hour,omitempty"`
}

// JWTClaims is the data type of claims in Encode() function (jwtauth.go)
type JWTClaims map[string]interface{}

const (
	ClaimAccountIdKey   = "_id"
	ClaimUserIdKey      = "user_id"
	ClaimUsernameKey    = "username"
	ClaimPermissionsKey = "permissions"
	ClaimExpiredInKey   = "exp"
	ClaimIssuedAtKey    = "iat"
)

// MakeTokenAuth creates token authentication
func MakeTokenAuth(jwtAlgo jwa.SignatureAlgorithm, jwtSecret string) (*jwtauth.JWTAuth, error) {
	// generates token auth
	tokenAuth := jwtauth.New(jwtAlgo.String(), []byte(jwtSecret), nil)

	return tokenAuth, nil
}

// CreateAccessToken creates an access token
func CreateAccessToken(tokenAuth *jwtauth.JWTAuth, jwtClaims JWTClaims) string {
	// extracts token
	_, tokenString, _ := tokenAuth.Encode(jwtClaims)

	return tokenString
}

// CheckPasswordHash compares the input password with the stored password,
// and finally validates if they are match or not
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
