package mauth

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/pkg/errors"
)

func ElevateToAdmin(ctx context.Context, client *auth.Client, uid string) error {
	claims := map[string]interface{}{"admin": true}
	err := client.SetCustomUserClaims(ctx, uid, claims)
	if err != nil {
		return errors.Wrap(err, "error setting custom claims")
	}
	return nil
}

func RevokeAdmin(ctx context.Context, client *auth.Client, uid string) error {
	err := client.SetCustomUserClaims(ctx, uid, nil)
	if err != nil {
		return errors.Wrap(err, "error revoking custom claims")
	}
	return nil
}

func VerifyAdmin(ctx context.Context, client *auth.Client, uid string) error {
	user, err := client.GetUser(ctx, uid)
	if err != nil {
		return errors.Wrap(err, "failed to get user during verify admin")
	}

	if admin, ok := user.CustomClaims["admin"]; ok {
		if admin.(bool) {
			return nil
		}
	}
	return errors.New("not admin")
}
