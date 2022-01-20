package mauth

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"os"
	"github.com/pkg/errors"
)

var (
	firebaseConfigFile = os.Getenv("FIREBASE_CONFIG_FILE")
)

func InitAuth() (*auth.Client, error) {
	opt := option.WithCredentialsFile(firebaseConfigFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing firebase auth (create firebase app)")
	}

	client, errAuth := app.Auth(context.Background())
	if errAuth != nil {
		return nil, errors.Wrap(errAuth, "error initializing firebase auth (creating client)")
	}

	return client, nil
}
