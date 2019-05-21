package models

import (
	"database/sql"
	"fmt"
)

type ProductOrder struct {
	ProductID       int
	QuantityOrdered int
}

// Product defines a product
type Order struct {
	ID               int `db:"id" json:"id"`
	CustomerID       int `db:"customer" json:"customer"`
	Customer         Customer
	ProductsOrdered  []ProductOrder `db:"products" json:"products"`
	PickupLocationID int            `db:"pickup_location" json:"pickup_location"`
	PickupDate       string         `db:"pickup_date" json:"pickup_date"`
	OrderTotal       float64        `db:"order_total" json:"order_total"`
}

// AllOrders returns all products in the database.
func AllOrders() ([]*Order, error) {
	rows, err := db.Query("SELECT * FROM orders, customers where customer.id = order.customerId")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	orders := make([]*Order, 0)
	for rows.Next() {
		order := new(Order)
		err := rows.Scan(&order.ID, &order.CustomerID, &order.Customer, &order.ProductsOrdered, &order.PickupLocationID, &order.PickupDate, &order.OrderTotal)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrder returns a orders from the database.
func GetOrder(id int) (*Order, error) {
	row := db.QueryRow("SELECT * FROM orders WHERE id=$1", id)

	order := new(Order)
	err := row.Scan(&order.ID, &order.CustomerID, &order.ProductsOrdered, &order.PickupLocationID, &order.PickupDate, &order.OrderTotal)
	if err != nil {
		return nil, err
	}
	return order, nil
}

//AddOrder add a orders to the database.
func AddOrder(order Order) (result sql.Result, err error) {

	insertOrder := fmt.Sprint("INSERT INTO orders (name, type, description, ingredients, price) VALUES ($1, $2, $3, $4, $5)")

	result, err = db.Exec(insertOrder, order.CustomerID, order.ProductsOrdered, order.PickupLocationID, order.PickupDate, order.OrderTotal)

	return
}

//EditOrder edits a orders in the database using PUT.
func EditOrder(order Order) (result sql.Result, err error) {
	editOrder := fmt.Sprint("UPDATE orders SET name=$1, type=$2, description=$3, ingredients=$4, price=$5 WHERE id = $6")

	result, err = db.Exec(editOrder, order.ProductsOrdered, order.PickupLocationID, order.PickupDate, order.OrderTotal, order.CustomerID)

	return
}

//DeleteOrder deletes a orders from the database using DELETE.
func DeleteOrder(id int) (result sql.Result, err error) {

	deleteOrder := fmt.Sprint("DELETE FROM orders WHERE id = $1")

	result, err = db.Exec(deleteOrder, id)

	return
}
