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

func main() {
	initDB()
	http.HandleFunc("/device", CrudDeviceHandler)
	http.HandleFunc("/list-devices", GetAllDevicesHandler)
	http.HandleFunc("/search-device", SearchDeviceHandler)
}

func GetAllDevicesHandler(w http.ResponseWriter, r *http.Request) {
	devices, err := repository.FindAllDevices()
	if err != nil {
		log.Println("Error finding devices:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

func CrudDeviceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		path := strings.TrimPrefix(r.URL.Path, "/device/")
		deviceID, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid device ID", http.StatusBadRequest)
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

	case http.MethodPut:
		path := strings.TrimPrefix(r.URL.Path, "/device/")
		deviceID, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "Invalid device ID", http.StatusBadRequest)
			return
		}

		deviceFromDB, err := repository.FindDeviceByID(deviceID)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				http.Error(w, fmt.Sprintf("Device with id %v not found", deviceID), http.StatusNotFound)
				return
			}
			log.Println("Error finding device:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var deviceDTO Device
		err = json.NewDecoder(r.Body).Decode(&deviceDTO)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		deviceFromDB.Name = deviceDTO.Name
		deviceFromDB.Brand = deviceDTO.Brand

		_, err = repository.UpdateDevice(deviceFromDB)
		if err != nil {
			log.Println("Error updating device:", err)
			w.WriteHeader(http.StatusInternalServerError)
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
			log.Println("Error deleting device:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func SearchDeviceHandler(w http.ResponseWriter, r *http.Request) {
	brand := r.URL.Query().Get("brand")
	if brand == "" {
		http.Error(w, "Brand is required", http.StatusBadRequest)
		return
	}

	devices, err := repository.FindDeviceByBrand(brand)
	if err != nil {
		log.Println("Error finding devices:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}
