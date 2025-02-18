package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// Valid implements the jwt.Claims interface.


// GetAudience implements the jwt.Claims interface.
func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{}, nil
}

// GetExpirationTime implements the jwt.Claims interface.
func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

// GetIssuedAt implements the jwt.Claims interface.
func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

// GetNotBefore implements the jwt.Claims interface.
func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Now()), nil
}

// GetIssuer implements the jwt.Claims interface.
func (p *Payload) GetIssuer() (string, error) {
	return "", nil
}

// GetSubject implements the jwt.Claims interface.
func (p *Payload) GetSubject() (string, error) {
	return p.Username, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
	 return errors.New("Token has expired")
	}
	return nil
   }

// creates a new token payload with a specific username and duration
func NewPayLoad(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil

}
