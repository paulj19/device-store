package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	initDB()
	code := m.Run()
	os.Exit(code)
}

func Test_AddDeviceHandler(t *testing.T) {
	t.Run("should add device", func(t *testing.T) {

		deviceName := strconv.Itoa(rand.Intn(100)) + "TEST DEVICE"
		var device Device = Device{
			Name:  deviceName,
			Brand: "Test Brand",
		}
		deviceStub, err := json.Marshal(device)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/device", bytes.NewReader(deviceStub))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		fmt.Println("Request:", handler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
		if !strings.Contains(rr.Body.String(), deviceName) {
			t.Errorf("expected message %v, got %v", deviceName, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "Test Brand") {
			t.Errorf("expected message %v, got %v", "Test Brand", rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "creation_time") {
			t.Errorf("expected message %v, got %v", "creation_time", rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "id") {
			t.Errorf("expected message %v, got %v", "id", rr.Body.String())
		}
	})

	t.Run("should return 422 with message device already exists", func(t *testing.T) {
		deviceName := strconv.Itoa(rand.Intn(100)) + "TEST DEVICE"
		var device Device = Device{
			Name:  deviceName,
			Brand: "Test Brand",
		}
		deviceStub, err := json.Marshal(device)
		req, err := http.NewRequest("POST", "/device", bytes.NewReader(deviceStub))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)

		// Make another duplicate request
		deviceStub, err = json.Marshal(device)
		req, err = http.NewRequest("POST", "/device", bytes.NewReader(deviceStub))
		if err != nil {
			t.Fatal(err)
		}
		rr = httptest.NewRecorder()
		handler = http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, got %d", http.StatusUnprocessableEntity, rr.Code)
		}
		if strings.Contains(rr.Body.String(), fmt.Sprintf("Device %v already exists", deviceName)) {
			t.Errorf("expected message %v, got %v", fmt.Sprintf("Device %v already exists", device), rr.Body.String())
		}
	})
}

func Test_ReadDeviceHandler(t *testing.T) {
	t.Run("should return device", func(t *testing.T) {
		deviceName := strconv.Itoa(rand.Intn(100000)) + "TEST DEVICE"
		var device Device = Device{
			Name:  deviceName,
			Brand: "Test Brand",
		}
		device, err := repository.SaveDevice(device)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("GET", "/device?id="+strconv.Itoa(device.ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
		if !strings.Contains(rr.Body.String(), device.Name) {
			t.Errorf("expected message %v, got %v", device.Name, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "Test Brand") {
			t.Errorf("expected message %v, got %v", "Test Brand", rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "creation_time") {
			t.Errorf("expected message %v, got %v", "creation_time", rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "id") {
			t.Errorf("expected message %v, got %v", "id", rr.Body.String())
		}
	})
	t.Run("should return 404 not found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/device?id=100000", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "Device with id 100000 not found") {
			t.Errorf("expected message %v, got %v", "Device with id 100000 not found", rr.Body.String())
		}
	})
	t.Run("should return 400 bad request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/device?id=abc", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
	t.Run("should return 400 bad request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/device", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}
