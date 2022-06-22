package main

import (
	"golangSecond/routes/booking"
	"golangSecond/routes/car"
	"golangSecond/routes/customer"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	customer.Customer(router)
	car.Car(router)
	booking.Booking(router)
	
	router.Run(":8080")
}