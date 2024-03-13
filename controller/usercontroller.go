package controller

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/PhanLuc1/tech-heim-backend/database"
	"github.com/PhanLuc1/tech-heim-backend/models"
	generate "github.com/PhanLuc1/tech-heim-backend/tokens"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var refreshTokenMap = make(map[string]bool)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Passowrd is Incorerct"
		valid = false
	}
	return valid, msg
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
func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User
		var founduser models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(400, gin.H{"message": err.Error()})
			return
		}
		fmt.Println(*user.Email)
		err := database.Client.QueryRow("SELECT * FROM user WHERE email = ?", user.Email).Scan(
			&founduser.Id,
			&founduser.Email,
			&founduser.First_Name,
			&founduser.Last_Name,
			&founduser.Password,
			&founduser.PhoneNumber,
			&founduser.Address,
		)
		if err != nil {
			ctx.JSON(500, gin.H{"Message": "Email is not avilable"})
			return
		}
		PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)
		if !PasswordIsValid {
			ctx.JSON(401, gin.H{"error": msg})
			return
		}
		token, refreshToken, _ := generate.TokenGeneration(*founduser.Id)
		refreshTokenMap[refreshToken] = true
		ctx.IndentedJSON(http.StatusOK, gin.H{"token": token, "refreshToken": refreshToken})
	}
}
