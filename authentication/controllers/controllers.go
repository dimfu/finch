package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func checkDBMiddleware(ctx *gin.Context) (*pgxpool.Pool, error) {
	db, ok := ctx.MustGet("db").(*pgxpool.Pool)
	if !ok {
		return nil, errors.New("db not found in context")
	}
	return db, nil
}
