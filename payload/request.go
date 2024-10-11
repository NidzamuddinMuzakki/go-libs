package payload

type (
	Device struct {
		DeviceType string  `example:"Android" json:"device_type" validate:"required"`
		DeviceID   string  `example:"EA7583CD-A667-48BC-B806-42ECB2B48606" json:"device_id" validate:"required"`
		Longitude  float64 `example:"1234" json:"longitude" validate:"required"`
		Latitude   float64 `example:"-1223" json:"latitude" validate:"required"`
	}
	ReqGeneral struct {
		Type    string      `json:"type"`
		Service string      `json:"service"`
		Data    interface{} `json:"data"`
	}
	ReqListGeneral struct {
		Device        Device         `json:"device"`
		Criteria      []Criteria     `json:"criteria"`
		SortCriteria  []SortCriteria `json:"sort_criteria"`
		PageNum       int            `json:"page_num"`
		RecordPerPage int            `json:"record_per_page"`
	}
	Criteria struct {
		Value        string `json:"value" binding:"required"`
		AnotherValue string `json:"another_value" binding:"required_if=Operator BETWEEN"`
		Field        string `json:"field" binding:"required"`
		Operator     string `json:"operator" binding:"required"`
	}
	SortCriteria struct {
		SortOrder string `json:"sort_order" binding:"required,oneof=asc desc ASC DESC"`
		Field     string `json:"field" binding:"required"`
	}
)
