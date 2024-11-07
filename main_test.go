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

	defer repository.DeleteAllDevices()
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
		handler.ServeHTTP(rr, req)
		var deviceResponse Device
		err = json.Unmarshal(rr.Body.Bytes(), &deviceResponse)
		if err != nil {
			t.Fatal(err)
		}

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
		if deviceResponse.Name != device.Name {
			t.Errorf("expected name %v, got %v", device.Name, deviceResponse.Name)
		}
		if deviceResponse.Brand != device.Brand {
			t.Errorf("expected brand %v, got %v", device.Brand, deviceResponse.Brand)
		}
		if deviceResponse.ID == 0 {
			t.Errorf("expected id to be non zero, got %v", deviceResponse.ID)
		}
		if deviceResponse.CreationTime.IsZero() {
			t.Errorf("expected creation time to be non zero, got %v", deviceResponse.CreationTime)
		}
	})

	t.Run("should return 422 with message device already exists", func(t *testing.T) {
		deviceName := strconv.Itoa(rand.Intn(100000)) + "TEST DEVICE"
		var device Device = Device{
			Name:  deviceName,
			Brand: "Test Brand",
		}
		device, err := repository.SaveDevice(device)
		if err != nil {
			t.Fatal(err)
		}
		// Make another duplicate request
		deviceStub, err := json.Marshal(device)
		req, err := http.NewRequest("POST", "/device", bytes.NewReader(deviceStub))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, got %d", http.StatusUnprocessableEntity, rr.Code)
		}
		if strings.Contains(rr.Body.String(), fmt.Sprintf("Device %v already exists", deviceName)) {
			t.Errorf("expected message %v, got %v", fmt.Sprintf("Device %v already exists", device), rr.Body.String())
		}
	})
	t.Run("should return 400 bad request for emtpy name or brand", func(t *testing.T) {
		var device Device = Device{
			Name:  "",
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
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "Name and brand are required") {
			t.Errorf("expected message %v, got %v", "Name and brand are required", rr.Body.String())
		}
	})

	repository.DeleteAllDevices()
}

func Test_GetDeviceHandler(t *testing.T) {
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

		req, err := http.NewRequest("GET", "/device/"+strconv.Itoa(device.ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)
		var deviceResponse Device
		err = json.Unmarshal(rr.Body.Bytes(), &deviceResponse)
		if err != nil {
			t.Fatal(err)
		}
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
		if deviceResponse.Name != device.Name {
			t.Errorf("expected name %v, got %v", device.Name, deviceResponse.Name)
		}
		if deviceResponse.Brand != device.Brand {
			t.Errorf("expected brand %v, got %v", device.Brand, deviceResponse.Brand)
		}
		if deviceResponse.ID != device.ID {
			t.Errorf("expected id %v, got %v", device.ID, deviceResponse.ID)
		}
		if deviceResponse.CreationTime.IsZero() {
			t.Errorf("expected creation time to be non zero, got %v", deviceResponse.CreationTime)
		}
	})
	t.Run("should return 404 not found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/device/100000", nil)
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
		req, err := http.NewRequest("GET", "/device/abc", nil)
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
	repository.DeleteAllDevices()
}

func Test_ListDevicesHandler(t *testing.T) {
	repository.DeleteAllDevices()
	t.Run("should return list of devices", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			deviceName := strconv.Itoa(rand.Intn(100000)) + "TEST DEVICE"
			var device Device = Device{
				Name:  deviceName,
				Brand: "Test Brand",
			}
			device, err := repository.SaveDevice(device)
			if err != nil {
				t.Fatal(err)
			}
		}
		req, err := http.NewRequest("GET", "/list-devices", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAllDevicesHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
		var devices []Device
		err = json.Unmarshal(rr.Body.Bytes(), &devices)
		if err != nil {
			t.Fatal(err)
		}
		if len(devices) != 10 {
			t.Errorf("expected 10 devices, got %d", len(devices))
		}
		repository.DeleteAllDevices()
	})
	t.Run("should return empty list of devices", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/list-devices", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAllDevicesHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
		var devices []Device
		err = json.Unmarshal(rr.Body.Bytes(), &devices)
		if err != nil {
			t.Fatal(err)
		}
		if len(devices) != 0 {
			t.Errorf("expected 0 devices, got %d", len(devices))
		}
	})

	repository.DeleteAllDevices()
}

