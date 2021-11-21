package utils

type Initreq struct {
	Amount              string `json:"amount" validate:"required,numeric"`
	Virtualaddress      string `json:"virtualAddress"`
	Merchanttype        string `json:"merchantType" validate:"required"`
	Issl                string `json:"issl" validate:"required,min=0,max=1,numeric"`
	Message             string `json:"message"`
	Authorization       string `json:"Authorization" validate:"required"`
	Serviceprovidertype string `json:"Serviceprovidertype"`
	PgMerchantId        string
}
type Amount struct {
	Amountindouble   float64
	OrderNo_in_float int64
}
