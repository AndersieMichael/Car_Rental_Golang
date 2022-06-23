package main

import (
	"golangSecond/routes/booking"
	"golangSecond/routes/car"
	"golangSecond/routes/customer"
	"golangSecond/routes/driver"
	"golangSecond/routes/driver_icentive"
	"golangSecond/routes/membership"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	customer.Customer(router)
	car.Car(router)
	booking.Booking(router)
	membership.Membership(router)
	driver.Driver(router)
	driver_incentive.Incentive(router)

	router.Run(":8080")
}