func Test_UpdateDevice(t *testing.T) {
	t.Run("should update device", func(t *testing.T) {
		deviceName := strconv.Itoa(rand.Intn(100000)) + "TEST DEVICE"
		var device Device = Device{
			Name:  deviceName,
			Brand: "Test Brand",
		}
		device, err := repository.SaveDevice(device)
		if err != nil {
			t.Fatal(err)
		}
		device.Name = "Updated Name"
		device.Brand = "Updated Brand"
		deviceStub, err := json.Marshal(device)
		req, err := http.NewRequest("PUT", "/device/"+strconv.Itoa(device.ID), bytes.NewReader(deviceStub))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var deviceResponse Device
		err = json.Unmarshal(rr.Body.Bytes(), &deviceResponse)
		if err != nil {
			t.Fatal(err)
		}
		if deviceResponse.Name != device.Name {
			t.Errorf("expected name %v, got %v", device.Name, deviceResponse.Name)
		}
		if deviceResponse.Brand != device.Brand {
			t.Errorf("expected brand %v, got %v", device.Brand, deviceResponse.Brand)
		}
		if deviceResponse.ID != device.ID {
			t.Errorf("expected id %v, got %v", device.ID, deviceResponse.ID)
		}
		if deviceResponse.CreationTime.IsZero() {
			t.Errorf("expected creation time to be non zero, got %v", deviceResponse.CreationTime)
		}
	})

	t.Run("should return 404 not found", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/device/100000", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("should return 400 bad request for emtpy name or brand", func(t *testing.T) {
		var device Device = Device{
			Name:  "Test Device",
			Brand: "Test Brand",
		}
		device, err := repository.SaveDevice(device)
		if err != nil {
			t.Fatal(err)
		}
		device.Name = ""
		device.Brand = ""
		deviceStub, err := json.Marshal(device)
		req, err := http.NewRequest("PUT", "/device/"+strconv.Itoa(device.ID), bytes.NewReader(deviceStub))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "Name and brand are required") {
			t.Errorf("expected message %v, got %v", "Name and brand are required", rr.Body.String())
		}
	})
	repository.DeleteAllDevices()
}

func Test_DeleteDevice(t *testing.T) {
	t.Run("should delete device", func(t *testing.T) {
		deviceName := strconv.Itoa(rand.Intn(100000)) + "TEST DEVICE"
		var device Device = Device{
			Name:  deviceName,
			Brand: "Test Brand",
		}
		device, err := repository.SaveDevice(device)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("DELETE", "/device/"+strconv.Itoa(device.ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
		_, err = repository.FindDeviceByID(device.ID)
		if err == nil || !strings.Contains(err.Error(), "no rows in result set") {
			t.Errorf("expected device to be deleted, got %v", err)
		}
	})

	t.Run("should return 404 not found", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/device/1111111", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CrudDeviceHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}
	})
	repository.DeleteAllDevices()
}

func Test_SearchDeviceHandler(t *testing.T) {
	t.Run("should search device", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			var device Device = Device{
				Name:  strconv.Itoa(rand.Intn(100000)) + "TEST DEVICE",
				Brand: "Test Brand",
			}
			device, err := repository.SaveDevice(device)
			if err != nil {
				t.Fatal(err)
			}
		}
		req, err := http.NewRequest("GET", "/search-device?brand=TEST DEVICE", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetAllDevicesHandler)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
		var devices []Device
		err = json.Unmarshal(rr.Body.Bytes(), &devices)
		if err != nil {
			t.Fatal(err)
		}
		if len(devices) != 10 {
			t.Errorf("expected 10 devices, got %d", len(devices))
		}
		repository.DeleteAllDevices()
	})

	t.Run("should return 404 not found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/search-device?brand=ABC", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(SearchDeviceHandler)
		handler.ServeHTTP(rr, req)
		var devices []Device
		err = json.Unmarshal(rr.Body.Bytes(), &devices)
		if err != nil {
			t.Fatal(err)
		}
		if len(devices) != 0 {
			t.Errorf("expected 0 devices, got %d", len(devices))
		}
	})
	repository.DeleteAllDevices()
}
