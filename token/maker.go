package token

import "time"

// Maker is an interface for managing tokens
// A JWT and PASETO struct will be mad to implement this interface
type Maker interface {

	// Creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// chekc if input token is valid or not
	// if valid the method will return the payload object data
	VerifyToken(token string) (*Payload, error)

}