package routes

import (
	"net/http"

	"github.com/Ant767/AuthBackend/auth"
	"github.com/Ant767/AuthBackend/db"
	"github.com/gin-gonic/gin"
)

type RegisterRequestBody struct {
	Username string `json:"username"`
	Handle   string `json:"handle"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterRoute(c *gin.Context, resendKey string) {
	var body RegisterRequestBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := auth.RegisterAccount(resendKey, db.GetUsersCollection(), body.Handle, body.Username, body.Password, body.Email)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": true, "message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": false, "message": "User successfully created!"})
	}
}
