package driver_incentive

import (
	"database/sql"
	"golangSecond/utilities/db"
	Err "golangSecond/utilities/error"
	"golangSecond/utilities/webhook"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Incentive(router *gin.Engine) {
	parent := "/incentive"
	route := router.Group(parent)

	// GET INCENTIVE
	//=============================================================
	route.GET("/get", func(c *gin.Context) {
		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()

		incentive_get_result, err := getIncentive(db)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getIncentive function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent+" | request.go | get")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}
		// ASSEMBLY RESPONSE
		//=============================================================
		data := Response{
			Message: "Success",
			Data:    incentive_get_result,
		}

		c.JSON(200, data)

	})

	// GET INCENTIVE BY ID
	//=============================================================
	route.GET("/get/:id", func(c *gin.Context) {
		id := c.Param("id")
		incentive_id, _ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()

		incentive_get_id_result, err := getIncentiveByID(db, incentive_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getIncentiveByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent+" | request.go | get/id")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Incentive not exist
		if incentive_get_id_result.Driver_incentive_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error incentive not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, parent+" | request.go | get/id")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ASSEMBLY RESPONSE
		//=============================================================
		data := Response{
			Message: "Success",
			Data:    incentive_get_id_result,
		}

		c.JSON(200, data)
	})

	// ADD INCENTIVE
	//=============================================================
	route.POST("/add", func(c *gin.Context) {
		var body IncentiveForm

		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent+" | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()

		incentive_add_result, err := AddIncentive(db, body)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in addIncentive function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent+" | request.go | get")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ASSEMBLY RESPONSE
		//=============================================================
		data := Response{
			Message: "Success",
			Data:    incentive_add_result,
		}
		c.JSON(200, data)

	})

	// UPDATE INCENTIVE
	//=============================================================
	route.PUT("/update/:id", func(c *gin.Context) {
		id := c.Param("id")

		var body IncentiveForm

		incentive_id, _ := strconv.Atoi(id)

		err := c.ShouldBindBodyWith(&body, binding.JSON)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent+" | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()

		incentive_get_id_result, err := getIncentiveByID(db, incentive_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getIncentiveByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent+" | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Incentive not exist
		if incentive_get_id_result.Driver_incentive_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error incentive not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, parent+" | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		err = UpdateIncentive(db, incentive_id, body)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in updateIncentive function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent+" | request.go | update")
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

	// DELETE INCENTIVE
	//=============================================================
	route.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		incentive_id, _ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		defer db.Close()

		incentive_get_id_result, err := getIncentiveByID(db, incentive_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getIncentiveByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent+" | request.go | delete")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if Incentive not exist
		if incentive_get_id_result.Driver_incentive_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error incentive not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, parent+" | request.go | delete")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		err = deleteIncentive(db, incentive_id)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in deleteDriver function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent+" | request.go | delete")
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