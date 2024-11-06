package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Repository interface {
	SaveDevice(device Device) (Device, error)
	FindDeviceByID(id int) (Device, error)
}
type RepositoryImpl struct {
	db *sql.DB
}

func (r RepositoryImpl) FindDeviceByID(id int) (Device, error) {
	// query db for device with id

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
	// insert into db with creation time now
	query := "INSERT INTO devices (name, brand, creation_time) VALUES (?, ?, NOW())"
	result, err := r.db.Exec(query, device.Name, device.Brand)
	if err != nil {
		return Device{}, err
	}
	fmt.Println("Device saved successfully")
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