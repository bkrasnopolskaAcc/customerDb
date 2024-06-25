package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    _ "github.com/lib/pq"
    "github.com/gorilla/mux"
)

var db *sql.DB

func initDB() {
   // Connect to the default postgres database to create the new database
   defaultConnStr := "user=CarRental password=CarRental sslmode=disable"
   defaultDB, err := sql.Open("postgres", defaultConnStr)
   if err != nil {
	   log.Fatal(err)
   }
   defer defaultDB.Close()

   _, err = defaultDB.Exec("CREATE DATABASE customerdb")
   if err != nil && err.Error() != `pq: database "customerdb" already exists` {
	   log.Fatal(err)
   }

   // Connect to the newly created customerdb database
   connStr := "user=CarRental password=CarRental dbname=customerdb sslmode=disable"
   db, err = sql.Open("postgres", connStr)
   if err != nil {
	   log.Fatal(err)
   }

   err = db.Ping()
   if err != nil {
	   log.Fatal(err)
   }

   // Create customers table if it does not exist
   createTableQuery := `
   CREATE TABLE IF NOT EXISTS customers (
	   id SERIAL PRIMARY KEY,
	   name VARCHAR(100),
	   email VARCHAR(100)
   )`
   _, err = db.Exec(createTableQuery)
   if err != nil {
	   log.Fatal(err)
   }

   fmt.Println("Connected to the database successfully and ensured the customers table exists!")
}

func main() {
    initDB()
    defer db.Close()

    router := mux.NewRouter()
    router.HandleFunc("/customers", getCustomers).Methods("GET")
    router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
    router.HandleFunc("/customers", createCustomer).Methods("POST")
    router.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
    router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8000", router))
}