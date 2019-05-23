package models

import (
	"database/sql"
	"fmt"
)

// OrderHeader collects basic order information
type OrderHeader struct {
	OrderID          int    `db:"id" json:"id"`
	OrderDate        string `db:"order_date" json:"order_date"`
	PickupLocationID int    `db:"pickup_location" json:"pickup_location"`
	PickupDate       string `db:"pickup_date" json:"pickup_date"`
	CustomerID       int    `db:"customer_id" json:"customer_id"`
}

// A ProductOrdered collects product information for association with and order
type ProductOrdered struct {
	OrderID         int     `db:"order_id" json:"order_id"`
	ProductID       int     `db:"product_id" json:"product_id"`
	ProductName     string  `db:"name" json:"name"`
	ProductPrice    float64 `db:"price" json:"price"`
	QuantityOrdered int     `db:"quantity" json:"quantity"`
}

// OrderTotal holds total order amount
type OrderTotal struct {
	Total float64
}

// Order defines a product
type Order struct {
	OrderHeader     OrderHeader
	Customer        Customer
	Location        *Location
	ProductsOrdered []ProductOrdered
	Total           float64
}

func createOrderHeaders() ([]*OrderHeader, error) {
	rows, err := db.Query("SELECT orders.id, orders.pickup_location_id, orders.pickup_date, orders.customer_id FROM orders;")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	orderHeaders := make([]*OrderHeader, 0)
	for rows.Next() {
		orderHeader := new(OrderHeader)

		err = rows.Scan(&orderHeader.OrderID, &orderHeader.PickupLocationID, &orderHeader.PickupDate, &orderHeader.CustomerID)
		if err != nil {
			return nil, err
		}
		orderHeaders = append(orderHeaders, orderHeader)
	}
	return orderHeaders, nil
}

func getAllOrdersProducts() ([]*ProductOrdered, error) {
	rows, err := db.Query("SELECT products_orders.order_id, products_orders.product_id, products.name, products.price, products_orders.quantity FROM orders LEFT JOIN products_orders ON orders.id = products_orders.order_id LEFT JOIN products ON products_orders.product_id = products.id;")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	productsOrdered := make([]*ProductOrdered, 0)

	for rows.Next() {
		productOrdered := new(ProductOrdered)
		err = rows.Scan(&productOrdered.OrderID, &productOrdered.ProductID, &productOrdered.ProductName, &productOrdered.ProductPrice, &productOrdered.QuantityOrdered)
		if err != nil {
			return nil, err
		}
		productsOrdered = append(productsOrdered, productOrdered)
	}

	return productsOrdered, nil
}

func getCustomersWithOrders() ([]*Customer, error) {
	rows, err := db.Query("SELECT customers.id, customers.first_name, customers.last_name, customers.phone, customers.email, customers.postalcode FROM customers JOIN orders ON orders.customer_id = customers.id;")

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

	return customers, nil
}

func calculateOrderTotal(productsOnOrder []ProductOrdered) (orderTotal float64) {

	for i := 0; i < len(productsOnOrder); i++ {
		orderTotal += productsOnOrder[i].ProductPrice * float64(productsOnOrder[i].QuantityOrdered)
	}
	return
}

// AllOrders returns all orders in the database.
func AllOrders() ([]*Order, error) {
	orderHeaders, err := createOrderHeaders()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	orderedProducts, err := getAllOrdersProducts()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	customers, err := getCustomersWithOrders()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	orders := make([]*Order, 0)
	for i := 0; i < len(orderHeaders); i++ {
		order := new(Order)
		order.OrderHeader = *orderHeaders[i]

		for c := 0; c < len(customers); c++ {
			if customers[c].ID == orderHeaders[i].CustomerID {
				order.Customer = *customers[c]
			}
		}

		order.Location, err = GetLocation(order.OrderHeader.PickupLocationID)
		if err != nil {
			return nil, err
		}

		for j := 0; j < len(orderedProducts); j++ {
			if orderHeaders[i].OrderID == orderedProducts[j].OrderID {
				order.ProductsOrdered = append(order.ProductsOrdered, *orderedProducts[j])
			}
			order.Total = calculateOrderTotal(order.ProductsOrdered)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// GetOrder returns an order from the database.
func GetOrder(id int) (*Order, error) {
	orderHeaderRow := db.QueryRow("SELECT orders.id, orders.pickup_location_id, orders.pickup_date, orders.customer_id FROM orders WHERE id=$1", id)

	productRows, err := db.Query("SELECT products_orders.order_id, products_orders.product_id, products.name, products.price, products_orders.quantity FROM products JOIN products_orders ON products.id = products_orders.product_id WHERE products_orders.order_id=$1;", id)
	if err != nil {
		return nil, err
	}
	orderHeader := new(OrderHeader)
	err = orderHeaderRow.Scan(&orderHeader.OrderID, &orderHeader.PickupLocationID, &orderHeader.PickupDate, &orderHeader.CustomerID)
	if err != nil {
		return nil, err
	}

	products := make([]ProductOrdered, 0)
	for productRows.Next() {
		product := new(ProductOrdered)
		err = productRows.Scan(&product.OrderID, &product.ProductID, &product.ProductName, &product.ProductPrice, &product.QuantityOrdered)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}

	customer, err := GetCustomer(orderHeader.CustomerID)
	if err != nil {
		return nil, err
	}
	total := calculateOrderTotal(products)

	order := new(Order)
	order.OrderHeader = *orderHeader
	order.Customer = *customer
	order.Location, err = GetLocation(orderHeader.PickupLocationID)
	if err != nil {
		return nil, err
	}
	order.ProductsOrdered = products
	order.Total = total
	return order, nil
}

//AddOrder add a orders to the database.
// func AddOrder(order Order) (result sql.Result, err error) {

// 	insertOrder := fmt.Sprint("INSERT INTO orders (customer_id, order_date, location_id, pickup_date) VALUES ($1, $2, $3, $4, $5)")

// 	result, err = db.Exec(insertOrder, order.CustomerID, order.ProductsOrdered, order.PickupLocationID, order.PickupDate, order.OrderTotal)

// 	return
// }

//EditOrder edits a orders in the database using PUT.
// func EditOrder(order Order) (result sql.Result, err error) {
// 	editOrder := fmt.Sprint("UPDATE orders SET name=$1, type=$2, description=$3, ingredients=$4, price=$5 WHERE id = $6")

// 	result, err = db.Exec(editOrder, order.ProductsOrdered, order.PickupLocationID, order.PickupDate, order.OrderTotal, order.CustomerID)

// 	return
// }

//DeleteOrder deletes a orders from the database using DELETE.
func DeleteOrder(id int) (result sql.Result, err error) {

	deleteOrder := fmt.Sprint("DELETE FROM orders WHERE id = $1")

	result, err = db.Exec(deleteOrder, id)

	return
}
