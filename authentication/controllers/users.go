package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dimfu/finch/authentication/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func Signup(ctx *gin.Context) {
	db, err := checkDBMiddleware(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		ctx.Error(err)
		return
	}

	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		ctx.Error(err)
		return

	}
	if err := user.Create(db); err != nil {
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
	db, err := checkDBMiddleware(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		ctx.Error(err)
		return
	}

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

	user, err = user.FindByUsername(db)
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

	log.Println(user)
}
