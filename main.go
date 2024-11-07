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

func initDB() {
	dsn := "user:password@tcp(localhost:3306)/device_store?parseTime=true"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	repository = RepositoryImpl{
		db: db,
	}
}

func main() {
	initDB()
	http.HandleFunc("/device/", CrudDeviceHandler)
	http.HandleFunc("/devices", CrudDevicesHandler)
	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func CrudDeviceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		device, err := GetDeviceById(w, r)
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)

	case http.MethodPost:
		var newDevice Device
		err := json.NewDecoder(r.Body).Decode(&newDevice)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		if newDevice.Name == "" || newDevice.Brand == "" {
			http.Error(w, "Name and brand are required", http.StatusBadRequest)
			return
		}

		newDevice, err = repository.SaveDevice(newDevice)
		if err != nil {
			log.Printf("Error adding device: %v", err)
			if strings.Contains(err.Error(), "Error 1062 (23000): Duplicate entry") {
				http.Error(w, fmt.Sprintf("Device %v already exists", newDevice), http.StatusUnprocessableEntity)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("Device added: %v", newDevice)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newDevice)

	case http.MethodPut:
		deviceFromDB, err := GetDeviceById(w, r)
		if err != nil {
			return
		}
		var deviceDTO Device
		err = json.NewDecoder(r.Body).Decode(&deviceDTO)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if deviceDTO.Name == "" || deviceDTO.Brand == "" {
			http.Error(w, "Name and brand are required", http.StatusBadRequest)
			return
		}

		deviceFromDB.Name = deviceDTO.Name
		deviceFromDB.Brand = deviceDTO.Brand

		_, err = repository.UpdateDevice(deviceFromDB)
		if err != nil {
			log.Printf("Error updating device: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deviceFromDB)

	case http.MethodDelete:
		path := strings.TrimPrefix(r.URL.Path, "/device/")
		deviceID, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid device ID", http.StatusBadRequest)
			return
		}

		err = repository.DeleteDevice(deviceID)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				http.Error(w, fmt.Sprintf("Device with id %v not found", deviceID), http.StatusNotFound)
				return
			}
			log.Printf("Error deleting device: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func CrudDevicesHandler(w http.ResponseWriter, r *http.Request) {
	brand := r.URL.Query().Get("brand")
	var devices []Device
	var err error

	if brand == "" {
		devices, err = repository.FindAllDevices()
	} else {
		devices, err = repository.FindDevicesByBrand(brand)
	}

	if err != nil {
		log.Printf("Error finding devices: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func GetDeviceById(w http.ResponseWriter, r *http.Request) (Device, error) {
	path := strings.TrimPrefix(r.URL.Path, "/device/")
	deviceID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return Device{}, err
	}
	device, err := repository.FindDeviceByID(deviceID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, fmt.Sprintf("Device with id %v not found", deviceID), http.StatusNotFound)
			return Device{}, err
		}
		log.Printf("Error finding device: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return Device{}, err
	}
	return device, nil
}
