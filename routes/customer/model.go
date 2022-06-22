package customer

type Response struct {
	Message        string      `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Error_Key      string      `json:"error_key,omitempty"`
	Error_Message  string      `json:"error_message,omitempty"`
	Secondary_Data interface{} `json:"secondary_data,omitempty"`
}

type GetCustomer struct {
	Customer_ID  int    `json:"customer_id" db:"customer_id"`
	Name         *string `json:"name" db:"name"`
	Nik          *string `json:"nik" db:"nik"`
	Phone_number *string `json:"phone_number" db:"phone_number"`
}

type CustomerForm struct {
	Membership_ID *int    `json:"membership_id" binding:"required"`
	Name          *string `json:"name" binding:"required"`
	Nik           *string `json:"nik" binding:"required"`
	Phone_number  *string `json:"phone_number" binding:"required"`
	Password      *string `json:"password" binding:"required"`
}

type CustomerUpdateForm struct {
	Membership_ID *int    `json:"membership_id" `
	Name          *string `json:"name" `
	Nik           *string `json:"nik" `
	Phone_number  *string `json:"phone_number" `
}
