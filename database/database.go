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
var Ctx = context.Background()

func SetupFirebase() {
	opt := option.WithCredentialsFile("database/serviceAccount.json")
	config := &firebase.Config{ProjectID: "test-e7825"}
	app, err := firebase.NewApp(Ctx, config, opt)
	if err != nil {
		log.Fatalln(err)
	}

	Client, err = app.Firestore(Ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// defer Client.Close()
}
