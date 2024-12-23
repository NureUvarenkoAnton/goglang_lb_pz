package api

type PositionMessage struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type RegisterPayload struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	UserType string `json:"user_type" form:"user_type" binding:"-"`
}

type BatchRegisterPayload struct {
	Payloads []RegisterPayload `json:"data"`
}
