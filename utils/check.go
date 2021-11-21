package utils

type Mobile_Response struct {
	Status          string
	StatusDesc      string
	Mebussname      []string
	PgMerchantID    []string
	CrtDate         []string
	IntegrationType []string
	MerVirtualAdd   []string
	LegalStrName    []string
	MasterName      []string
	AdminName       []string
	ListOfMaps      []string
}
type Check_mobile_number struct {
	StrCntMobile string `json:"strCntMobile" validate:"required,len=10,numeric"`
}
type Mobiledata struct {
	Mobile_count int64
}
