package controller

import (
	"log"
	"regexp"

	"github.com/PhanLuc1/tech-heim-backend/database"
	"github.com/PhanLuc1/tech-heim-backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(400, gin.H{"message": err.Error()})
			return
		}
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		re := regexp.MustCompile(emailRegex)
		if !re.MatchString(*user.Email) {
			ctx.JSON(422, gin.H{"error": "Invalid email format"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		query := "INSERT INTO user (email, firstName, lastName, password) VALUES (?, ?, ?, ?)"
		_, err := database.Client.Query(query, user.Email, user.First_Name, user.Last_Name, user.Password)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(201, gin.H{"message": "Your account has been created"})
	}
}
