package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

const minSecretKeySize = 32



func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid length of secret Key: required lenght %d", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}


func (jwtmaker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {

	payload, err := NewPayLoad(username, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(jwtmaker.secretKey))

	return token, payload, nil



}


func (jwtmaker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,fmt.Errorf("unexpected signing method: %v: ", token.Header["alg"])
		}

		return []byte(jwtmaker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{} , keyFunc)
	if err != nil {
		return nil, fmt.Errorf("problem parsing the token %w", err)
	}

	// extract the claims
	if payload, ok := jwtToken.Claims.(*Payload); ok && jwtToken.Valid {
		return payload, nil

	} else {
		return nil, fmt.Errorf("invalid Token")
	}
	
}
