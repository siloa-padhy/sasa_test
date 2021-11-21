package utils

// type Checkvpa_request struct {
// 	MerVirtualAdd string
// 	VAReqType     string
// }
type Checkvpa_response struct {
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
type Vpatype struct {
	Vpa_count  int64
	User_count int64
}
type Dynamicres struct {
	PspRefNo           string
	UpiTransRefNo      string
	ApprovalNumber     string
	NpciTransId        string
	CustRefNo          string
	Amount             string
	TxnAuthDate        string
	ResponseCode       string
	Status             string
	CommitResponseInfo AddInfo
	PayerVPA           string
	PayeeVPA           string
	OrderNo            string
	CurrentStatusDesc  string
	TxnNote            string
	TxnType            string
}
type AddInfo struct {
	AddInfo1 string
	AddInfo2 string
	AddInfo3 string
	AddInfo4 string
}
