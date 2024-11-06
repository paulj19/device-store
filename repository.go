package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Repository interface {
	SaveDevice(device Device) (Device, error)
	FindDeviceByID(id int) (Device, error)
	FindDeviceByBrand(brand string) ([]Device, error)
	FindAllDevices() ([]Device, error)
	UpdateDevice(device Device) (Device, error)
	DeleteDevice(id int) error
	DeleteAllDevices()
}
type RepositoryImpl struct {
	db *sql.DB
}

func (r RepositoryImpl) FindDeviceByID(id int) (Device, error) {
	var creationTimeRaw []byte
	query := "SELECT * FROM devices WHERE id = ?"
	row := r.db.QueryRow(query, id)
	var device Device
	err := row.Scan(&device.ID, &device.Name, &device.Brand, &creationTimeRaw)

	if err != nil {
		return Device{}, err
	}

	creationTimeStr := string(creationTimeRaw)
	creationTime, err := time.Parse("2006-01-02 15:04:05", creationTimeStr)
	if err != nil {
		return Device{}, fmt.Errorf("error parsing creation_time: %v", err)
	}

	device.CreationTime = creationTime
	return device, nil
}

func (r RepositoryImpl) SaveDevice(device Device) (Device, error) {
	query := "INSERT INTO devices (name, brand, creation_time) VALUES (?, ?, NOW())"
	result, err := r.db.Exec(query, device.Name, device.Brand)
	if err != nil {
		return Device{}, err
	}
	deviceID, err := result.LastInsertId()
	if err != nil {
		return Device{}, err
	}

	device, err = r.FindDeviceByID(int(deviceID))
	if err != nil {
		return Device{}, err
	}
	return device, nil
}

func (respository RepositoryImpl) FindDeviceByBrand(brand string) ([]Device, error) {
	query := "SELECT * FROM devices WHERE brand = ?"
	rows, err := respository.db.Query(query, brand)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var device Device
		var creationTimeRaw []byte
		err := rows.Scan(&device.ID, &device.Name, &device.Brand, &creationTimeRaw)
		if err != nil {
			return nil, err
		}
		creationTimeStr := string(creationTimeRaw)
		creationTime, err := time.Parse("2006-01-02 15:04:05", creationTimeStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing creation_time: %v", err)
		}

		device.CreationTime = creationTime
		devices = append(devices, device)
	}
	return devices, nil
}

func (r RepositoryImpl) FindAllDevices() ([]Device, error) {
	query := "SELECT * FROM devices"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var device Device
		var creationTimeRaw []byte
		err := rows.Scan(&device.ID, &device.Name, &device.Brand, &creationTimeRaw)
		if err != nil {
			return nil, err
		}
		creationTimeStr := string(creationTimeRaw)
		creationTime, err := time.Parse("2006-01-02 15:04:05", creationTimeStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing creation_time: %v", err)
		}

		device.CreationTime = creationTime
		devices = append(devices, device)
	}
	return devices, nil
}

func (r RepositoryImpl) UpdateDevice(device Device) (Device, error) {
	query := "UPDATE devices SET name = ?, brand = ? WHERE id = ?"
	_, err := r.db.Exec(query, device.Name, device.Brand, device.ID)
	if err != nil {
		return Device{}, err
	}
	return device, nil
}

func (r RepositoryImpl) DeleteDevice(id int) error {
	_, err := r.FindDeviceByID(id)
	if err != nil {
		return err
	}
	query := "DELETE FROM devices WHERE id = ?"
	_, err = r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository RepositoryImpl) DeleteAllDevices() {
	query := "DELETE FROM devices"
	_, err := repository.db.Exec(query)
	if err != nil {
		log.Println("Error deleting devices:", err)
	}
}
