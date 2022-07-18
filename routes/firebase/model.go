package firebase

type Response struct {
	Message        string      `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Error_Key      string      `json:"error_key,omitempty"`
	Error_Message  string      `json:"error_message,omitempty"`
	Secondary_Data interface{} `json:"secondary_data,omitempty"`
}

type Upload_Image struct {
	Filename             string `binding:"required"`
	Image                string `binding:"required"`
	Mime_Type            string `binding:"required"`
}