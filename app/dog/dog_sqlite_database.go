package dog

import (
	"database/sql"
	"fmt"
)

var _ DogDatabase = (*dogDatabase)(nil)

type dogDatabase struct {
	connection *sql.DB
}

func NewDogDatabase(connection *sql.DB) *dogDatabase {
	return &dogDatabase{connection: connection}
}

func (db *dogDatabase) Insert(dog Dog) error {
	_, err := db.connection.Exec("INSERT INTO dogs (id, name, breed) VALUES (?, ?, ?)", dog.ID, dog.Name, dog.Breed)
	if err != nil {
		return fmt.Errorf("failed to insert dogs, error: %w", err)
	}
	return nil
}

func (db *dogDatabase) SelectAll() ([]Dog, error) {
	rows, err := db.connection.Query(`SELECT id, name, breed FROM dogs;`)
	if err != nil {
		return nil, fmt.Errorf("failed to query dogs, error: %w", err)
	}
	defer rows.Close()

	dogs := []Dog{}
	for rows.Next() {
		var dog Dog
		if err := rows.Scan(&dog.ID, &dog.Name, &dog.Breed); err != nil {
			return nil, fmt.Errorf("failed to scan query result, error: %w", err)
		}

		dogs = append(dogs, dog)
	}

	return dogs, nil
}

func (db *dogDatabase) Delete(id string) error {
	_, err := db.connection.Exec("DELETE FROM dogs WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete dogs, error: %w", err)
	}

	return nil
}
