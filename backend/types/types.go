package types

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TimeStamp struct {
	StampType    string `json:"stamp_type"`
	UserID       uint32 `json:"user_id"`
	TimeStampID  uint32 `json:"timestamp_id"`
	Year         int16  `json:"year"`
	Month        uint8  `json:"month"`
	Day          uint8  `json:"day"`
	Hour         uint8  `json:"hour"`
	Minute       uint8  `json:"minute"`
	Second       uint8  `json:"second"`
}
