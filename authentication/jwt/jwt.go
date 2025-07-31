package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dimfu/finch/authentication/models"
	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = []byte(os.Getenv("SECRET_KEY"))

var (
	ACCESS_LIFETIME_DURATION  = 1 * time.Hour
	REFRESH_LIFETIME_DURATION = 90 * 24 * time.Hour
	NO_SECRET_KEY             = errors.New("Error key environment is not set")
)

type Token struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

func (t *Token) ToRefreshToken() *models.RefreshToken {
	return &models.RefreshToken{
		UserID:    t.UserID,
		TokenHash: t.RefreshToken,
	}
}

func Generate(userID string, expTime int64) (Token, error) {
	token := Token{}
	var err error
	if len(SECRET_KEY) == 0 {
		return token, NO_SECRET_KEY
	}

	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = expTime

	// declare the token with the algorithm for signing, and the claims
	tokClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.AccessToken, err = tokClaims.SignedString(SECRET_KEY)

	token.UserID = userID

	if err != nil {
		return token, err
	}

	return createRefreshToken(token)
}

func ValidateAccessToken(accessToken string) (models.User, error) {
	user := models.User{}
	if len(SECRET_KEY) == 0 {
		return user, NO_SECRET_KEY
	}
	// validate whether the provided accessToken is using the same method or not
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return SECRET_KEY, nil
	})

	if err != nil {
		return user, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		user.ID = payload["user_id"].(string)
		return user, nil
	}

	return user, errors.New("invalid token")
}

func ValidateRefreshToken(refreshToken string) (models.User, error) {
	// do the same thing as access token validation but the purpose of doing this is
	// to get the payload from the refresh token such as; user_id and the expiration date
	user := models.User{}
	if len(SECRET_KEY) == 0 {
		return user, NO_SECRET_KEY
	}
	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return SECRET_KEY, nil
	})

	if err != nil {
		return user, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return user, errors.New("invalid token")
	}

	// reparse the access token without verifying its signature because we already know its already valid
	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	// payload["token"] is the access token that we want to parse
	token, _, err = parser.ParseUnverified(payload["token"].(string), claims)
	if err != nil {
		return user, err
	}

	payload, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return user, errors.New("invalid token")
	}

	user.ID = payload["user_id"].(string)

	// if its valid we're going to use user_id to refresh the refresh token later
	return user, nil
}

func createRefreshToken(token Token) (Token, error) {
	var err error

	claims := jwt.MapClaims{}
	// use the refresh token as access token
	claims["token"] = token.AccessToken
	claims["exp"] = time.Now().Add(REFRESH_LIFETIME_DURATION).Unix()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token.RefreshToken, err = refreshToken.SignedString(SECRET_KEY)
	if err != nil {
		return token, err
	}

	return token, nil
}
