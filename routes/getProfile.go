package routes

import (
	"context"
	"net/http"

	"github.com/Ant767/AuthBackend/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetProfile(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	handle := c.Params.ByName("handle")
	collection := db.GetUsersCollection()

	if auth != "" && handle == "me" {
		var result bson.M
		filter := bson.D{
			{"token", auth},
		}

		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query MongoDB"})
			}
			return
		}
		c.JSON(200, result)
	} else {
		var result bson.M
		filter := bson.D{
			{"handle", handle},
		}
		projection := bson.D{
			{"token", 0},
			{"email", 0},
			{"password", 0},
			{"verificationCode", 0},
		}
		opts := options.FindOne().SetProjection(projection)
		err := collection.FindOne(context.TODO(), filter, opts).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query MongoDB"})
			}
			return
		}
		c.JSON(200, result)
	}
}
