package main

import (
    "database/sql"
    "encoding/json"
    "fmt" 
    "log"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var db *sql.DB
var err error

func main() {
    db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/ehe") //ganti ini
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    http.HandleFunc("/users", getUsers)
    fmt.Println("Server listening on :8080") 
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, email FROM users") //sesuaikan dengan tablenya dan data yang ingin diambil
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        users = append(users, user)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}
