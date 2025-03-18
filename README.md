## Overview  

This project is a **modular monolith** REST API for a social media platform where users can create posts and comment on them. The architecture is designed to be modular, allowing for easy scaling and future migration to microservices if needed.

### Key Features:  
- **Modular Design** – Each module encapsulates its own logic, making the system maintainable and scalable.  
- **RESTful API** – Provides endpoints for managing posts, comments, and user interactions.  
- **WebSocket Support** – Enables real-time updates and notifications.  
- **Flexible Communication** – Modules interact via interfaces (synchronous) and an event bus using Go channels (asynchronous).  
- **Scalability** – Designed to handle growing business requirements while maintaining performance.

## Usage  

The application can be run locally or using **Docker Compose**.

### Running the Application  

#### Option 1: Run Locally  
1. Clone the repository:  
   ```sh
   git clone https://github.com/vladovsiychuk/modular-monolith-go.git
   cd modular-monolith-go
   ```
2. Install dependencies
   ```sh
   go mod tidy
   ```
3. Start **Docker Desktop** app
4. Start the required databases using docker-compose
   ```sh
   docker-compose up postgres mongo redis -d
   ```
5. Run the application:
   ```sh
   go run ./cmd/main.go
   ```
#### Option 2: Run With Docker Compose
1. Build the Docker image (the image must be named app):
   ```sh
   docker build -t app .
   ```
2. Start the services with Docker Compose:
   ```sh
   docker-compose up
   ```

## Testing
1. Run only unit tests:
   ```sh
   go test ./...
   ```
2. Running unit and integration tests together (requires running **Docker Desktop**)
   ```sh
   go test -tags=integration ./...
   ```

