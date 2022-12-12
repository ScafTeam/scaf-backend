package database

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var Client *firestore.Client

// Use a service account

func SetupFirebase() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("database/serviceAccount.json")
	config := &firebase.Config{ProjectID: "test-e7825"}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalln(err)
	}

	Client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// defer Client.Close()
}
