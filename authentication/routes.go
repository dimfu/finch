package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dimfu/finch/services/internal/jwt"
	"github.com/dimfu/finch/services/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func Signup(ctx *gin.Context) {
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		ctx.Error(err)
		return

	}
	if err := user.Create(); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "User already exist.",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func Signin(ctx *gin.Context) {
	user := &models.User{}
	if err := ctx.BindJSON(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		ctx.Error(err)
		return
	}

	if errs := user.ValidateCreds(); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"errors": errs,
		})
		return
	}

	user, err := user.FindByUsername()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("User with username %s not found", user.Username),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		ctx.Error(err)
		return
	}

	if err := user.CompareHashAndPassword(); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		ctx.Error(err)
		return
	}

	expirationTime := time.Now().Add(jwt.ACCESS_LIFETIME_DURATION)
	token, err := jwt.Generate(user.ID, expirationTime.Unix())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		ctx.Error(err)
		return
	}

	ctx.SetCookie(
		"access_token",
		token.AccessToken,
		int(time.Until(expirationTime).Seconds()), // maxAge in seconds
		"/",
		"",   // domain
		true, // secure (set to false if developing over HTTP)
		true, // httpOnly — important!
	)

	ctx.SetCookie(
		"refresh_token",
		token.RefreshToken,
		int(time.Until(time.Now().Add(jwt.REFRESH_LIFETIME_DURATION)).Seconds()), // maxAge in seconds
		"/",
		"",   // domain
		true, // secure (set to false if developing over HTTP)
		true, // httpOnly — important!
	)

	if err := token.ToRefreshToken().Insert(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Signed in successfully",
	})
}

func Refresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token not found",
		})
		return
	}
	user, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		ctx.Error(err)
		return
	}

	expirationTime := jwt.ACCESS_LIFETIME_DURATION
	token, err := jwt.Generate(user.ID, time.Now().Add(expirationTime).Unix())
	if err := token.ToRefreshToken().CreateOrUpdate(refreshToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		ctx.Error(err)
		return
	}

	ctx.SetCookie(
		"access_token",
		token.AccessToken,
		int(time.Until(time.Now().Add(expirationTime)).Seconds()), // maxAge in seconds
		"/",
		"",   // domain
		true, // secure (set to false if developing over HTTP)
		true, // httpOnly — important!
	)

	ctx.SetCookie(
		"refresh_token",
		token.RefreshToken, // This is the NEW refresh token that was stored in DB
		int(time.Until(time.Now().Add(jwt.REFRESH_LIFETIME_DURATION)).Seconds()), // maxAge in seconds
		"/",
		"",   // domain
		true, // secure (set to false if developing over HTTP)
		true, // httpOnly — important!
	)
}
