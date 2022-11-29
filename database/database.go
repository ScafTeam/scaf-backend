package database

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var client *firestore.Client

// Use a service account
var ctx = context.Background()

func SetupFirebase() {
	opt := option.WithCredentialsFile("database/serviceAccount.json")
	config := &firebase.Config{ProjectID: "ProjectID"}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
}

func GetData(collection string, doc string) (map[string]interface{}, error) {
	dsnap, err := client.Collection(collection).Doc(doc).Get(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	fmt.Println(dsnap.Data())
	return dsnap.Data(), nil
}
