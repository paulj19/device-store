# Device Store

Device Store is a simple RESTful API for managing devices. It allows you to create, read, update, and delete devices in a MySQL database.

### Design decisions:
   - I used relational DB because it seemed like a natural choice given the requirements, device name had clear direct relation with brand and this makes DB operations efficient
   - Device name and brand cannot be null or empty and request to create or set them so will fail. This ensures data consistency
   - combination of name and brand is unique and is set a composite key, this ensures no duplicates
   - get and search for list of devices are separate endpoints to ensure separation of concerns and to ensure response reflects single and list device output

## Features

- Add a new device
- Get a device by ID
- Get all devices
- Update a device
- Delete a device
- Search devices by brand

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/paulj19/device-store.git
   cd device-store
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```
   
## Tests
  To run all tests docker compose should be running, run the following commands
  ```sh 
  docker compose up -d 
  go test
  ```
  To run single test run the following command
  ```sh
  docker compose up -d 
  go test -run Test_AddDeviceHandler
  ```
  
  Once tests are done you can stop the docker compose
  ```sh
  docker compose down --volumes
  ```
## Usage

1. Run the server:

   ```sh
   go run main.go
   ```

2. The server will start on `http://localhost:8080`. You can use `curl` or any API client to interact with the API.

### Endpoints

- **Add a new device**

  ```sh
  curl -X POST -H "Content-Type: application/json" -d '{"name": "test device", "brand": "test brand"}' http://localhost:8080/device/
  ```

- **Get a device by ID**

  ```sh
  curl -X GET http://localhost:8080/device/{id}
  ```

- **Get all devices**

  ```sh
  curl -X GET http://localhost:8080/devices
  ```

- **Update a device**

  ```sh
  curl -X PUT -H "Content-Type: application/json" -d '{"name": "updated device", "brand": "updated brand"}' http://localhost:8080/device/{id}
  ```

- **Delete a device**

  ```sh
  curl -X DELETE http://localhost:8080/device/{id}
  ```

- **Search devices by brand**

  ```sh
  curl -X GET http://localhost:8080/devices?brand={brand}
  ```
