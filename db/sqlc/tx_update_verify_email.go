package db

import (
	"context"
	"log"
)

type VerifyEmailTxParams struct {
	UpdateVerifyEmailParams
}

type VerifyEmailTxResults struct {
	User

}


func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResults, error) {
	var result VerifyEmailTxResults

	err := store.execTX(ctx, func(q *Queries) error {

		var err error

		verifyEmailResult, err := q.UpdateVerifyEmail(ctx, arg.UpdateVerifyEmailParams)
		if err != nil {
			// log.Printf("error updating verify email", err)
			return err
		}

		userVerifiedResult, err := q.UpdateUserVerification(ctx, verifyEmailResult.Username)
		if err != nil {
			// log.Printf("error updating user", err)
			return err
		}

		result.User = userVerifiedResult

		return nil
	})
		return result, err

}