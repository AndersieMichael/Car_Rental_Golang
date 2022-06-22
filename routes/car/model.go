package car

type Response struct {
	Message        string      `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Error_Key      string      `json:"error_key,omitempty"`
	Error_Message  string      `json:"error_message,omitempty"`
	Secondary_Data interface{} `json:"secondary_data,omitempty"`
}

type GetCar struct{
	Car_ID 				int `json:"car_id" db:"cars_id"`
	Name 				string `json:"car_name" db:"name"`
	Rent_price_daily	int `json:"rent_price_daily" db:"rent_price_daily"`
	Stock 				int `json:"stock" db:"stock"`
}

type CarForm struct{
	Name 				string `json:"car_name" binding:"required"`
	Rent_price_daily	int `json:"rent_price_daily" binding:"required"`
	Stock 				int `json:"stock" binding:"required"`
}

