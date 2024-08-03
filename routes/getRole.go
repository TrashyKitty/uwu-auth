package routes

import (
	"context"
	"errors"

	"github.com/Ant767/AuthBackend/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetRoleNumber(handle string) (error, int32) {
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
	collection := db.GetUsersCollection()
	opts := options.FindOne().SetProjection(projection)
	err := collection.FindOne(context.TODO(), filter, opts).Decode(&result)
	if err != nil {
		return errors.New("Invalid role"), 0
	}
	role, hasRole := result["role"].(int32)
	if hasRole {
		return nil, role
	} else {
		return nil, 0
	}
}

func GetRole(c *gin.Context) {
	err, result := GetRoleNumber(c.Params.ByName("handle"))
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"role": result})
}
