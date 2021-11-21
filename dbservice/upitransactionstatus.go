package dbservice

import "time"

type Upitransaction struct {
	Add_info1               string
	Add_info2               string
	Add_info3               string
	Add_info4               string
	Amount                  float64
	Api_service_provider_id int64
	Api_user_id             int64
	Api_wallet_id           int64
	Approval_number         string
	Client_unique_id        string
	Created_date            time.Time
	Credit_debit            bool
	Id                      int64
	Is_auto_settled         string
	Is_sl                   bool
	Merchant_type           string
	Npci_trans_id           string
	Operation_performed     string
	Origin_identifier       string
	Param_a                 string
	Param_b                 string
	Param_c                 string
	Rrn                     string
	Service_provider_id     int64
	Stan                    int
	Status                  string
	Status_desc             string
	Tnx_id_static           string
	Transaction_note        string
	Transaction_status_code string
	Transaction_type        string
	Txn_auth_date           string
	Updated_date            time.Time
	Upi_trans_ref_no        string
	User_id                 int64
	Wallet_id               int64
	//Client_unique_id        int64
}
