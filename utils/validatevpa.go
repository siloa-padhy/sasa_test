package utils

type Vparequest struct {
	VirtualAddress string `json:"virtualAddress" validate:"required"`
	VAReqType      string `json:"vAReqType" validate:"required"`
	PgMerchantId   string `json:"pgMerchantId"`
	PspRefNo       string `json:"pspRefNo"`
	Authorization  string `json:"Authorization"`
}
type Vparesponse struct {
	Status      string
	StatusDesc  string
	RequestInfo string
	PspId       string
	PspRefNo    string
}
