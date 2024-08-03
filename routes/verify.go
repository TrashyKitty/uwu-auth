package routes

import (
	"context"

	"github.com/Ant767/AuthBackend/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func VerifyAccount(c *gin.Context) {
	code := c.Params.ByName("code")
	filter := bson.D{{"verificationCode", code}}
	update := bson.D{
		{"$set", bson.D{
			{"verified", true},
		}},
	}

	collection := db.GetUsersCollection()

	_, updateErr := collection.UpdateOne(context.TODO(), filter, update)
	if updateErr != nil {
		c.String(200, "Failed to verify account")
		return
	} else {
		c.String(200, "Verified account!")
		return
	}
}
