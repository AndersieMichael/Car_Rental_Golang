package booking

import (
	"database/sql"
	"golangSecond/utilities/db"
	Err "golangSecond/utilities/error"
	"golangSecond/utilities/webhook"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Booking(router *gin.Engine){
	route := router.Group("/booking")

	// GET BOOKING
	//=============================================================	
	route.GET("/get",func(c *gin.Context) {
		// CONNECT DB
		//=============================================================
		db := db.Connect()	
		
		booking_get_result, err := getBooking(db)
		if err != nil {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "booking | request.go | get")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getBooking function",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}
		// ASSEMBLY RESPONSE
		//=============================================================
		defer db.Close()
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
		book_id,_ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()

		book_get_id_result, err := getBookingByID(db,book_id)
		if err != nil && err != sql.ErrNoRows {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "Booking | request.go | get/id")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getBookingByID function",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Booking not exist
		if book_get_id_result.Booking_ID == 0 {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", "Booking ID Not Found", "Booking | request.go | get/id")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error booking not Found",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ASSEMBLY RESPONSE
		//=============================================================
		defer db.Close()
		data := Response{
			Message: "Success",
			Data:    book_get_id_result,
		}

		c.JSON(200, data)
	})

	// ADD BOOKING
	//=============================================================	
	route.POST("/add",func(c *gin.Context) {
		var body BookingForm

		err := c.ShouldBindBodyWith(&body,binding.JSON)
		if err!=nil{
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "book | request.go | add")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		
		book_add_result, err := addBooking(db,body)
		if err != nil {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "booking | request.go | get")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in addBooking function",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ASSEMBLY RESPONSE
		//=============================================================
		defer db.Close()
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

		err := c.ShouldBindBodyWith(&body,binding.JSON)
		if err!=nil{
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "book | request.go | update")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		
		book_get_id_result, err := getBookingByID(db,book_id)
		if err != nil && err != sql.ErrNoRows {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "Booking | request.go | update")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getBookingByID function",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Booking not exist
		if book_get_id_result.Booking_ID == 0 {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", "Booking ID Not Found", "Booking | request.go | update")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error booking not Found",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		err = updateBooking(db,book_id,body)
		if err != nil {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "booking | request.go | update")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in updateBooking function",
			}
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

	// DELETE BOOKING
	//=============================================================		
	route.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		book_id, _ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		
		book_get_id_result, err := getBookingByID(db,book_id)
		if err != nil && err != sql.ErrNoRows {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "Booking | request.go | delete")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getBookingByID function",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Booking not exist
		if book_get_id_result.Booking_ID == 0 {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", "Booking ID Not Found", "Booking | request.go | delete")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error booking not Found",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		err = deleteBooking(db,book_id)
		if err != nil {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "booking | request.go | delete")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in deleteBooking function",
			}
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