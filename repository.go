package main

import (
	"database/sql"
	"log"
)

type Repository interface {
	SaveDevice(device Device) (Device, error)
	FindDeviceByID(id int) (Device, error)
	FindDevicesByBrand(brand string) ([]Device, error)
	FindAllDevices() ([]Device, error)
	UpdateDevice(device Device) (Device, error)
	DeleteDevice(id int) error
	DeleteAllDevices()
}
type RepositoryImpl struct {
	db *sql.DB
}

func (r RepositoryImpl) FindDeviceByID(id int) (Device, error) {
	query := "SELECT * FROM devices WHERE id = ?"
	row := r.db.QueryRow(query, id)
	var device Device
	err := row.Scan(&device.ID, &device.Name, &device.Brand, &device.CreationTime)
	if err != nil {
		return Device{}, err
	}
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

func (r RepositoryImpl) FindDevicesByBrand(brand string) ([]Device, error) {
	query := "SELECT * FROM devices WHERE brand = ?"
	rows, err := r.db.Query(query, brand)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []Device
	for rows.Next() {
		var device Device
		err := rows.Scan(&device.ID, &device.Name, &device.Brand, &device.CreationTime)
		if err != nil {
			return nil, err
		}
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
		err := rows.Scan(&device.ID, &device.Name, &device.Brand, &device.CreationTime)
		if err != nil {
			return nil, err
		}
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

// helper function just for tests
func (r RepositoryImpl) DeleteAllDevices() {
	query := "DELETE FROM devices"
	_, err := r.db.Exec(query)
	if err != nil {
		log.Fatal("Error deleting devices:", err)
	}
}
