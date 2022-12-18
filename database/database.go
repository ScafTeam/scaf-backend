package database

import (
	"bufio"
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var Client *firestore.Client
var Key, ID string

// Use a service account
func readConfig() (string, string) {
	file, err := os.Open("database/config.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	key := scanner.Text()
	scanner.Scan()
	id := scanner.Text()
	return key, id
}

func SetupFirebase() {
	Key, ID = readConfig()
	log.Print(ID)
	log.Print(Key)
	ctx := context.Background()
	opt := option.WithCredentialsFile("database/serviceAccount.json")
	config := &firebase.Config{ProjectID: ID}
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
