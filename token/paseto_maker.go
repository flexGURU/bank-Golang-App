package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto *paseto.V2
	symmetricKey []byte

}

func NewPasetoMaker(symmetricKey string) (Maker, error)  {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid keysize. Length of the key should be %d: ", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto: paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
	
}



func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayLoad(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err

}



func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	if err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil); err != nil {
		return nil, fmt.Errorf("invalid token %v", err)
	}

	if err := payload.Valid(); err != nil {
		return nil, err

	}

	return payload, nil




	
}
