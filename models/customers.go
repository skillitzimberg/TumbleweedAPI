package models

import (
	"database/sql"
	"fmt"
)

// CustomerRecord gives form to customer data from form (no ID).
type CustomerRecord struct {
	FirstName  string `db:"first_name" json:"firstName"`
	LastName   string `db:"last_name" json:"lastName"`
	Phone      string `db:"phone" json:"phone"`
	Email      string `db:"email" json:"email"`
	PostalCode string `db:"postalCode" json:"postalCode"`
}

// Customer gives form to a customer's data returned from the database (includes ID).
type Customer struct {
	ID         int    `db:"id" json:"id"`
	FirstName  string `db:"first_name" json:"firstName"`
	LastName   string `db:"last_name" json:"lastName"`
	Phone      string `db:"phone" json:"phone"`
	Email      string `db:"email" json:"email"`
	PostalCode string `db:"postalCode" json:"postalCode"`
}

// AllCustomers returns all customers in the database.
func AllCustomers() ([]*Customer, error) {
	rows, err := db.Query("SELECT * FROM customers")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	customers := make([]*Customer, 0)
	for rows.Next() {
		customer := new(Customer)
		err := rows.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Phone, &customer.Email, &customer.PostalCode)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return customers, nil
}

// GetCustomer returns a customer from the database.
func GetCustomer(id int) (*Customer, error) {
	row := db.QueryRow("SELECT * FROM customers WHERE id=$1", id)

	customer := new(Customer)
	err := row.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Phone, &customer.Email, &customer.PostalCode)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

//AddCustomer add a customer to the database.
func AddCustomer(customerRecord CustomerRecord) (result sql.Result, err error) {
	fmt.Println("Got to AddCustomer")

	insertCustomer := fmt.Sprint("INSERT INTO customers (first_name, last_name, phone, email, postalCode) VALUES ($1, $2, $3, $4, $5)")

	result, err = db.Exec(insertCustomer, customerRecord.FirstName, customerRecord.LastName, customerRecord.Phone, customerRecord.Email, customerRecord.PostalCode)

	return
}

//EditCustomer edits a customer in the database using PUT.
func EditCustomer(customer Customer) (result sql.Result, err error) {
	editCustomer := fmt.Sprint("UPDATE customers SET first_name=$1, last_name=$2, phone=$3, email=$4, postalCode=$5 WHERE id = $6")

	result, err = db.Exec(editCustomer, customer.FirstName, customer.LastName, customer.Phone, customer.Email, customer.PostalCode, customer.ID)

	return
}

//DeleteCustomer deletes a customer from the database using DELETE.
func DeleteCustomer(id int) (result sql.Result, err error) {

	deleteCustomer := fmt.Sprint("DELETE FROM customers WHERE id = $1")

	result, err = db.Exec(deleteCustomer, id)

	return
}
