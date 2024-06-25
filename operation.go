package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type Customer struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, email FROM customers")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var customers []Customer
    for rows.Next() {
        var c Customer
        if err := rows.Scan(&c.ID, &c.Name, &c.Email); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        customers = append(customers, c)
    }

    json.NewEncoder(w).Encode(customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])

    var c Customer
    err := db.QueryRow("SELECT id, name, email FROM customers WHERE id = $1", id).Scan(&c.ID, &c.Name, &c.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Customer not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    json.NewEncoder(w).Encode(c)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
    var c Customer
    json.NewDecoder(r.Body).Decode(&c)

    err := db.QueryRow("INSERT INTO customers(name, email) VALUES($1, $2) RETURNING id", c.Name, c.Email).Scan(&c.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(c)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])

    var c Customer
    json.NewDecoder(r.Body).Decode(&c)

    _, err := db.Exec("UPDATE customers SET name = $1, email = $2 WHERE id = $3", c.Name, c.Email, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(c)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])

    _, err := db.Exec("DELETE FROM customers WHERE id = $1", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Customer with ID %d deleted", id)
}