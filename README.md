# go-crud-api

This project is a **RESTful API** built using Go, featuring a user entity with full CRUD operations. It is designed with
a focus on clean architecture and includes logging, configuration management, and MongoDB integration.

## **Features**

- **User Entity**:
    - Implements a user entity with a dedicated repository, service, and handler.
    - Supports **CRUD** operations (Create, Read, Update, Delete) for user data.

- **Logging**:
    - Utilizes **Logrus** for logging.
    - Logs are written both to the console and a file for persistent storage.

- **Configuration**:
    - Managed via **cleanenv**.
    - Configures both the server and MongoDB database connections.

- **MongoDB Integration**:
    - A MongoDB client is implemented to connect to the database and ping the Mongo server to ensure connectivity.

- **Error Handling**:
    - Custom error types are implemented to handle various application-specific errors.

- **Middleware**:
    - Middleware functions are used to enhance request handling

## **Getting Started**

### **Prerequisites**

- Go 1.22 or later
- Docker and Docker Compose
- MongoDB

1.**Clone the repository**:
   ```
   git clone https://github.com/yourusername/go-rest-api.git
   cd go-rest-api   
   ```
2.**Build the Docker Compose setup**:

```
docker-compose up --build
```

3.**Run the MongoDB service only**:
```
docker-compose up mongodb
```

4.**Run the application locally**:
```
go run cmd/main/app.go
```



