package database

import (
	"context"
	"encoding/base64"
	"fmt"
)

type User struct {
	Email          string
	PublicKeyEnc   string
	PrivateKeyEnc  string
	PublicKeySign  string
	PrivateKeySign string
	Totp           []byte
}

const (
	queryInsertUser = `
		INSERT INTO public."user"
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	queryUserByEmail = `
		SELECT *
		FROM public."user"
		WHERE email = $1
	`
)

func (db *DB) InsertUser(user *User) error {
	user.Totp = []byte(base64.StdEncoding.EncodeToString(user.Totp))
	_, err := db.pool.Exec(
		context.Background(),
		queryInsertUser,
		user.Email,
		user.PublicKeyEnc,
		user.PrivateKeyEnc,
		user.PublicKeySign,
		user.PrivateKeySign,
		user.Totp,
	)
	if err != nil {
		err = fmt.Errorf("inserting user: %w", err)
	}

	return err
}

func (db *DB) UserByEmail(email string) (User, error) {
	var user User
	err := db.pool.QueryRow(
		context.Background(),
		queryUserByEmail,
		email,
	).Scan(
		&user.Email,
		&user.PublicKeyEnc,
		&user.PrivateKeyEnc,
		&user.PublicKeySign,
		&user.PrivateKeySign,
		&user.Totp,
	)
	if err != nil {
		err = fmt.Errorf("retrieving user by email: %w", err)
	}
	user.Totp, _ = base64.StdEncoding.DecodeString(string(user.Totp))

	return user, err
}
