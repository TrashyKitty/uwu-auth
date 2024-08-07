package auth

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/Ant767/AuthBackend/utils"
	"github.com/google/uuid"
	"github.com/resend/resend-go/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func isValidHandle(handle string) bool {
	// Define the regex pattern to match lowercase letters, numbers, _, -, and .
	pattern := `^[a-z0-9_.-]+$`

	// Compile the regex
	re := regexp.MustCompile(pattern)

	// Check if the handle matches the pattern
	return re.MatchString(handle)
}
func RegisterAccount(resendKey string, collection *mongo.Collection, handle string, username string, password string, email string) error {
	if !isValidHandle(handle) {
		return errors.New("Handle can only be lowercase leetters, numbers, underscores (_), hyphens (-), and periods (.)")
	}
	handleFilter := bson.D{{"handle", handle}}
	handleFilter2 := bson.D{{"email", email}}

	var result bson.M
	var result2 bson.M
	err := collection.FindOne(context.TODO(), handleFilter).Decode(&result)
	err2 := collection.FindOne(context.TODO(), handleFilter2).Decode(&result2)

	if err != nil && err2 != nil {
		if err == mongo.ErrNoDocuments && err2 == mongo.ErrNoDocuments {
			hashedPassword, _ := utils.HashPassword(password)

			userID := uuid.New()
			verificationCode := uuid.New()
			token := utils.MakeToken(userID.String())
			if token == "" {
				return errors.New("Failed to generate token")
			}
			document := bson.D{
				{"username", username},
				{"handle", handle},
				{"password", hashedPassword},
				{"email", email},
				{"role", 0},
				{"id", userID.String()},
				{"token", token},
				{"verificationCode", verificationCode.String()},
				{"verified", false},
			}
			_, err3 := collection.InsertOne(context.TODO(), document)
			if err3 != nil {
				return err3
			}

			client := resend.NewClient(resendKey)

			params := &resend.SendEmailRequest{
				To:      []string{email},
				From:    "accounts@trashdev.org",
				Subject: "Verify your account",
				Html:    fmt.Sprintf("Verify your trashdev account by clicking <a href=\"https://auth.trashdev.org/verify/%s\">this link</a>", verificationCode.String()),
			}

			sent, _ := client.Emails.Send(params)

			fmt.Println(sent.Id)

			return nil
		}

		return errors.New("an error occurred")
	}
	return errors.New("user already found")
}
