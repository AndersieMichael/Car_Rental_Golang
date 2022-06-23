package driver

type Response struct {
	Message        string      `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Error_Key      string      `json:"error_key,omitempty"`
	Error_Message  string      `json:"error_message,omitempty"`
	Secondary_Data interface{} `json:"secondary_data,omitempty"`
}

type GetDriver struct {
	Driver_ID		int   `json:"driver_id" db:"driver_id"`
	Name       		string  `json:"name" db:"name"`
	Nik           	string  `json:"nik" db:"nik"`
	Phone_number	string `json:"phone_number" db:"phone_number"`
	Daily_cost      int `json:"daily_cost" db:"daily_cost"`
}

type DriverForm struct {
	Name       		string  `json:"name" binding:"required"`
	Nik           	string  `json:"nik" binding:"required"`
	Phone_number	string `json:"phone_number" binding:"required"`
	Daily_cost      int `json:"daily_cost" binding:"required"`
}