package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var (
	App  *firebase.App
	Auth *auth.Client
)

func InitFirebase() error {
	opt := option.WithCredentialsFile("/app/internal/firebase/serviceAccountKey.json") // your Firebase admin SDK JSON file
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v", err)
	}

	App = app
	Auth = client
	log.Println("✅ Firebase initialized")

	return nil
}
