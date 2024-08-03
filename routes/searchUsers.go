package routes

import (
	"context"
	"net/http"

	"github.com/Ant767/AuthBackend/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SearchUsers(c *gin.Context) {
	username := c.Query("username")
	handle := c.Query("handle")

	filter := bson.D{}
	if handle != "" {
		filter = append(filter, bson.E{"handle", bson.M{"$regex": handle, "$options": "i"}}) // Case-insensitive
	}
	if username != "" {
		filter = append(filter, bson.E{"username", bson.M{"$regex": username, "$options": "i"}})
	}

	projection := bson.D{
		{"password", 0},         // Include the 'handle' field
		{"email", 0},            // Include the 'handle' field
		{"token", 0},            // Include the 'handle' field
		{"verificationCode", 0}, // Include the 'handle' field
	}

	collection := db.GetUsersCollection()

	opts := options.Find().SetProjection(projection)
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query MongoDB"})
		return
	}
	defer cursor.Close(context.TODO())
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode MongoDB results"})
		return
	}

	// Return only the 'handle' field in the response
	handles := make([]interface{}, len(results))
	for i, result := range results {
		if handle, ok := result["handle"]; ok {
			handles[i] = handle
		} else {
			handles[i] = nil
		}
	}

	// Return results as JSON
	c.JSON(http.StatusOK, handles)
}
