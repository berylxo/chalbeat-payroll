package main

import (
	"log"
	"net/http"

	"payroll-service/routes"
)

func main() {
	r := routes.SetupRouter()

	log.Println("Payroll service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
