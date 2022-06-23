package membership

type Response struct {
	Message        string      `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Error_Key      string      `json:"error_key,omitempty"`
	Error_Message  string      `json:"error_message,omitempty"`
	Secondary_Data interface{} `json:"secondary_data,omitempty"`
}

type GetMembership struct {
	Membership_ID  int    `json:"membership_id" db:"membership_id"`
	Name           string `json:"name" db:"name"`
	Daily_discount float64    `json:"daily_discount" db:"daily_discount"`
}

type MembershipFormat struct {
	Name           string `json:"name" db:"name"`
	Daily_discount float64    `json:"daily_discount" db:"daily_discount"`
}
