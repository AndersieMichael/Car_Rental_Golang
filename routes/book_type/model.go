package booktype

type Response struct {
	Message        string      `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Error_Key      string      `json:"error_key,omitempty"`
	Error_Message  string      `json:"error_message,omitempty"`
	Secondary_Data interface{} `json:"secondary_data,omitempty"`
}

type GetType struct{
	Booktype_id 	int `json:"booktype_id" db:"booktype_id"`
	Name 			string `json:"name" db:"name"`
	Description		string `json:"description" db:"description"`
}