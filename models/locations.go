package models

import (
	"database/sql"
	"fmt"
)

// Location defines a location
type Location struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Address     string `db:"address" json:"address"`
	City        string `db:"city" json:"city"`
	State       string `db:"state" json:"state"`
	PostalCode  string `db:"postal_code" json:"postal_code"`
}

// AllLocations returns all locations in the database.
func AllLocations() ([]*Location, error) {
	rows, err := db.Query("SELECT * FROM locations")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	locations := make([]*Location, 0)
	for rows.Next() {
		location := new(Location)
		err := rows.Scan(&location.ID, &location.Name, &location.Description, &location.Address, &location.City, &location.State, &location.PostalCode)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return locations, nil
}

// GetLocation returns a locations from the database.
func GetLocation(id int) (*Location, error) {
	row := db.QueryRow("SELECT * FROM locations WHERE id=$1", id)

	location := new(Location)
	err := row.Scan(&location.ID, &location.Name, &location.Description, &location.Address, &location.City, &location.State, &location.PostalCode)
	if err != nil {
		return nil, err
	}
	return location, nil
}

//AddLocation add a locations to the database.
func AddLocation(location Location) (result sql.Result, err error) {

	insertLocation := fmt.Sprint("INSERT INTO locations (name, type, description, ingredients, price) VALUES ($1, $2, $3, $4, $5)")

	result, err = db.Exec(insertLocation, location.ID, location.Name, location.Description, location.Address, location.City, location.State, location.PostalCode)

	return
}

//EditLocation edits a locations in the database using PUT.
func EditLocation(location Location) (result sql.Result, err error) {
	editLocation := fmt.Sprint("UPDATE locations SET name=$1, description=$2, address=$3, city=$4, state=$5, postal_code=$6 WHERE id = $7")

	result, err = db.Exec(editLocation, location.Name, location.Description, location.Address, location.City, location.State, location.PostalCode, location.ID)

	return
}

//DeleteLocation deletes a locations from the database using DELETE.
func DeleteLocation(id int) (result sql.Result, err error) {

	deleteLocation := fmt.Sprint("DELETE FROM locations WHERE id = $1")

	result, err = db.Exec(deleteLocation, id)

	return
}
