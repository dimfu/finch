package models

import (
	"context"

	"github.com/dimfu/finch/authentication/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RefreshToken struct {
	TokenHash string `db:"token_hash"`
	UserID    string `db:"user_id"`
	Revoked   bool   `db:"revoked"`
}

func (r *RefreshToken) Insert() error {
	db := db.Pool
	ctx := context.Background()
	if _, err := db.Exec(ctx,
		`INSERT INTO refresh_tokens (token_hash, user_id, revoked) 
		VALUES ($1, $2, $3)`,
		r.TokenHash, r.UserID, false,
	); err != nil {
		return err
	}
	return nil
}

func (r *RefreshToken) RevokeByHash() error {
	db := db.Pool
	ctx := context.Background()
	if _, err := db.Exec(ctx,
		`UPDATE refresh_tokens SET revoked = TRUE WHERE token_hash = $1`, r.TokenHash,
	); err != nil {
		return err
	}
	return nil
}

// TODO: Add delete refresh tuple method once we have the session list ui for it

func (r *RefreshToken) CreateOrUpdate(prevToken string) error {
	db := db.Pool
	ctx := context.Background()
	var rowId uuid.UUID
	err := db.QueryRow(ctx, "SELECT id FROM refresh_tokens WHERE token_hash = $1 AND revoked = false", prevToken).Scan(&rowId)
	switch err {
	case nil:
		_, err := db.Exec(ctx, "UPDATE refresh_tokens SET token_hash = $1, updated_at = now() WHERE id = $2", r.TokenHash, rowId.String())
		return err
	case pgx.ErrNoRows:
		return r.Insert()
	default:
		return err
	}
}
