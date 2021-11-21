package dbservice

import (
	"time"
)

type Indusupiproporties struct {
	UserPropertyKey   IndusUpiPropertiedKey
	Account_number    string
	Bank_name         string
	Bene_name         string
	Beni_status       string
	Created_date      time.Time
	Direct_payee_vpa1 string
	Dynamic_wallet1   int64
	Ifsc_code         string
	IsSl              bool
	Merchant_code     string
	Merchant_key      string
	Mobile_number     string
	Onboarding_status string
	Payee_vpa         string
	Pg_merchant_id    string
	Static_wallet1    int64
	Status            string
	Status_desc       string
	Store_name        string
	Sub_status        string
	Upi_update_status string
	Upi_update_date   time.Time
	User_id           int64
	User_name         string
}

type IndusUpiPropertiedKey struct {
	UserName string
	PayeeVPA string
}
type Listkeys struct {
	Status           string
	TransDetails     []string
	PeginationConfig string
	Statusdesc       string
	Pgmerchantid     string
	Ismerchant       string
}
type Pojo2 struct {
	Status      string
	StatusDesc  string
	RequestInfo string
	PayeeType   string
}
