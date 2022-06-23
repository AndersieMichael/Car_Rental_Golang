package booking

import (
	"database/sql"
	"golangSecond/routes/car"
	"golangSecond/routes/customer"
	inc "golangSecond/routes/driver_icentive"
	"golangSecond/routes/driver"
	"golangSecond/routes/membership"
	"golangSecond/utilities/db"
	Err "golangSecond/utilities/error"
	"golangSecond/utilities/webhook"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Booking(router *gin.Engine) {
	route := router.Group("/booking")

	// GET BOOKING
	//=============================================================
	route.GET("/get", func(c *gin.Context) {
		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()

		booking_get_result, err := getBooking(db)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getBooking function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | get")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}
		// ASSEMBLY RESPONSE
		//=============================================================
		data := Response{
			Message: "Success",
			Data:    booking_get_result,
		}

		c.JSON(200, data)

	})

	// GET BOOKING BY ID
	//=============================================================
	route.GET("/get/:id", func(c *gin.Context) {
		id := c.Param("id")
		book_id, _ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()

		book_get_id_result, err := getBookingByID(db, book_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getBookingByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "Booking | request.go | get/id")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Booking not exist
		if book_get_id_result.Booking_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error booking not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "Booking | request.go | get/id")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ASSEMBLY RESPONSE
		//=============================================================
		data := Response{
			Message: "Success",
			Data:    book_get_id_result,
		}

		c.JSON(200, data)
	})

	// ADD BOOKING
	//=============================================================
	route.POST("/add", func(c *gin.Context) {
		var body BookingForm

		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "book | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()

		book_add_result, err := addBooking(db, body)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in addBooking function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ASSEMBLY RESPONSE
		//=============================================================
		data := Response{
			Message: "Success",
			Data:    book_add_result,
		}
		c.JSON(200, data)

	})

	// ADD BOOKING FOR V2
	//=============================================================
	route.POST("/v2/add", func(c *gin.Context) {
		var body BookingFormV2
		var body2 inc.IncentiveForm
		var discount int = 0
		var incentive int = 0
		var total int = 0
		var driver_payment int = 0

		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "book | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()

		// GET CUSTOMER DATA
		//=============================================================
		customer_data_result, err := customer.GetCustomerBooking(db, *body.Customer_ID)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in GetCustomerBooking function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Customer_ID not exist
		if customer_data_result.Customer_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error Car not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// GET CAR DATA
		//=============================================================
		car_data_result, err := car.GetCarByID(db, *body.Cars_ID)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getCarByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if CAR not exist
		if car_data_result.Car_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error Car not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		//change from unix
		startT := time.Unix(body.Start_time, 0).Format("2006-01-02")
		endT := time.Unix(body.End_time, 0).Format("2006-01-02")

		//change to time from string
		firstT, _ := time.Parse("2006-01-02", (startT))
		secondT, _ := time.Parse("2006-01-02", (endT))

		//calculate days range
		//=============================================================
		rangeT := (secondT.Sub(firstT).Hours() / 24) + 1

		//Calculate total_cost
		//=============================================================
		total = int(rangeT) * int(car_data_result.Rent_price_daily)

		// GET MEMBERSHIP DATA
		//=============================================================
		membership_data_result, err := membership.GetMembershipByID(db, customer_data_result.Membership_ID)

		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in GetMembershipByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if membership not exist
		//=============================================================
		if membership_data_result.Membership_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error membership not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		//if membership discount not 0 (membership id 7)
		//=============================================================
		if membership_data_result.Daily_discount != 0 {

			discount = total * int(membership_data_result.Daily_discount/100)

		}

		//checking booktypeID [HARDCODE]
		//=============================================================
		if body.Booktype_ID == 2 {
			if body.Driver_ID == nil {
				data := Response{
					Message:       "Failed",
					Error_Key:     "error_insert_driverID",
					Error_Message: "Driver cannot be null for booktype id 2",
				}
				webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | add")
				Err.HandleError(err)
				c.JSON(200, data)
				return //END
			}

			// GET DRIVER DATA
			//=============================================================
			driver_data_result, err := driver.GetDriverByID(db, *body.Driver_ID)

			if err != nil && err != sql.ErrNoRows {
				data := Response{
					Message:       "Failed",
					Error_Key:     "error_internal_server",
					Error_Message: "error in GetDriverByID function",
				}
				webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | add")
				Err.HandleError(err)
				c.JSON(200, data)
				return //END
			}

			// check if driver not exist
			//=============================================================
			if driver_data_result.Driver_ID == 0 {
				data := Response{
					Message:       "Failed",
					Error_Key:     "error_id_not_found",
					Error_Message: "error driver not Found",
				}
				webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | add")
				Err.HandleError(err)
				c.JSON(200, data)
				return //END
			}

			//CALCULATE TOTAL_DRIVER_COST
			//=============================================================
			driver_payment = int(driver_data_result.Daily_cost) * int(rangeT)

		} else if body.Booktype_ID != 1 { //error selain 1 dan 2
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_insert_booktype_id",
				Error_Message: "book type must be 1 or 2",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		

		// ADD BOOKING DATA
		//=============================================================
		book_add_result, err := addBookingV2(db, body,
			startT,
			endT,
			total,
			discount,
			driver_payment)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in addBookingV2 function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}
		
		if body.Booktype_ID == 2 {

			//CALCULATE INCENTIVE
			//=============================================================
			incentive = total * 5 / 100
			body2.Booking_ID = book_add_result
			body2.Incentive = incentive

			// ADD INCENTIVE DATA
			//=============================================================
			_, err := inc.AddIncentive(db,body2)
			if err != nil {
				data := Response{
					Message:       "Failed",
					Error_Key:     "error_internal_server",
					Error_Message: "error in AddIncentive function",
				}
				webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | add")
				Err.HandleError(err)
				c.JSON(200, data)
				return //END
			}
		}

		// ASSEMBLY RESPONSE
		//=============================================================
		data := Response{
			Message: "Success",
			Data:    book_add_result,
		}
		c.JSON(200, data)
	})

	// UPDATE BOOKING
	//=============================================================
	route.PUT("/update/:id", func(c *gin.Context) {
		id := c.Param("id")

		var body BookingForm

		book_id, _ := strconv.Atoi(id)

		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "book | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()

		book_get_id_result, err := getBookingByID(db, book_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getBookingByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "Booking | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Booking not exist
		if book_get_id_result.Booking_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error booking not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "Booking | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		err = updateBooking(db, book_id, body)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in updateBooking function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ASSEMBLY RESPONSE
		//=============================================================
		data := Response{
			Message: "Success",
		}
		c.JSON(200, data)

	})

	// DELETE BOOKING
	//=============================================================
	route.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		book_id, _ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()

		book_get_id_result, err := getBookingByID(db, book_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getBookingByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "Booking | request.go | delete")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Booking not exist
		if book_get_id_result.Booking_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error booking not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "Booking | request.go | delete")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		err = deleteBooking(db, book_id)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in deleteBooking function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | delete")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ASSEMBLY RESPONSE
		//=============================================================
		defer db.Close()
		data := Response{
			Message: "Success",
		}
		c.JSON(200, data)
	})
}
