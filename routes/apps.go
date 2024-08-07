package routes

import (
	"context"
	"net/http"

	"github.com/Ant767/AuthBackend/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type App struct {
	Name        string `json:"name"`
	AppIconURL  string `json:"appIconURL"`
	Description string `json:"description"`
	Deprecated  bool   `json:"deprecated"`
	AppID       string `json:"app_id"`
}

var apps = []App{
	{
		Name:        "testy :3",
		AppIconURL:  "https://media.discordapp.net/attachments/1205541941937442858/1269069384425406647/1018-3447020280.jpg?ex=66aeb877&is=66ad66f7&hm=d37848c69e1f1479ce6e7a1be9a87a6f6036c290cc1aec000e460563a6ba94b0&=&format=webp&width=1013&height=675",
		Description: "testy :3",
		Deprecated:  false,
		AppID:       "testy",
	},
}

type Role struct {
	Name  string `json:"name"`
	Power int    `json:"power"`
}

var roles = map[int]Role{
	0: {
		Name:  "Member",
		Power: 0,
	},
	1: {
		Name:  "Owner",
		Power: 10,
	},
	2: {
		Name:  "Founder",
		Power: 20,
	},
}

func GetAppsList(c *gin.Context) {
	c.JSON(200, apps)
}
func GetRolesList(c *gin.Context) {
	c.JSON(200, roles)
}

type CreateAppAssociationRequest struct {
	Token string `json:"token"`
	AppID string `json:"app_id"`
}

func findAppByID(apps []App, id string) *App {
	for _, app := range apps {
		if app.AppID == id {
			return &app
		}
	}
	return nil
}

func CreateAppAssociation(c *gin.Context) {
	var body CreateAppAssociationRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app := findAppByID(apps, body.AppID)

	if app != nil {
		filter := bson.M{
			"token": body.Token,
		}

		update := bson.M{
			"$addToSet": bson.M{
				"appAssociations": app.AppID,
			},
		}

		collection := db.GetUsersCollection()
		_, err := collection.UpdateOne(context.Background(), filter, update)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": true})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error": false})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"error": true})
		return
	}
}

func GetAppByID(c *gin.Context) {
	appID := c.Params.ByName("id")
	app := findAppByID(apps, appID)
	if app != nil {
		c.JSON(http.StatusOK, gin.H{"error": false, "app": app})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": true})
	}
}
