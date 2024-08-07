package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/Ant767/AuthBackend/db"
	"github.com/Ant767/AuthBackend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/time/rate"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}
	defer conn.Close()

	for {
		var v interface{}
		err := conn.ReadJSON(&v)
		if err != nil {
			fmt.Println("Non-JSON message received or connection closed")
			return
		}
		data, ok := v.(map[string]interface{})
		if !ok {
			fmt.Println("Invalid data format received")
			return
		}

		if value, exists := data["type"]; exists && value == "send_message" {
			// Handle sending message logic here
		}
	}
}

func RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(60, 2)
	return func(c *gin.Context) {
		if c.Request.Method == "POST" && c.Request.URL.Path == "/register" {
			if limiter.Allow() {
				c.Next()
			} else {
				c.JSON(http.StatusTooManyRequests, gin.H{"message": "Rate limit exceeded"})
			}
		} else {
			c.Next()
		}
	}
}

type ValidTokenRequestBody struct {
	Token string `json:"token"`
}

type Config struct {
	MongoDBUrl string `json:"mongodb_url"`
	Port       int    `json:"port"`
	ResendKey  string `json:"resend_key"`
}

func loadConfig() (Config, error) {
	fileName := "config.json"
	if os.Getenv("AUTHENV") == "PRODUCTION" {
		if runtime.GOOS == "windows" {
			fileName = "C:\\etc\\uwu-auth\\config.json"
		} else {
			fileName = "/etc/uwu-auth/config.json"
		}
	}

	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Config file not found, creating a sample config")
			sampleConfig := Config{
				MongoDBUrl: "mongodb://127.0.0.1",
				Port:       80,
			}
			file, err = os.Create("config.json")
			if err != nil {
				return Config{}, fmt.Errorf("error creating config file: %w", err)
			}
			defer file.Close()

			jsonData, err := json.MarshalIndent(sampleConfig, "", "  ")
			if err != nil {
				return Config{}, fmt.Errorf("error marshalling JSON: %w", err)
			}
			_, err = file.Write(jsonData)
			if err != nil {
				return Config{}, fmt.Errorf("error writing to config file: %w", err)
			}
			return sampleConfig, nil
		}
		return Config{}, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return Config{}, fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	return config, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	db.CreateDBClient()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(RateLimiter())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	r.POST("/register", func(c *gin.Context) {
		routes.RegisterRoute(c, config.ResendKey)
	})
	r.POST("/login", routes.LoginRoute)
	r.POST("/profile/update-avatar", routes.UpdateProfilePicture)
	r.POST("/profile/update-bio", routes.UpdateBio)
	r.POST("/profile/update-status", routes.UpdateStatus)
	r.POST("/profile/update-pronouns", routes.UpdatePronouns)
	r.POST("/profile/update-banner", routes.UpdateProfileBanner)
	r.GET("/apps", routes.GetAppsList)
	r.GET("/app/:id", routes.GetAppByID)
	r.GET("/verify/:code", routes.VerifyAccount)
	r.GET("/role/:handle", routes.GetRole)
	r.GET("/profiles/:handle", routes.GetProfile)
	r.GET("/users/search", routes.SearchUsers)
	r.GET("/roles", routes.GetRolesList)
	r.POST("/create-app-association", routes.CreateAppAssociation)
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"version": 1})
	})
	r.POST("/promote", handlePromote)
	r.POST("/verify", handleVerify)
	r.POST("/unverify", handleUnverify)
	r.POST("/is-valid-token", handleValidToken)
	r.GET("/uploads/:upload", handleFileUpload)
	r.GET("/message", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})
	r.POST("/profile/post", routes.PostToProfileFeed)
	r.Run(fmt.Sprintf(":%d", config.Port))
}

func handlePromote(c *gin.Context) {
	password := c.Request.FormValue("password")
	if password != "i like feet" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid password"})
		return
	}
	handle := c.Request.FormValue("handle")
	role := c.Request.FormValue("role")
	roleInt, err := strconv.Atoi(role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid role"})
		return
	}
	filter := bson.D{{"handle", handle}}
	update := bson.D{{"$set", bson.D{{"role", roleInt}}}}

	collection := db.GetUsersCollection()
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

func handleVerify(c *gin.Context) {
	password := c.Request.FormValue("password")
	if password != "i like feet" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid password"})
		return
	}
	handle := c.Request.FormValue("handle")
	filter := bson.D{{"handle", handle}}
	update := bson.D{{"$set", bson.D{{"verified", true}}}}

	collection := db.GetUsersCollection()
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
}

func handleUnverify(c *gin.Context) {
	password := c.Request.FormValue("password")
	if password != "i like feet" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid password"})
		return
	}
	handle := c.Request.FormValue("handle")
	filter := bson.D{{"handle", handle}}
	update := bson.D{{"$set", bson.D{{"verified", false}}}}

	collection := db.GetUsersCollection()
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User unverified successfully"})
}

func handleValidToken(c *gin.Context) {
	var body ValidTokenRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
		return
	}
	filter := bson.D{{"token", body.Token}}
	usersCollection := db.GetUsersCollection()
	var result bson.M
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true})
}

func handleFileUpload(c *gin.Context) {
	fileName := c.Params.ByName("upload")
	c.File(fmt.Sprintf("./uploads/%s", fileName))
}
