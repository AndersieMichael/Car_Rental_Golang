package booking

type Response struct {
	Message        string      `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Error_Key      string      `json:"error_key,omitempty"`
	Error_Message  string      `json:"error_message,omitempty"`
	Secondary_Data interface{} `json:"secondary_data,omitempty"`
}

type GetBooking struct {
	Booking_ID        int   `json:"booking_id" db:"booking_id"`
	Customer_ID       *int  `json:"customer_id" db:"customer_id"`
	Cars_ID           *int  `json:"cars_id" db:"cars_id"`
	Start_time        int64 `json:"start_time" db:"start_time"`
	End_time          int64 `json:"end_time" db:"end_time"`
	Total_cost        int   `json:"total_cost" db:"total_cost"`
	Finished          bool  `json:"finished" db:"finished"`
	Discount          int   `json:"discount" db:"discount"`
	Booktype_ID       *int  `json:"booktype_id" db:"booktype_id"`
	Driver_ID         *int  `json:"driver_id" db:"driver_id"`
	Total_driver_cost int   `json:"total_driver_cost" db:"total_driver_cost"`
}

type BookingForm struct {
	Customer_ID       *int  `json:"customer_id" db:"customer_id"`
	Cars_ID           *int  `json:"cars_id" db:"cars_id"`
	Start_time        int64 `json:"start_time" db:"start_time"`
	End_time          int64 `json:"end_time" db:"end_time"`
	Total_cost        int   `json:"total_cost" db:"total_cost"`
	Finished          bool  `json:"finished" db:"finished"`
	Discount          int   `json:"discount" db:"discount"`
	Booktype_ID       *int  `json:"booktype_id" db:"booktype_id"`
	Driver_ID         *int  `json:"driver_id" db:"driver_id"`
	Total_driver_cost int   `json:"total_driver_cost" db:"total_driver_cost"`
}
