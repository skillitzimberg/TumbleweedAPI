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
	http.ListenAndServe(":3000", nil)
}

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
