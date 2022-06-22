package car

import (
	"database/sql"
	"golangSecond/utilities/db"
	Err "golangSecond/utilities/error"
	"golangSecond/utilities/webhook"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Car(router *gin.Engine) {
	route := router.Group("/car")


	// GET CAR
	//=============================================================	
	route.GET("/get", func(c *gin.Context) {
		// CONNECT DB
		//=============================================================
		db := db.Connect()

		result, err := getCar(db)
		if err != nil {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "car | request.go | get")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getCar function",
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
			Data:    result,
		}

		c.JSON(200, data)
	})

	// GET CAR BY ID
	//=============================================================		
	route.GET("/get/:id", func(c *gin.Context) {
		id := c.Param("id")
		car_id, _ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		result, err := getCarByID(db, car_id)
		if err != nil && err != sql.ErrNoRows {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "car | request.go | get/id")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getCarByID function",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if CAR not exist
		if result.Car_ID == 0 {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", "Car ID Not Found", "car | request.go | get/id")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error Car not Found",
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
			Data:    result,
		}

		c.JSON(200, data)
	})

	// ADD CAR
	//=============================================================		
	route.POST("/add", func(c *gin.Context) {
		var body CarForm

		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "car | request.go | add")
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

		car_add_result, err := addCar(db, body)
		if err != nil {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "car | request.go | add")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in addCar function",
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
			Data:    car_add_result,
		}
		c.JSON(200, data)
	})

	// UPDATE CAR
	//=============================================================		
	route.PUT("/update/:id", func(c *gin.Context) {
		id := c.Param("id")

		var body CarForm

		car_id, _ := strconv.Atoi(id)

		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "car | request.go | update")
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

		car_get_id_result, err := getCarByID(db, car_id)
		if err != nil && err != sql.ErrNoRows {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "car | request.go | update")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getCarByID function",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if CAR not exist
		if car_get_id_result.Car_ID == 0 {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", "Car ID Not Found", "car | request.go | update")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error Car not Found",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		err = updateCar(db, car_id, body)
		if err != nil {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "car | request.go | update")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in updateCar function",
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

	// DELETE CAR
	//=============================================================		
	route.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		car_id, _ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()

		car_get_id_result, err := getCarByID(db, car_id)
		if err != nil && err != sql.ErrNoRows {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "car | request.go | delete")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getCarByID function",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if CAR not exist
		if car_get_id_result.Car_ID == 0 {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", "Car ID Not Found", "car | request.go | delete")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error Car not Found",
			}
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		err = deleteCar(db, car_id)

		if err != nil {
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, "Failed", err.Error(), "car | request.go | delete")
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in deleteCar function",
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
