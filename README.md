# Fiber HRMS (Human Resource Management System)

Fiber HRMS is a lightweight Human Resource Management System built using the Fiber web framework for Go and MongoDB as the database backend. It provides endpoints to perform CRUD (Create, Read, Update, Delete) operations on employee records.

## Prerequisites

Before running this application, ensure you have the following installed:

- Go (v1.16 or higher)
- MongoDB
- Fiber (`github.com/gofiber/fiber/v2`)
- MongoDB Go driver (`go.mongodb.org/mongo-driver`)

## Installation

1. Clone the repository:

```
git clone https://github.com/your_username/fiber-hrms.git
```

2. Navigate to the project directory:

```
cd fiber-hrms
```

3. Install dependencies:

```
go mod download
```

4. Set up MongoDB:
   - Ensure MongoDB is running locally or provide the appropriate connection URI in the code (`mongoURI` constant).

5. Build and run the application:

```
go run main.go
```

## Usage

Once the application is running, you can interact with it through HTTP requests. Here are the available endpoints:

- **GET /employee**: Fetch all employees.
- **POST /employee**: Create a new employee.
- **PUT /employee/:id**: Update an existing employee.
- **DELETE /employee/:id**: Delete an employee.

Example usage:

```
curl -X GET http://localhost:3000/employee
```

## Configuration

You can configure the MongoDB connection URI and other settings by modifying the constants defined in the code.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

This project is licensed under the [MIT License](LICENSE).
