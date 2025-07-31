package models

import (
	"context"
	"strings"
	"time"

	"github.com/dimfu/finch/authentication/db"
	"github.com/guregu/null"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string      `json:"id,omitempty" db:"id"`
	Username     string      `json:"username" db:"username"`
	Email        string      `json:"email,omitempty" db:"email"`
	Password     string      `json:"password"` // used only internally, not stored
	PasswordHash null.String `json:"passwordHash,omitempty" db:"password_hash"`
	FirstName    null.String `json:"firstName,omitempty" db:"first_name"`
	LastName     null.String `json:"lastName,omitempty" db:"last_name"`
	CreatedAt    time.Time   `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt    time.Time   `json:"updatedAt,omitempty" db:"updated_at"`
}

func (u *User) ValidateCreds() map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(u.Username) == "" {
		errors["username"] = "Username is required"
	}
	if strings.TrimSpace(u.Password) == "" {
		errors["password"] = "Password is required"
	}

	return errors
}

func (u *User) FindByUsername() (*User, error) {
	db := db.Pool
	var user User
	ctx := context.Background()
	err := db.QueryRow(ctx, `
		SELECT id, username, email, password_hash 
		FROM users 
		WHERE username = $1`, u.Username).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)

	// in case we need it for authentication
	user.Password = u.Password

	if err != nil {
		// return the old user instance instead
		return u, err
	}
	return &user, err
}

func (u *User) CompareHashAndPassword() error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash.String), []byte(u.Password))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Create() error {
	db := db.Pool
	// TODO: more strict sanitization
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO users (username, email, first_name, last_name, password_hash)
		 VALUES ($1, $2, $3, $4, $5)`,
		u.Username, u.Email, u.FirstName, u.LastName, string(hash),
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
