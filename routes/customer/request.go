package customer

import (
	"database/sql"
	"golangSecond/utilities/db"
	Err "golangSecond/utilities/error"
	"golangSecond/utilities/webhook"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Customer(router *gin.Engine) {
	route := router.Group("/customer")

	//get customer

	route.GET("/get", func(c *gin.Context) {
		// CONNECT DB
		//=============================================================
		db := db.Connect()

		//get customer

		result, err := getCustomer(db)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getCustomer function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "customer | request.go | get")
			Err.HandleError(err)
			c.JSON(200, data)
			return
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

	//get customer by id

	route.GET("/get/:id", func(c *gin.Context) {
		id := c.Param("id")

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		c_id, _ := strconv.Atoi(id)

		//get customer by id
		result, err := getCustomerByID(db, c_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getCustomer function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "customer | request.go | get/id")
			Err.HandleError(err)
			c.JSON(200, data)
			return
		}

		// check if customer not exist
		if result.Customer_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error Customer not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "customer | request.go | get/id")
			Err.HandleError(err)
			c.JSON(200, data)
			return
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

	//add customer

	route.POST("/add", func(c *gin.Context) {
		var body CustomerForm

		//check validation

		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "customer | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		result, err := addCustomer(db, body)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getCustomer function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "customer | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return
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

	//update customer

	route.PUT("/update/:id", func(c *gin.Context) {
		id := c.Param("id")
		var body CustomerUpdateForm
		c_id, _ := strconv.Atoi(id)

		//validation

		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "customer | request.go | Update")
			Err.HandleError(err)
			c.JSON(200, data)
			return
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()

		//checking customer id

		check_customer_result, err := getCustomerByID(db, c_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getCustomerByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "customer | request.go | Update")
			Err.HandleError(err)
			c.JSON(200, data)
			return
		}

		// if customer not exist

		if check_customer_result.Customer_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error Customer not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "customer | request.go | Update")
			Err.HandleError(err)
			c.JSON(200, data)
			return
		}

		//Update Customer

		err = updateCustomer(db, c_id, body)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in updateCustomer function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "customer | request.go | Update")
			Err.HandleError(err)
			c.JSON(200, data)
			return
		}
		// ASSEMBLY RESPONSE
		//=============================================================
		defer db.Close()
		data := Response{
			Message: "Success",
		}
		c.JSON(200, data)
	})

	//delete customer

	route.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")
		c_id, _ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()

		check_customer_result, err := getCustomerByID(db, c_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getCustomerByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "customer | request.go | delete")
			Err.HandleError(err)
			c.JSON(200, data)
			return
		}

		// check if customer exist

		if check_customer_result.Customer_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error Customer not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "customer | request.go | delete")
			Err.HandleError(err)
			c.JSON(200, data)
			return
		}

		//delete customer by id

		err = deleteCustomer(db, c_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in deleteCustomer function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "customer | request.go | delete")
			Err.HandleError(err)
			c.JSON(200, data)
			return
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
