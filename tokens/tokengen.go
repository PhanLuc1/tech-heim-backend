package tokens

import "github.com/dgrijalva/jwt-go"

type SignedDetails struct {
	User_id int
	jwt.StandardClaims
}
