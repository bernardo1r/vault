package database

import (
	"context"
	"fmt"
)

type ItemIndex struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

const queryIndexByEmail = `
	SELECT *
	FROM public."index"($1)
`

func (db *DB) IndexByEmail(email string) ([]ItemIndex, error) {
	rows, err := db.pool.Query(
		context.Background(),
		queryIndexByEmail,
		email,
	)
	if err != nil {
		return nil, fmt.Errorf("running 'index' function: %w", err)
	}

	index := make([]ItemIndex, 0)
	defer rows.Close()
	for rows.Next() {
		var item ItemIndex
		err = rows.Scan(
			&item.Id,
			&item.Label,
		)
		if err != nil {
			return nil, fmt.Errorf("reading index item value: %w", err)
		}
		index = append(index, item)
	}

	return index, nil
}
