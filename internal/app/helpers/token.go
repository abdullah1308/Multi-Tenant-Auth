package helpers

import (
	"errors"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type SignedAccessDetails struct {
	UserID   uint
	Username string
	OrgID    uint
	UserType string
	jwt.StandardClaims
}

type SignedRefreshDetails struct {
	UserID   uint
	Username string
	OrgID    uint
	jwt.StandardClaims
}

var ACCESS_TOKEN_SECRET string = os.Getenv("ACCESS_TOKEN_SECRET")
var REFRESH_TOKEN_SECRET string = os.Getenv("REFRESH_TOKEN_SECRET")

func GenerateAllTokens(userID uint, username string, orgID uint, userType string) (signedToken string, signedRefreshToken string, err error) {
	accessClaims := &SignedAccessDetails{
		UserID:   userID,
		Username: username,
		OrgID:    orgID,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}

	refreshClaims := &SignedRefreshDetails{
		UserID:   userID,
		Username: username,
		OrgID:    orgID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(ACCESS_TOKEN_SECRET))

	if err != nil {
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(REFRESH_TOKEN_SECRET))

	if err != nil {
		return
	}

	return token, refreshToken, nil
}

func ValidateAccessToken(signedToken string) (claims *SignedAccessDetails, err error) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedAccessDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(ACCESS_TOKEN_SECRET), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*SignedAccessDetails)
	if !ok {
		err = errors.New("the token is invalid")
		return
	}

	if claims.ExpiresAt < time.Now().Unix() {
		err = errors.New("token is expired")
		return
	}

	return claims, err
}

func ValidateRefreshToken(signedToken string) (claims *SignedRefreshDetails, err error) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedRefreshDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(REFRESH_TOKEN_SECRET), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*SignedRefreshDetails)
	if !ok {
		err = errors.New("the token is invalid")
		return
	}

	if claims.ExpiresAt < time.Now().Unix() {
		err = errors.New("token is expired")
		return
	}

	return claims, err
}
