# Weather by Zipcode with OpenTelemetry

This project consists of two microservices (`service_a` and `service_b`) implemented in Go, which use OpenTelemetry for tracing. The services communicate with each other to provide weather information based on a given zipcode.

## Project Structure

```plaintext
lab2-weather-by-zipcode-otel/
├── .docker/
│   └── otel-collector-config.yaml
├── service_a/
│   ├── Dockerfile
│   ├── go.mod
│   ├── main.go
├── service_b/
│   ├── Dockerfile
│   ├── go.mod
│   ├── main.go
├── shared/
│   ├── go.mod
│   ├── otel.go
├── docker-compose.yml
└── README.md
```

## Services
### service_a
service_a validates a provided zipcode and forwards the request to service_b. It is responsible for receiving the request, extracting the OpenTelemetry context, and initiating a new trace.

### service_b
service_b receives the forwarded request from service_a, retrieves the location for the provided zipcode, fetches the weather data, and responds back to service_a.

### shared
The shared directory contains common code used by both services, such as the OpenTelemetry setup.

## Running the Project
To run the project, use Docker Compose. This will build and start both services along with the necessary observability tools (Jaeger, Zipkin, and OpenTelemetry Collector).

### Prerequisites
- Docker
- Docker Compose

### Steps
1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/lab2-weather-by-zipcode-otel.git
    cd lab2-weather-by-zipcode-otel
    ```

2. Build and start the services:
    ```bash
    docker-compose -f docker-compose.yml up --build
    ```

3. Access Jaeger UI at http://localhost:16686

4. Access Zipkin UI at http://localhost:9411

## Testing
To test the setup, you can make a POST request to service_a to get the weather information for a given zipcode.

### Example Request
**Endpoint:** http://localhost:8080/weather
**Method:** POST
**Body:**

```json
{
    "zipcode": "80440000"
}
```

### Example Response
```json
{"temp_C":14.4,"temp_F":57.92,"temp_K":287.54999999999995}
```