package booking

import (
	"database/sql"
	"golangSecond/routes/car"
	"golangSecond/routes/customer"
	"golangSecond/routes/driver"
	inc "golangSecond/routes/driver_icentive"
	"golangSecond/routes/membership"
	"golangSecond/utilities/db"
	Err "golangSecond/utilities/error"
	"golangSecond/utilities/webhook"
	"math"
	"strconv"

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

		// if car stock is empty
		if car_data_result.Stock == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_stock_is_empty",
				Error_Message: "error Car stock is empty",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		//check if end time lower than start time
		if body.End_time < body.Start_time {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_time_range",
				Error_Message: "end time cannot lower than start time",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		//calculate days range
		//=============================================================
		rangeDate := body.End_time - body.Start_time
		rangeT := int(math.Ceil(float64(rangeDate)/(1000*60*60*24))) + 1

		//Calculate total_cost
		//=============================================================
		total = rangeT * car_data_result.Rent_price_daily

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
		if membership_data_result.Daily_discount != 7 {

			temp := float64(total) * membership_data_result.Daily_discount
			discount = int(temp)
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

			// GET Booking Driver data checking date
			//=============================================================
			check_driver_date_result, err := getBookingByDriverID(db, *body.Driver_ID, body.Start_time, body.End_time)
			if err != nil && err != sql.ErrNoRows {
				data := Response{
					Message:       "Failed",
					Error_Key:     "error_internal_server",
					Error_Message: "error in check_driver_date_result function",
				}
				webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | add")
				Err.HandleError(err)
				c.JSON(200, data)
				return //END
			}

			// check if found same date
			//=============================================================
			if check_driver_date_result.Booking_ID != 0 {
				data := Response{
					Message:       "Failed",
					Error_Key:     "error_date_book",
					Error_Message: "error driver already have booking for that date",
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

		// EXECUTE ALL PROCES
		//=============================================================

		// UPDATE CAR STOCK
		//=============================================================
		stock := car_data_result.Stock - 1

		err = car.UpdateCarQuantity(db, *body.Cars_ID, stock)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in UpdateCarQuantity function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ADD BOOKING DATA
		//=============================================================
		book_add_result, err := addBookingV2(db, body,
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
			_, err := inc.AddIncentive(db, body2)
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

	// UPDATE BOOKING FOR V2
	//=============================================================
	route.PUT("/v2/update/:id", func(c *gin.Context) {
		var body2 inc.IncentiveForm
		var discount int = 0
		var incentive int = 0
		var total int = 0
		var driver_payment int = 0
		id := c.Param("id")

		var body BookingFormV2

		book_id, _ := strconv.Atoi(id)

		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "Booking | request.go | update")
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

		// GET CUSTOMER DATA
		//=============================================================
		customer_data_result, err := customer.GetCustomerBooking(db, *body.Customer_ID)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in GetCustomerBooking function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | update")
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
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | update")
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
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | update")
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
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// if car stock is empty
		if car_data_result.Stock == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_stock_is_empty",
				Error_Message: "error Car stock is empty",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		//check if end time lower than start time
		if body.End_time < body.Start_time {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_time_range",
				Error_Message: "end time cannot lower than start time",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		//calculate days range
		//=============================================================
		rangeDate := body.End_time - body.Start_time
		rangeT := int(math.Ceil(float64(rangeDate)/(1000*60*60*24))) + 1

		//Calculate total_cost
		//=============================================================
		total = rangeT * car_data_result.Rent_price_daily

		// GET MEMBERSHIP DATA
		//=============================================================
		membership_data_result, err := membership.GetMembershipByID(db, customer_data_result.Membership_ID)

		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in GetMembershipByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | update")
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
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		//if membership discount not 0 (membership id 7)
		//=============================================================
		if membership_data_result.Daily_discount != 7 {

			temp := float64(total) * membership_data_result.Daily_discount
			discount = int(temp)
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
				webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | update")
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
				webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | update")
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
				webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | update")
				Err.HandleError(err)
				c.JSON(200, data)
				return //END
			}

			// GET Booking Driver data checking date
			//=============================================================
			check_driver_date_result, err := getBookingByDriverID(db, *body.Driver_ID, body.Start_time, body.End_time)
			if err != nil && err != sql.ErrNoRows {
				data := Response{
					Message:       "Failed",
					Error_Key:     "error_internal_server",
					Error_Message: "error in check_driver_date_result function",
				}
				webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | update")
				Err.HandleError(err)
				c.JSON(200, data)
				return //END
			}

			// check if found same date
			//=============================================================
			if check_driver_date_result.Booking_ID != 0 {
				data := Response{
					Message:       "Failed",
					Error_Key:     "error_date_book",
					Error_Message: "error driver already have booking for that date",
				}
				webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | update")
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
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "booking | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}else{
			body.Driver_ID = nil
		}

		// EXECUTE ALL PROCES
		//=============================================================

		// UPDATE CAR STOCK
		//=============================================================
		stock := car_data_result.Stock - 1

		err = car.UpdateCarQuantity(db, *body.Cars_ID, stock)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in UpdateCarQuantity function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ADD BOOKING DATA
		//=============================================================
		err = updateBookingV2(db, book_id, body,
			total,
			discount,
			driver_payment)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in updateBookingV2 function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		if body.Booktype_ID == 2 {

			//CALCULATE INCENTIVE
			//=============================================================
			incentive = total * 5 / 100
			body2.Booking_ID = book_id
			body2.Incentive = incentive

			incentive_get_id_result, err := inc.GetIncentiveByBOOKID(db, book_id)
			if err != nil && err != sql.ErrNoRows {
				data := Response{
					Message:       "Failed",
					Error_Key:     "error_internal_server",
					Error_Message: "error in GetIncentiveByBOOKID function",
				}
				webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | get/id")
				Err.HandleError(err)
				c.JSON(200, data)
				return //END
			}

			// check if Incentive not exist
			if incentive_get_id_result.Driver_incentive_ID == 0 {

				// ADD INCENTIVE DATA
				//=============================================================
				_, err := inc.AddIncentive(db, body2)
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
			}else{
				
				// UPDATE INCENTIVE DATA
				//=============================================================
				err = inc.UpdateIncentivebyBOOK(db, book_id, body2)
				if err != nil {
					data := Response{
						Message:       "Failed",
						Error_Key:     "error_internal_server",
						Error_Message: "error in UpdateIncentivebyBOOK function",
					}
					webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | update")
					Err.HandleError(err)
					c.JSON(200, data)
					return //END
				}
			}
		}else{
			err := inc.DeleteIncentivebyBOOK(db,book_id)
			if err != nil {
				data := Response{
					Message:       "Failed",
					Error_Key:     "error_internal_server",
					Error_Message: "error in DeleteIncentive function",
				}
				webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "booking | request.go | update")
				Err.HandleError(err)
				c.JSON(200, data)
				return //END
			}
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
