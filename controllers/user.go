package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm/clause"
)

var (
	appJSON     = "application/json"
	errEmail    = "ERROR: duplicate key value violates unique constraint \"idx_users_email\" (SQLSTATE 23505)"
	errUsername = "ERROR: duplicate key value violates unique constraint \"idx_users_username\" (SQLSTATE 23505)"
)

func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)

	User := models.User{}

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	if !(len(User.Password)>6){
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Password must be at least 6 characters",
		})
		return
	}

	hashedPass := helpers.HashPass(User.Password)
	User.Password = hashedPass

	err := db.Create(&User).Error

	if err != nil {
		message := err.Error()
		if message == errEmail || message == errUsername {
			message = "Username or Email Already Registered"
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": message,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       User.ID,
		"email":    User.Email,
		"age":      User.Age,
		"username": User.Username,
	})
}

func UserLogin(c *gin.Context) {
	db := database.GetDB()

	contentType := helpers.GetContentType(c)

	var password string

	var User models.User

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	password = User.Password

	err := db.Debug().Where("email = ?", User.Email).Take(&User).Error

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email/password (fetch)",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))

	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email/password",
		})
		return
	}

	token := helpers.GenerateToken(User.ID, User.Email, User.Username)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func UserUpdate(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("UserData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	var User models.User
	var currentUser models.User

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	currentUser.ID = uint(userData["id"].(float64))
	currentUser.Email = string(userData["email"].(string))
	currentUser.Username = string(userData["username"].(string))

	err := db.Model(&User).Clauses(clause.Returning{}).Where("id=?", currentUser.ID).Updates(models.User{
		Email:    User.Email,
		Username: User.Username,
	}).Error

	if err != nil {
		message := err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         User.ID,
		"email":      User.Email,
		"username":   User.Username,
		"age":        User.Age,
		"updated_at": User.UpdatedAt,
	})

}

func UserDelete(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("UserData").(jwt.MapClaims)

	userID := uint(userData["id"].(float64))
	var User models.User

	err := db.Model(&User).Delete(&User, userID).Error

	if err != nil {
		message := err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "your account has been sucessfully deleted",
	})

}
