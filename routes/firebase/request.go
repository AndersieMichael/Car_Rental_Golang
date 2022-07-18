package firebase

import (
	"encoding/base64"
	"golangSecond/utilities/db"
	Err "golangSecond/utilities/error"
	"golangSecond/utilities/webhook"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func MyFirebase(router *gin.Engine){
	parent := "/firebase"
	routes := router.Group(parent)

	routes.POST("/add", func(c *gin.Context) {
		var upload Upload_Image
		var mime string

		err := c.ShouldBindBodyWith(&upload, binding.JSON)
		if err !=nil{
			Err.HandleError(err)
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
		
		db := db.Connect()
		defer db.Close()
		tx,_ := db.Beginx()
		defer tx.Rollback()

		if(upload.Mime_Type)=="jpeg"{
			mime = "image/jpeg"
		}else if upload.Mime_Type == "jpg"{
			mime = "image/jpeg"
		}else if upload.Mime_Type =="png"{
			mime = "image/png"
		}
		decoded, _ := base64.StdEncoding.DecodeString(upload.Image)
		url,err := UploadExcel([]byte(decoded),upload.Filename,mime)
		if err !=nil{
			Err.HandleError(err)
			data := Response{
				Message:       "Failed",
				Error_Key:     "error_internal_server",
				Error_Message: err.Error(),
			}
			webhook.PostToWebHook(c.Request.Method, c.Request.Host+c.Request.URL.Path, data.Error_Key, err.Error(), parent +" | request.go | add")
			Err.HandleError(err)
			c.JSON(200, data)
			return //END
		}
		
		// ASSEMBLY RESPONSE
		//=============================================================
		data := Response{
			Message: "Success",
			Data: url,
		}
		tx.Commit()
		c.JSON(200, data)

	})
}