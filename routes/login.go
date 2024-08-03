package routes

import (
	"context"
	"errors"
	"net/http"

	"github.com/Ant767/AuthBackend/db"
	"github.com/Ant767/AuthBackend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginRequestBody struct {
	HandleOrEmail string `json:"identifier"`
	Password      string `json:"password"`
}

func getUserByHandleOrEmail(handleOrEmail string) (bson.M, error) {
	usersCollection := db.GetUsersCollection()
	filter := bson.D{{"email", handleOrEmail}}
	filterHandle := bson.D{{"handle", handleOrEmail}}
	var result bson.M
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		var result2 bson.M
		err2 := usersCollection.FindOne(context.TODO(), filterHandle).Decode(&result2)
		if err2 != nil {
			return nil, errors.New("User not found")
		} else {
			return result2, nil
		}
	} else {
		return result, nil
	}
}
func LoginRoute(c *gin.Context) {
	var body LoginRequestBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := getUserByHandleOrEmail(body.HandleOrEmail)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
	} else {
		hashedPassword := user["password"].(string)
		token := user["token"].(string)
		userID := user["id"].(string)
		if utils.IsCorrectPassword(hashedPassword, body.Password) {
			c.JSON(http.StatusOK, gin.H{"error": false, "message": "Successfully logged in!", "token": token, "userID": userID})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": true, "message": "Invalid password"})
		}
	}

}
