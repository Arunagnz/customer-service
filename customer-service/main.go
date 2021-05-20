package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arunagnz/customer-service/models"
	"github.com/arunagnz/customer-service/postgres"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

var store *postgres.Store
var err error

func init() {
	// load environment variables
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file : ", err)
	// }
	connectDB()
}

func connectDB() {
	// host := os.Getenv("DB_HOST")
	// port := os.Getenv("DB_PORT")
	// user := os.Getenv("DB_USER")
	// password := os.Getenv("DB_PASSWORD")
	// dbname := os.Getenv("DB_NAME")
	host := "host.docker.internal"
	port := "5432"
	user := "postgres"
	password := "postgres"
	dbname := "crud_demo"

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, password, port)

	// Openning connection to database
	store, err = postgres.NewStore(dbURI)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to database successfully")
	}
}

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return show(req)
	case "POST":
		return create(req)
	case "PUT":
		return update(req)
	case "DELETE":
		return delete(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cc, err := store.Customers()
	if err != nil {
		return serverError(err)
	}
	if cc == nil {
		return clientError(http.StatusNotFound)
	}

	js, err := json.Marshal(cc)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	c := new(models.Customer)
	err := json.Unmarshal([]byte(req.Body), c)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	if c.Name == "" || c.Email == "" || c.Password == "" || c.PhoneNumber == "" || c.Address == "" {
		return clientError(http.StatusBadRequest)
	}

	c.ID = uuid.New()

	err = store.CreateCustomer(c)
	if err != nil {
		return serverError(err)
	}

	js, err := json.Marshal(c)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(js),
	}, nil
}

func update(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	c := new(models.Customer)
	err := json.Unmarshal([]byte(req.Body), c)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	err = store.UpdateCustomer(c)
	if err != nil {
		return serverError(err)
	}

	values := map[string]string{"message": "Customer updated successfully", "id": c.ID.String()}

	js, err := json.Marshal(values)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func delete(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	c := new(models.Customer)
	err := json.Unmarshal([]byte(req.Body), c)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	err = store.DeleteCustomer(c.ID)
	if err != nil {
		return serverError(err)
	}

	values := map[string]string{"message": "Customer deleted successfully", "id": c.ID.String()}

	js, err := json.Marshal(values)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(router)
}
