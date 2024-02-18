# Log Handler

A simple log handler

## Setup and Usage

To set up the project, follow the steps below:

1. Clone the repository: `git clone https://github.com/MrRancy/LogAssignment.git`
2. Install dependencies: `go mod tidy`
3. Build the project: `go build`

### Running the Application

To run the application, execute the compiled binary:

```bash
./logAssignment
```

## API Endpoints

### Health Check

- **Endpoint**: `/healthz`
- **Method**: GET
- **Description**: Check the health status of the application.
- **Handler Function**: `Health`

### Log Endpoint

- **Endpoint**: `/log`
- **Method**: POST
- **Description**: Receive and process log payloads.
- **Handler Function**: `LogHandler`
- **Middleware**: `LogMiddleware`

#### Log Payload Format

The log endpoint expects the log payload to be in JSON format. The payload should adhere to the following structure:

```json
{
  "user_id": 1,
  "total": 1.65,
  "title": "delectus aut autem",
  "meta": {
    "logins": [
      {
        "time": "2020-08-08T01:52:50Z",
        "ip": "0.0.0.0"
      }
    ],
    "phone_numbers": {
      "home": "555-1212",
      "mobile": "123-5555"
    }
  },
  "completed": false
}
```

If the JSON payload is invalid, a `400 Bad Request` response will be returned with the error message.

### Logging Controller

The logging controller (`Controller`) provides the following functions:

#### Health Check Function

- **Function**: `Health`
- **Description**: Returns a simple "OK" response indicating the health of the application.

#### Log Handler Function

- **Function**: `LogHandler`
- **Description**: Handles incoming log payloads, validates the JSON format, and stores the logs in the cache.
- **Response Format**:

```json
{
  "Message": "Log payload received successfully",
  "Data": [
    // Array of logged data
  ]
}
```

## Middleware

### Log Middleware

- **Function**: `LogMiddleware`
- **Description**: Middleware function to handle logging-related tasks before reaching the `LogHandler`. This can include logging incoming requests or performing additional validations.

## Dependencies

List any external dependencies or libraries used in the project.

	github.com/gin-gonic/gin v1.9.1
	go.uber.org/zap v1.26.0
