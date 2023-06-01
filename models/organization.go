package models

import "database/sql"

type Organization struct {
	ID   int
	Name string
}

func CreateOrganization(db *sql.DB, name string) (Organization, error) {
	result, err := db.Exec("INSERT INTO Organization (name) VALUES (?)",
		name)

	if err != nil {
		return Organization{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Organization{}, err
	}
	return Organization{
		ID: int(id), Name: name,
	}, nil
}
