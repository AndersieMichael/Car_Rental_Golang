package driver_incentive

type Response struct {
	Message        string      `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Error_Key      string      `json:"error_key,omitempty"`
	Error_Message  string      `json:"error_message,omitempty"`
	Secondary_Data interface{} `json:"secondary_data,omitempty"`
}

type GetIncentive struct{
	Driver_incentive_ID int `json:"driver_incentive_id" db:"driver_incentive_id"`
	Booking_ID 			int `json:"booking_id" db:"booking_id"`
	Incentive			int `json:"incentive" db:"incentive"`
}

type IncentiveForm struct{
	Booking_ID 			int `json:"booking_id"`
	Incentive			int `json:"incentive" binding:"required"`
}