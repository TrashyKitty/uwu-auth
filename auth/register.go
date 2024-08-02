package auth

import (
	"context"
	"errors"

	"github.com/Ant767/AuthBackend/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterAccount(collection *mongo.Collection, handle string, username string, password string, email string) error {
	handleFilter := bson.D{{"handle", handle}}
	handleFilter2 := bson.D{{"email", email}}

	var result bson.M
	var result2 bson.M
	err := collection.FindOne(context.TODO(), handleFilter).Decode(&result)
	err2 := collection.FindOne(context.TODO(), handleFilter2).Decode(&result2)

	if err != nil && err2 != nil {
		if err == mongo.ErrNoDocuments && err2 == mongo.ErrNoDocuments {
			_, hashedPassword := utils.HashPassword(password)

			userID := uuid.New()
			token := utils.MakeToken(userID.String())
			if token == "" {
				return errors.New("Failed to generate token")
			}
			document := bson.D{
				{"username", username},
				{"handle", handle},
				{"password", hashedPassword},
				{"email", email},
				{"id", userID.String()},
				{"token", token},
			}
			_, err3 := collection.InsertOne(context.TODO(), document)
			if err3 != nil {
				return err3
			}
			return nil
		}

		return errors.New("an error occurred")
	}
	return errors.New("user already found")
}
