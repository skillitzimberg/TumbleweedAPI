package models

import (
	"database/sql"
	"fmt"
)

// Product defines a product
type Product struct {
	ID          int     `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Type        string  `db:"productType" json:"productType"`
	Description string  `db:"description" json:"description"`
	Ingredients string  `db:"ingredients" json:"ingredients"`
	Price       float64 `db:"price" json:"price"`
}

// AllProducts returns all products in the database.
func AllProducts() ([]*Product, error) {
	rows, err := db.Query("SELECT * FROM products")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	products := make([]*Product, 0)
	for rows.Next() {
		product := new(Product)
		err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Description, &product.Ingredients, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// GetProduct returns a products from the database.
func GetProduct(id int) (*Product, error) {
	row := db.QueryRow("SELECT * FROM products WHERE id=$1", id)

	product := new(Product)
	err := row.Scan(&product.ID, &product.Name, &product.Type, &product.Description, &product.Ingredients, &product.Price)
	if err != nil {
		return nil, err
	}
	return product, nil
}

//AddProduct add a products to the database.
func AddProduct(product Product) (result sql.Result, err error) {

	insertProduct := fmt.Sprint("INSERT INTO products (name, type, description, ingredients, price) VALUES ($1, $2, $3, $4, $5)")

	result, err = db.Exec(insertProduct, product.Name, product.Type, product.Description, product.Ingredients, product.Price)

	return
}

//EditProduct edits a products in the database using PUT.
func EditProduct(product Product) (result sql.Result, err error) {
	editProduct := fmt.Sprint("UPDATE products SET name=$1, type=$2, description=$3, ingredients=$4, price=$5 WHERE id = $6")

	result, err = db.Exec(editProduct, product.Name, product.Type, product.Description, product.Ingredients, product.Price, product.ID)

	return
}

//DeleteProduct deletes a products from the database using DELETE.
func DeleteProduct(id int) (result sql.Result, err error) {

	deleteProduct := fmt.Sprint("DELETE FROM products WHERE id = $1")

	result, err = db.Exec(deleteProduct, id)

	return
}
