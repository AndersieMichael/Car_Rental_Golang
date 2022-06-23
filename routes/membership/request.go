package membership

import (
	"database/sql"
	"golangSecond/utilities/db"
	Err "golangSecond/utilities/error"
	"golangSecond/utilities/webhook"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Membership(router *gin.Engine){
	route := router.Group("/membership")

	// GET MEMBERSHIP
	//=============================================================	
	route.GET("/get",func (c *gin.Context)  {
		// CONNECT DB
		//=============================================================
		db := db.Connect()

		membership_get_result, err := getMembership(db)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getMembership function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "membership | request.go | get")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}
		// ASSEMBLY RESPONSE
		//=============================================================
		defer db.Close()
		data := Response{
			Message: "Success",
			Data:    membership_get_result,
		}

		c.JSON(200, data)

	})


	// GET MEMBERSHIP BY ID
	//=============================================================	
	route.GET("/get/:id",func (c *gin.Context)  {

		id := c.Param("id")
		
		member_id,_ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()

		membership_get_id_result, err := GetMembershipByID(db,member_id)

		if err != nil && err != sql.ErrNoRows  {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getMembershipByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "membership | request.go | get/id")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if membership not exist
		if membership_get_id_result.Membership_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error membership not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "membership | request.go | get/id")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}				

		// ASSEMBLY RESPONSE
		//=============================================================
		defer db.Close()
		data := Response{
			Message: "Success",
			Data:    membership_get_id_result,
		}

		c.JSON(200, data)

	})	
	
	// ADD MEMBERSHIP
	//=============================================================	
	route.POST("/add",func(c *gin.Context) {
		var body MembershipFormat

		err := c.ShouldBindBodyWith(&body,binding.JSON)
		if err!=nil{
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "membership | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		
		membership_add_result, err := addMembership(db,body)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in addMembership function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "membership | request.go | get")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// ASSEMBLY RESPONSE
		//=============================================================
		defer db.Close()
		data := Response{
			Message: "Success",
			Data:    membership_add_result,
		}
		c.JSON(200, data)		
		
	})

	// UPDATE MEMBERSHIP
	//=============================================================		
	route.PUT("/update/:id", func(c *gin.Context) {
		id := c.Param("id")
		
		var body MembershipFormat

		member_id, _ := strconv.Atoi(id)

		err := c.ShouldBindBodyWith(&body,binding.JSON)
		if err!=nil{
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_param",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "membership | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		
		membership_get_id_result, err := GetMembershipByID(db,member_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getMembershipByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "membership | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if membership not exist
		if membership_get_id_result.Membership_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error membership not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "membership | request.go | update")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		err = updateMembership(db,member_id,body)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in updateMembership function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "membership | request.go | update")
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


	// DELETE MEMBERSHIP
	//=============================================================		
	route.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		member_id, _ := strconv.Atoi(id)

		// CONNECT DB
		//=============================================================
		db := db.Connect()
		
		membership_get_id_result, err := GetMembershipByID(db,member_id)
		if err != nil && err != sql.ErrNoRows {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in getMembershipByID function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "membership | request.go | delete")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		// check if membership not exist
		if membership_get_id_result.Membership_ID == 0 {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_id_not_found",
				Error_Message: "error membership not Found",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, data.Error_Message, "membership | request.go | delete")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}

		err = deleteMembership(db,member_id)
		if err != nil {
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: "error in deleteMembership function",
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), "membership | request.go | delete")
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