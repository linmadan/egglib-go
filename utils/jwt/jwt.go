package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

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

func Sign(claims jwt.Claims, secret string, method alg) (string, error) {
	var signingMethod jwt.SigningMethod
	switch method {
	default:
		signingMethod = jwt.SigningMethodHS256
	}
	token := jwt.NewWithClaims(signingMethod, claims)
	signingKey := []byte(secret)
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

func Valid(tokenString string, claims jwt.Claims, secret string, ) (bool, *jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
	if err != nil {
		return false, nil, err
	}
	if token.Valid {
		return true, token, nil
	} else {
		return false, nil, nil
	}
}

func IsExpired(err error) bool {
	ve, ok := err.(*jwt.ValidationError)
	if ok && ve.Errors == jwt.ValidationErrorExpired {
		return true
	}
	return false
}
