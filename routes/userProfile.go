package routes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Ant767/AuthBackend/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateImage(c *gin.Context, fileField string, databaseField string) {
	file, header, err := c.Request.FormFile(fileField)
	token := c.Request.FormValue("token")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileUuid := uuid.New()
	fileExtension := filepath.Ext(header.Filename)
	out, err := os.Create(fmt.Sprintf("./uploads/%s%s", fileUuid.String(), fileExtension))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filter := bson.D{{"token", token}}
	update := bson.D{
		{"$set", bson.D{
			{databaseField, fmt.Sprintf("%s%s", fileUuid.String(), fileExtension)},
		}},
	}

	collection := db.GetUsersCollection()

	_, updateErr := collection.UpdateOne(context.TODO(), filter, update)

	if updateErr != nil {
		c.JSON(http.StatusOK, gin.H{"error": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func UpdateProfilePicture(c *gin.Context) {
	UpdateImage(c, "avatar", "avatarFile")
}

func UpdateProfileBanner(c *gin.Context) {
	UpdateImage(c, "banner", "bannerFile")
}
