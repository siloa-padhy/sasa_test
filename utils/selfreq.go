package utils

type Selfonboard struct {
	PspId           string `json:"pspId"`
	PgMerchantId    string `json:"pgMerchantId"`
	Mebussname      string `json:"mebussname" validate:"required"`
	LegalStrName    string `json:"legalStrName" validate:"required"`
	MerVirtualAdd   string `json:"merVirtualAdd" validate:"required"`
	Awlmcc          string `json:"awlmcc"`
	StrCntMobile    string `json:"strCntMobile" validate:"required,len=10"`
	RequestUrl1     string `json:"requestUrl1"`
	RequestUrl2     string `json:"requestUrl2"`
	MerchantType    string `json:"merchantType"`
	IntegrationType string `json:"integrationType"`
	SettleType      string `json:"settleType"`
	PanNo           string `json:"panNo" validate:"required"`
	UserName        string `json:"userName"`
	ExtTID          string `json:"extTID"`
	ExtMID          string `json:"extMID"`
	MeEmailID       string `json:"meEmailID"`
	FirstName       string `json:"firstName" validate:"required"`
	LastName        string `json:"lastname" validate:"required"`
	AccNo           string `json:"accNo" validate:"required"`
	Ifsc            string `json:"ifsc" validate:"required"`
	PspRefNo        string `json:"pspRefNo"`
	VirtualAddress  string `json:"virtualAddress"`
	VAReqType       string `json:"vAReqType"`
	BankName        string `json:"bankName" validate:"required"`
	Authorization   string `json:"Authorization" validate:"required"`
}
type Response struct {
	Status          int    `json:"status"`
	StatusDesc      string `json:"statusDesc"`
	Mebussname      string
	PgMerchantID    string
	CrtDate         string
	IntegrationType string
	MerVirtualAdd   string
	LegalStrName    string
	MasterName      []string
	AdminName       []string
	ListOfMaps      []string
}
type Mobile struct {
	Mobileno int64
}
type Response1 struct {
	Qrcode     string
	Status     string
	Message    string
	AdminName  string
	MasterName string
	Vpaid      string
	Storename  string
}
type Response2 struct {
	Status     string
	StatusDesc string
}

// type Detail struct {
// 	detailname string
// }
type Auth struct {
	Header string
}

type SubMerchantOnBoardApiRequest struct {
	PspId           string `json:"pspId"`
	PgMerchantId    string `json:"pgMerchantId"`
	Mebussname      string `json:"mebussname"`
	LegalStrName    string `json:"legalStrName"`
	MerVirtualAdd   string `json:"merVirtualAdd"`
	Awlmcc          string `json:"awlmcc"`
	StrCntMobile    string `json:"strCntMobile"`
	RequestUrl1     string `json:"requestUrl1"`
	RequestUrl2     string `json:"requestUrl2"`
	MerchantType    string `json:"merchantType"`
	IntegrationType string `json:"integrationType"`
	SettleType      string `json:"settleType"`
	PanNo           string `json:"panNo"`
	UserName        string `json:"userName"`
	ExtTID          string `json:"extTID"`
	ExtMID          string `json:"extMID"`
	MeEmailID       string `json:"meEmailID"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	AccNo           string `json:"accNo"`
	Ifsc            string `json:"ifsc"`
	PspRefNo        string `json:"pspRefNo"`
	VirtualAddress  string `json:"virtualAddress"`
	VAReqType       string `json:"vAReqType"`
	BankName        string `json:"bankName"`
}

type CheckVirtualAdressApiResponse struct {
	Status      string
	StatusDesc  string
	RequestInfo string
	PayeeType   string
}
type Serviceproviderdetail struct {
	Detail          int64
	Servicename     string
	Detail_adm      int64
	Servicename_adm string
}

type SubMerchantonBoardApiResponse struct {
	Status          string
	StatusDesc      string
	Mebussname      string
	PgMerchantID    string
	CrtDate         string
	IntegrationType string
	MerVirtualAdd   string
	LegalStrName    string
	MasterName      string
	AdminName       string
}

type BankRequest struct {
	RequestMsg   string `json:"requestMsg"`
	PgMerchantId string `json:"pgMerchantId"`
}

type RequestInfoMap struct {
	PgMerchantId string `json:"pgMerchantId"`
	PspRefNo     string `json:"pspRefNo"`
}

type PayeeTypeMap struct {
	VirtualAddress string `json:"virtualAddress"`
}

type CheckvirtualAdressinstance struct {
	RequestInfo RequestInfoMap `json:"requestInfo"`
	PayeeType   PayeeTypeMap   `json:"payeeType"`
	VAReqType   string         `json:"vAReqType"`
}
type Upitransactionidcount struct {
	Idcount int64
	//Origin_idntf string
}
type Qr_response struct {
	Status       string
	StatusDesc   string
	TxnId        string
	MerchantId   string
	Amount       string
	Paymentstate string
	PayerVPA     string
	PayeeVPA     string
	Qrdata       string
	Userid       int64
}
