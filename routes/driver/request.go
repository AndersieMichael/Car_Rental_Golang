package driver

import (
	"database/sql"
	"golangSecond/utilities/db"
	Err "golangSecond/utilities/error"
	"golangSecond/utilities/webhook"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Driver(router *gin.Engine){
	parent := "/driver"
	route := router.Group(parent)

	// GET Driver
	//=============================================================	
	route.GET("/get",func(c *gin.Context) {
		// CONNECT DB
		//=============================================================
		db := db.Connect()	
		defer db.Close()

		driver_get_result, err := getDriver(db)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getDriver function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent + " | request.go | get")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}
		// ASSEMBLY RESPONSE
		//=============================================================
		data := Response{
			Message: "Success",
			Data:    driver_get_result,
		}

		c.JSON(200, data)

	})

	// GET DRIVER BY ID
	//=============================================================	
	route.GET("/get/:id", func(c *gin.Context) {
		id := c.Param("id")
		driver_id,_ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()

		driver_get_id_result, err := GetDriverByID(db,driver_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getDriverByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent +" | request.go | get/id")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Driver not exist
		if driver_get_id_result.Driver_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error Driver not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, parent + " | request.go | get/id")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ASSEMBLY RESPONSE
		//=============================================================
		data := Response{
			Message: "Success",
			Data:    driver_get_id_result,
		}

		c.JSON(200, data)
	})

	// ADD DRIVER
	//=============================================================	
	route.POST("/add",func(c *gin.Context) {
		var body DriverForm

		err := c.ShouldBindBodyWith(&body,binding.JSON)
		if err!=nil{
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent + " | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()
		
		driver_add_result, err := addDriver(db,body)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in addDriver function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent + " | request.go | get")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ASSEMBLY RESPONSE
		//=============================================================
		data := Response{
			Message: "Success",
			Data:    driver_add_result,
		}
		c.JSON(200, data)		
		
	})
	
	// UPDATE DRIVER
	//=============================================================		
	route.PUT("/update/:id", func(c *gin.Context) {
		id := c.Param("id")
		
		var body DriverForm

		driver_id, _ := strconv.Atoi(id)

		err := c.ShouldBindBodyWith(&body,binding.JSON)
		if err!=nil{
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent + " | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()
		
		driver_get_id_result, err := GetDriverByID(db,driver_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getDriverByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent + " | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Driver not exist
		if driver_get_id_result.Driver_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error Driver not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message,  parent + " | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		err = updateDriver(db,driver_id,body)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in updateDriver function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(),  parent + " | request.go | update")
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

	// DELETE DRIVER
	//=============================================================		
	route.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		driver_id, _ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()
		
		driver_get_id_result, err := GetDriverByID(db,driver_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getDriverByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent + " | request.go | delete")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Driver not exist
		if driver_get_id_result.Driver_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error driver not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, parent + " | request.go | delete")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		err = deleteDriver(db,driver_id)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in deleteDriver function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent + " | request.go | delete")
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


}