package database

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func CheckProjectNameUnique(user_email, project_name string) (bool, error) {
	iter := Client.Collection("projects").
		Where("name", "==", project_name).
		Where("author", "==", user_email).
		Documents(context.Background())

	var projectNum int

	for {
		_, err := iter.Next()
		if err == iterator.Done {
			return projectNum == 0, nil
		}
		projectNum++
		if err != nil {
			return false, err
		}
	}
}

func GetProjectDetail(user_email, project_name string) (*firestore.DocumentSnapshot, error) {
	iter := Client.Collection("projects").
		Where("name", "==", project_name).
		Where("author", "==", user_email).
		Documents(context.Background())

	var projectNum int
	var project *firestore.DocumentSnapshot

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		projectNum++
		if err != nil {
			return nil, err
		}
		project = doc
	}

	if projectNum == 0 || projectNum > 1 {
		return nil, nil
	}

	return project, nil
}
