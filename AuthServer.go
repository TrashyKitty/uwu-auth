package main

// "github.com/Ant767/AuthBackend/utils"

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Ant767/AuthBackend/db"
	"github.com/Ant767/AuthBackend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(60, 2)
	return func(c *gin.Context) {
		if c.Request.Method == "POST" && c.Request.URL.Path == "/register" {
			if limiter.Allow() {
				c.Next()
			} else {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"message": "Rate limit exceed",
				})
			}
		} else {
			c.Next()
		}

	}
}

type Config struct {
	MongoDBUrl string `json:"mongodb_url"`
	Port       int    `json:"port"`
	ResendKey  string `json:"resend_key"`
}

func main() {

	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		fmt.Println("Creating sample config")
		sampleConfig := Config{
			MongoDBUrl: "mongodb://127.0.0.1",
			Port:       80,
		}

		// Open a file for writing (or create it if it doesn't exist)
		file, err := os.Create("config.json")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		// Marshal the struct into JSON format
		jsonData, err := json.MarshalIndent(sampleConfig, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}

		// Write the JSON data to the file
		_, err = file.Write(jsonData)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		return
	}
	defer file.Close()

	// Read the file contents
	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Unmarshal the JSON data into the struct
	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	db.CreateDBClient()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(RateLimiter())
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept"},
	}))
	r.POST("/register", func(c *gin.Context) {
		routes.RegisterRoute(c, config.ResendKey)
	})
	r.POST("/login", routes.LoginRoute)
	r.POST("/profile/update-avatar", routes.UpdateProfilePicture)
	r.POST("/profile/update-banner", routes.UpdateProfileBanner)
	r.GET("/apps", routes.GetAppsList)
	r.GET("/app/:id", routes.GetAppByID)
	r.GET("/verify/:code", routes.VerifyAccount)
	r.GET("/role/:handle", routes.GetRole)
	r.GET("/profile/:handle", routes.GetProfile)
	r.GET("/users/search", routes.SearchUsers)
	r.POST("/create-app-association", routes.CreateAppAssociation)
	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": 1})
	})
	r.Run(fmt.Sprintf(":%d", config.Port))
}
