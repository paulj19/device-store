package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	http.HandleFunc("/device", CrudDeviceHandler)
}

func initDB() {
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

func CrudDeviceHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		deviceID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		device, err := repository.FindDeviceByID(deviceID)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				http.Error(w, fmt.Sprintf("Device with id %v not found", deviceID), http.StatusNotFound)
				return
			}
			log.Println("Error finding device:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)
	case http.MethodPost:
		var newDevice Device
		err := json.NewDecoder(r.Body).Decode(&newDevice)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newDevice, err = repository.SaveDevice(newDevice)

		if err != nil {
			log.Println("Error adding device:", err)
			if strings.Contains(err.Error(), "Error 1062 (23000): Duplicate entry") {
				http.Error(w, fmt.Sprintf("Device %v already exists", newDevice), http.StatusUnprocessableEntity)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("Device added:", newDevice)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newDevice)
	}
}
