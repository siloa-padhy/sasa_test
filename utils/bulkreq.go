package utils

import "time"

type Bulkreq struct {
	PspId           string `json:"pspId"`
	PgMerchantId    string `json:"PgMerchantId "`
	Mebussname      string `json:"mebussname" validate:"required"`
	LegalStrName    string `json:"legalStrName" validate:"required"`
	MerVirtualAdd   string `json:"merVirtualAdd" validate:"required"`
	Awlmcc          string `json:"awlmcc "`
	StrCntMobile    string `json:"strCntMobile" validate:"required"`
	RequestUrl1     string `json:"requestUrl1"`
	RequestUrl2     string `json:"requestUrl2"`
	MerchantType    string `json:"merchantType "`
	IntegrationType string `json:"integrationType"`
	SettleType      string `json:"settleType"`
	PanNo           string `json:"panNo" validate:"required,len=10,alphanum"`
	UserName        string `json:"userName" validate:"required"`
	ExtTID          string `json:"extTID"`
	ExtMID          string `json:"extMID"`
	MeEmailID       string `json:"meEmailID"`
	FirstName       string `json:"firstName" validate:"required"`
	LastName        string `json:"lastName" validate:"required"`
	AccNo           string `json:"accNo" validate:"required"`
	Ifsc            string `json:"ifsc" validate:"required"`
	PspRefNo        string `json:"pspRefNo"`
	VirtualAddress  string `json:"virtualAddress"`
	VAReqType       string `json:"vAReqType"`
	BankName        string `json:"bankName" validate:"required"`
	User_id         int64  `json:"user_id"`
	Authorization   string `json:"Authorization" validate:"required"`
	Payeevpa        string
}
type Generate struct {
	Authorization string
	Pa            string
	Pn            string
}
type Create struct {
	Fromdate     time.Time
	Todate       time.Time
	Fromindex    string
	Toindex      string
	PgMerchantId string
}
type Getfrmuser struct {
	UserName string ` validate:"required"`
	Pa       string
	Pn       string
}
type Statusen struct {
	PgMerchantId string ` validate:"required"`
	TxId         string
	MerchantType string
	CustRefNo    string
	PspRefNo     string
	NpciTranId   string
}
type Dynamic struct {
	PgMerchantId string
	Resp         string
}
type Static struct {
	PgMerchantId string
	Resp         string
}
type Bankdata struct {
	CustRefNo       string
	OrderNo         string //order number
	Status_callback string
	Payeevpa        string
	Status_desc     string
	MerchantType    string
	Amount          string
	PayerVPA        string
	TxnNote         string
	NpciTransId     string
	UpiTransRefNo   string
	TxnAuthDate     string
	ApprovalNumber  string
	AddInfo         AddInfo
	ResponseCode    string
}
type Bankcount struct {
	Rrn_count       int64
	Id_count        int64
	Payee_vpa_count int64
	Usermap_count   int64
	User_prop_count int64
	Feature_count   int64
}
type Bank_response struct {
	Status     string
	StatusDesc string
	ApiComment string
	CustRefNo  string
	OrderNo    string
	Rec        []string
}
type Data struct {
	User_admin_prop  string
	Username         string
	User_id          int64
	Feature_active   bool
	Admin_proporties string
	Amount           float64
	Amount_str       string
	PgMerchantID     string
	PayeeVPA         string
	Store_name       string
	Merchant_code    string
}
type Wallet struct {
	Status string //temporaryly used for wallet api Status response
}
type Trasactionpojo struct {
	pspRefNo               string
	upiTransRefNo          string
	npciTransId            string
	custRefNo              string
	amount                 string
	txnAuthDate            string
	responseCode           string
	approvalNumber         string
	status                 string
	statusDesc             string
	CommitResponseInfoPojo addInfo
	payerVPA               string
	payeeVPA               string
	orderNo                string
}
type addInfo struct {
	addInfo1 string
	addInfo2 string
	addInfo3 string
	addInfo4 string
}
