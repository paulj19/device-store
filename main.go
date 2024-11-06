package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Device struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Brand        string    `json:"brand"`
	CreationTime time.Time `json:"creation_time"`
}

var repository Repository

func main() {
	initDB()
	http.HandleFunc("/add-device", AddDeviceHandler)
}

func initDB() {

	// dsn := "user:password@tcp(127.0.0.1:3306)/"

	// // Connect to the MySQL server
	// db, err := sql.Open("mysql", dsn)
	// if err != nil {
	// 	log.Fatalf("Error connecting to the database: %v", err)
	// }
	// Create the database
	// _, err = db.Exec("CREATE DATABASE IF NOT EXISTS device_store")
	// if err != nil {
	// 	log.Fatalf("Error creating database: %v", err)
	// }
	// fmt.Println("Database created successfully")

	// // Connect to the newly created database
	// dsn = "user:password@tcp(127.0.0.1:3306)/device_store"
	// db, err = sql.Open("mysql", dsn)
	// if err != nil {
	// 	log.Fatalf("Error connecting to the device_store database: %v", err)
	// }

	// // Create a table
	// createTableSQL := `
	// CREATE TABLE IF NOT EXISTS devices (
	// id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
	// name VARCHAR(100) NOT NULL,
	// 			brand VARCHAR(100) NOT NULL,
	// creation_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	// );`
	// _, err = db.Exec(createTableSQL)
	// if err != nil {
	// 	log.Fatalf("Error creating table: %v", err)
	// }
	// fmt.Println("Table created successfully")
	dsn := "user:password@tcp(localhost:3306)/device_store"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = db.Ping()
	if err == nil {
		log.Println("Failed to connect to database", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	repository = RepositoryImpl{
		db: db,
	}
}

func AddDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var newDevice Device
	err := json.NewDecoder(r.Body).Decode(&newDevice)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("Request:")
	newDevice, err = repository.SaveDevice(newDevice)

	if err != nil {
		log.Println("Error adding device:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Device added:", newDevice)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newDevice)
}
