package main

// "github.com/Ant767/AuthBackend/utils"

import (
	"github.com/Ant767/AuthBackend/db"
	"github.com/Ant767/AuthBackend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.CreateDBClient()
	r := gin.Default()
	r.POST("/register", routes.RegisterRoute)
	r.Run()
}
