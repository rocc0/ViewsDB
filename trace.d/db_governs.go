package main

import (
	_ "github.com/lib/pq"
)

type government struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getGovernsList() (*[]government, error) {
	var (
		govs    []government
		govID   int
		govName string
	)
	res, err := db.Query("SELECT id, gov_name FROM governments")
	if err != nil {
		return nil, err
	}

	for res.Next() {

		if err = res.Scan(&govID, &govName); err != nil {
			return nil, err
		}

		govs = append(govs, government{govID, govName})
	}
	return &govs, nil
}
