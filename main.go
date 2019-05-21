package main

import (
	"TumbleweedAPI/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host     = "tumbleweed-db"
	port     = 5432
	user     = "postgres"
	password = "docker"
	dbname   = "tumbleweed"
)

var psqlDatabaseConnectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

func main() {
	models.InitDB(psqlDatabaseConnectionString)

	http.HandleFunc("/customers", getCustomers)
	http.HandleFunc("/customers/find", getCustomer)
	http.HandleFunc("/customers/add", addCustomer)
	http.HandleFunc("/customers/edit", editCustomer)
	http.HandleFunc("/customers/delete", deleteCustomer)

	http.HandleFunc("/products", getProducts)
	http.HandleFunc("/products/find", getProduct)
	http.HandleFunc("/products/add", addProduct)
	http.HandleFunc("/products/edit", editProduct)
	http.HandleFunc("/products/delete", deleteProduct)

	http.HandleFunc("/locations", getLocations)
	http.HandleFunc("/locations/find", getLocation)
	http.HandleFunc("/locations/add", addLocation)
	http.HandleFunc("/locations/edit", editLocation)
	http.HandleFunc("/locations/delete", deleteLocation)

	http.ListenAndServe(":3000", nil)
}

// Customer functions
func getCustomers(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "GET")

	cstmrs, err := models.AllCustomers()
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, cstmrs)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "GET")

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	customerID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}

	cstmr, err := models.GetCustomer(customerID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, cstmr)
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "POST")

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	postalCode := r.FormValue("postalCode")
	if firstName == "" || lastName == "" || phone == "" || email == "" || postalCode == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	record := models.CustomerRecord{
		FirstName:  firstName,
		LastName:   lastName,
		Phone:      phone,
		Email:      email,
		PostalCode: postalCode,
	}

	rowsAffected, err := models.AddCustomer(record)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, rowsAffected)
}

func editCustomer(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "PUT")

	id := r.FormValue("id")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	postalCode := r.FormValue("postalCode")
	if id == "" || firstName == "" || lastName == "" || phone == "" || email == "" || postalCode == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	cstmrID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err.Error())
	}
	customer := models.Customer{
		ID:         cstmrID,
		FirstName:  firstName,
		LastName:   lastName,
		Phone:      phone,
		Email:      email,
		PostalCode: postalCode,
	}

	rowsAffected, err := models.EditCustomer(customer)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, rowsAffected)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "DELETE")

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	cstmrID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err.Error())
	}

	rowsAffected, err := models.DeleteCustomer(cstmrID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, rowsAffected)
}

// Product functions
func getProducts(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "GET")

	prdcts, err := models.AllProducts()
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, prdcts)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "GET")

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}

	prdct, err := models.GetProduct(productID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, prdct)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "POST")

	name := r.FormValue("name")
	productType := r.FormValue("type")
	description := r.FormValue("description")
	ingredients := r.FormValue("ingredients")
	price := r.FormValue("price")
	if name == "" || productType == "" || description == "" || ingredients == "" || price == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	productPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println(err.Error())
	}

	product := models.Product{
		Name:        name,
		Type:        productType,
		Description: description,
		Ingredients: ingredients,
		Price:       productPrice,
	}

	rowsAffected, err := models.AddProduct(product)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, rowsAffected)
}

func editProduct(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "PUT")

	id := r.FormValue("id")
	name := r.FormValue("name")
	productType := r.FormValue("type")
	description := r.FormValue("description")
	ingredients := r.FormValue("ingredients")
	price := r.FormValue("price")
	if id == "" || name == "" || productType == "" || description == "" || ingredients == "" || price == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err.Error())
	}

	productPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println(err.Error())
	}
	product := models.Product{
		ID:          productID,
		Name:        name,
		Type:        productType,
		Description: description,
		Ingredients: ingredients,
		Price:       productPrice,
	}

	rowsAffected, err := models.EditProduct(product)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, rowsAffected)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "DELETE")

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	prdctID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err.Error())
	}

	rowsAffected, err := models.DeleteProduct(prdctID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, rowsAffected)
}

// Locations functions
func getLocations(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "GET")

	prdcts, err := models.AllLocations()
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, prdcts)
}

func getLocation(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "GET")

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}

	prdct, err := models.GetLocation(productID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, prdct)
}

func addLocation(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "POST")

	name := r.FormValue("name")
	description := r.FormValue("description")
	address := r.FormValue("address")
	city := r.FormValue("city")
	state := r.FormValue("state")
	postalCode := r.FormValue("postalCode")
	if name == "" || description == "" || address == "" || city == "" || state == "" || postalCode == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	location := models.Location{
		Name:        name,
		Description: description,
		Address:     address,
		City:        city,
		State:       state,
		PostalCode:  postalCode,
	}

	rowsAffected, err := models.AddLocation(location)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, rowsAffected)
}

func editLocation(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "PUT")

	id := r.FormValue("id")
	name := r.FormValue("name")
	description := r.FormValue("description")
	address := r.FormValue("address")
	city := r.FormValue("city")
	state := r.FormValue("state")
	postalCode := r.FormValue("postalCode")
	if id == "" || name == "" || description == "" || address == "" || city == "" || state == "" || postalCode == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	locationID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err.Error())
	}

	location := models.Location{
		ID:          locationID,
		Name:        name,
		Description: description,
		Address:     address,
		City:        city,
		State:       state,
		PostalCode:  postalCode,
	}

	rowsAffected, err := models.EditLocation(location)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, rowsAffected)
}

func deleteLocation(w http.ResponseWriter, r *http.Request) {
	checkHTTPMethod(w, r, "DELETE")

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	locationID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err.Error())
	}

	rowsAffected, err := models.DeleteLocation(locationID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshalAndWriteJSON(w, rowsAffected)
}

func marshalAndWriteJSON(w http.ResponseWriter, objectToMarshal interface{}) {
	js, err := json.Marshal(objectToMarshal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func checkHTTPMethod(w http.ResponseWriter, r *http.Request, method string) {
	if r.Method != method {
		http.Error(w, http.StatusText(405), 405)
		return
	}
}
