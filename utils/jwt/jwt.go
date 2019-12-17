package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	User                interface{} `json:"user"`
	ResourcePermissions interface{} `json:"resourcePermissions"`
	jwt.StandardClaims
}
type alg string

const (
	HS256 alg = "HS256"
)

var (
	// ErrTokenContextMissing denotes a token was not passed into the parsing
	// middleware's context.
	ErrTokenContextMissing = errors.New("token up for parsing was not passed through the context")

	// ErrTokenInvalid denotes a token was not able to be validated.
	ErrTokenInvalid = errors.New("JWT Token was invalid")

	// ErrTokenExpired denotes a token's expire header (exp) has since passed.
	ErrTokenExpired = errors.New("JWT Token is expired")

	// ErrTokenMalformed denotes a token was not formatted as a JWT token.
	ErrTokenMalformed = errors.New("JWT Token is malformed")

	// ErrTokenNotActive denotes a token's not before header (nbf) is in the
	// future.
	ErrTokenNotActive = errors.New("token is not valid yet")

	// ErrUnexpectedSigningMethod denotes a token was signed with an unexpected
	// signing method.
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
)

func Sign(payload map[string]interface{}, secret string, method alg) (string, error) {
	var signingMethod jwt.SigningMethod
	switch method {
	default:
		signingMethod = jwt.SigningMethodHS256
	}
	user, _ := payload["user"]
	resourcePermissions, _ := payload["resourcePermissions"]
	claims := UserClaims{
		User:                user,
		ResourcePermissions: resourcePermissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000,
		},
	}
	token := jwt.NewWithClaims(signingMethod, claims)
	signingKey := []byte(secret)
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}
