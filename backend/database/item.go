package database

import (
	"context"
	"fmt"
)

type Item struct {
	Email      string
	Id         string
	Label      string
	Key        string
	Credential string
}

const (
	queryInsertItem = `
		INSERT INTO public.item
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (email, id) DO UPDATE
			SET "label" = EXCLUDED."label",
			    "key" = EXCLUDED."key",
			    credential = EXCLUDED.Credential
	`
	queryItemByPK = `
		SELECT *
		FROM public.item
		WHERE email = $1 AND
		      id = $2
	`
	queryDeleteItem = `
		DELETE FROM public.item
		WHERE email = $1 AND
		      id = $2
	`
)

func (db *DB) InsertItem(item *Item) error {
	_, err := db.pool.Exec(
		context.Background(),
		queryInsertItem,
		item.Email,
		item.Id,
		item.Label,
		item.Key,
		item.Credential,
	)
	if err != nil {
		err = fmt.Errorf("inserting item: %w", err)
	}

	return err
}

func (db *DB) ItemByPK(email string, id string) (Item, error) {
	var item Item
	err := db.pool.QueryRow(
		context.Background(),
		queryItemByPK,
		email,
		id,
	).Scan(
		&item.Email,
		&item.Id,
		&item.Label,
		&item.Key,
		&item.Credential,
	)
	if err != nil {
		err = fmt.Errorf("retrieving item by pk: %w", err)
	}

	return item, err
}

func (db *DB) DeleteItem(email string, id string) error {
	_, err := db.pool.Exec(
		context.Background(),
		queryDeleteItem,
		email,
		id,
	)
	if err != nil {
		err = fmt.Errorf("deleting item: %w", err)
	}

	return err
}
