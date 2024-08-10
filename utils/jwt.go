package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"errors"
)

const secretKey="mySecretKey"
func GenerateToken(email string, userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"userId":userId,
		"exp": time.Now().Add(time.Hour*2).Unix(), //make token valid for only 2 hours
	})

	signedToken, err := token.SignedString([]byte(secretKey))
	return signedToken, err
}


func VerifyToken(token string) (string, error){

parsedToken, err := jwt.Parse(token, func(token *jwt.Token)(interface{}, error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC) //...HS256 is a type of HMAC
		
		//this checks for the signin methods to ensure it's same that was used during login
		if !ok{
			return nil, errors.New("Unexpected signing method")
		}

		return []byte(secretKey), nil
	})

	if err != nil{
		return "", errors.New("could not parse token: " + err.Error())
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid{
		return "", errors.New ("Invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok{
		return "", errors.New ("Invalid token claims")
	}

	// email:=claims["email"].(string)
	// userId:=int64(claims["userId"].(float64))
	userId:= claims["userId"].(string)

	return userId, nil
}
