package serviceimpl

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"main.go/dbservice"
	"main.go/utils"
)

func Upitransaction(transaction dbservice.Upitransaction, s *Server) (utils.Qr_response, error) {
	resp := utils.Qr_response{}
	transid := utils.Upitransactionidcount{}
	err = s.Session.Query(`select count(*) from upi_transaction_status where id = ? allow filtering`, transaction.Id).Scan(&transid.Idcount)
	if err != nil {
		resp.Status = "-1"
		resp.StatusDesc = "Error in Getting Data from DB " + err.Error()
		return resp, nil
	}
	if transid.Idcount == 0 {
		if transaction.Origin_identifier != "" {
			err := s.Session.Query("Insert INTO upi.upi_transaction_status(id,amount,api_service_provider_id,api_user_id,api_wallet_id,client_unique_id,created_date,credit_debit,is_sl,merchant_type,operation_performed,origin_identifier,param_a,param_b,param_c,service_provider_id,status,transaction_note,transaction_status_code,transaction_type,updated_date,user_id,wallet_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
				transaction.Id, transaction.Amount, transaction.Api_service_provider_id, transaction.Api_user_id, transaction.Api_wallet_id, transaction.Client_unique_id, transaction.Created_date, transaction.Credit_debit, transaction.Is_sl, transaction.Merchant_type, transaction.Operation_performed, transaction.Origin_identifier, transaction.Param_a, transaction.Param_b, transaction.Param_c, transaction.Service_provider_id, transaction.Status, transaction.Transaction_note, transaction.Transaction_status_code, transaction.Transaction_type, transaction.Updated_date, transaction.User_id, transaction.Wallet_id).Exec()
			if err != nil {
				resp.Status = "-1"
				resp.StatusDesc = "Data Insertion Error  " + err.Error()
				return resp, nil
			} else {
				fmt.Println("Data initiated Successfully with id  :", transaction.Id)
				resp.Status = "10"
				resp.StatusDesc = "Data initiated Successfully  "
				resp.Userid = transaction.Id
				resp.PayerVPA = transaction.Origin_identifier
				resp.MerchantId = transaction.Client_unique_id
				resp.TxnId = strconv.FormatInt(transaction.Id, 10)
				resp.Amount = strconv.FormatFloat(transaction.Amount, 'f', 1, 32)
				resp.Paymentstate = "INITIATED"
				return resp, nil
			}
		} else {
			err := s.Session.Query("Insert INTO upi.upi_transaction_status(id,amount,api_service_provider_id,api_user_id,api_wallet_id,client_unique_id,created_date,credit_debit,is_sl,merchant_type,operation_performed,param_a,param_b,param_c,service_provider_id,status,transaction_note,transaction_status_code,transaction_type,updated_date,user_id,wallet_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
				transaction.Id, transaction.Amount, transaction.Api_service_provider_id, transaction.Api_user_id, transaction.Api_wallet_id, transaction.Client_unique_id, transaction.Created_date, transaction.Credit_debit, transaction.Is_sl, transaction.Merchant_type, transaction.Operation_performed, transaction.Param_a, transaction.Param_b, transaction.Param_c, transaction.Service_provider_id, transaction.Status, transaction.Transaction_note, transaction.Transaction_status_code, transaction.Transaction_type, transaction.Updated_date, transaction.User_id, transaction.Wallet_id).Exec()
			if err != nil {
				resp.Status = "-1"
				resp.StatusDesc = "Data Insertion Error  " + err.Error()
				return resp, nil
			} else {

				fmt.Println("Data initiated Successfully with id  :", transaction.Id)
				resp.Status = "1"
				resp.StatusDesc = "Data initiated Successfully "
				resp.PayerVPA = fmt.Sprint(nil)
				return resp, nil
			}

		}

	} else {
		resp.Status = "-1"
		resp.StatusDesc = "Transaction id is already present in db"
	}
	return resp, err
}

//Full table data are :
//"Insert INTO upi.upi_transaction_status(id,add_info1,add_info2,add_info3,add_info4,amount,api_service_provider_id,api_user_id,api_wallet_id,approval_number,client_unique_id,created_date,credit_debit,id,is_auto_settled,is_sl,merchant_type,npci_trans_id,operation_performed,origin_identifier,param_a,param_b,param_c,rrn,service_provider_id,stan,status,status_desc,tnx_id_static,transaction_note,transaction_status_code,transaction_type,txn_auth_date,updated_date,upi_trans_ref_no,user_id,wallet_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
//Insert INTO upi.upi_transaction_status(id,add_info1,add_info2,add_info3,add_info4,amount,api_user_id,approval_number,client_unique_id,created_date,credit_debit,is_auto_settled,is_sl,merchant_type,operation_performed,origin_identifier,param_a,param_b,param_c,rrn,service_provider_id,status,status_desc,transaction_note,transaction_status_code,transaction_type,txn_auth_date,updated_date,upi_trans_ref_no,user_id,wallet_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
//before
//operation_performed,origin_identifier,param_a,param_b,param_c,rrn,service_provider_id,stan,status,status_desc,tnx_id_static,transaction_note,transaction_status_code,transaction_type,txn_auth_date,updated_date,upi_trans_ref_no,user_id,wallet_id
func Select_upi_transaction_status(bankdata utils.Bankdata, s *Server) (dbservice.Upitransaction, error) {
	var upi_transaction_status []dbservice.Upitransaction
	m := map[string]interface{}{}
	fmt.Println(bankdata.OrderNo)
	query := "select * from upi_transaction_status where id = ?"
	plandata := s.Session.Query(query, bankdata.OrderNo).Iter()
	for plandata.MapScan(m) {

		upi_transaction_status = append(upi_transaction_status, dbservice.Upitransaction{
			Id:                      m["id"].(int64),
			Add_info1:               m["add_info1"].(string),
			Add_info2:               m["add_info2"].(string),
			Add_info3:               m["add_info3"].(string),
			Add_info4:               m["add_info4"].(string),
			Amount:                  m["amount"].(float64),
			Api_service_provider_id: m["api_service_provider_id"].(int64),
			Api_user_id:             m["api_user_id"].(int64),
			Api_wallet_id:           m["api_wallet_id"].(int64),
			Approval_number:         m["approval_number"].(string),
			Client_unique_id:        m["client_unique_id"].(string),
			Created_date:            m["created_date"].(time.Time),
			Credit_debit:            m["credit_debit"].(bool),
			Is_auto_settled:         m["is_auto_settled"].(string),
			Is_sl:                   m["is_sl"].(bool),
			Merchant_type:           m["merchant_type"].(string),
			Npci_trans_id:           m["npci_trans_id"].(string),
			Operation_performed:     m["operation_performed"].(string),
			Origin_identifier:       m["origin_identifier"].(string),
			Param_a:                 m["param_a"].(string),
			Param_b:                 m["param_b"].(string),
			Param_c:                 m["param_c"].(string),
			Rrn:                     m["rrn"].(string),
			Service_provider_id:     m["service_provider_id"].(int64),
			Stan:                    m["stan"].(int),
			Status:                  m["status"].(string),
			Status_desc:             m["status"].(string),
			Tnx_id_static:           m["tnx_id_static"].(string),
			Transaction_note:        m["transaction_note"].(string),
			Transaction_status_code: m["transaction_status_code"].(string),
			Transaction_type:        m["transaction_type"].(string),
			Txn_auth_date:           m["txn_auth_date"].(string),
			Updated_date:            m["updated_date"].(time.Time),
			Upi_trans_ref_no:        m["upi_trans_ref_no"].(string),
			User_id:                 m["user_id"].(int64),
			Wallet_id:               m["wallet_id"].(int64),
		})
	}

	fmt.Println("upi_transaction_status  details", upi_transaction_status)
	if len(upi_transaction_status) == 0 {
		var emptyplan dbservice.Upitransaction
		return emptyplan, errors.New("no Data Available for User apimapping")
	}

	return upi_transaction_status[0], nil

}

//"Insert INTO upi.upi_transaction_status(txn_auth_date,updated_date,upi_trans_ref_no,user_id,wallet_id)
//VALUES(?,?,?,?,?,?)"

func Insert_upi_transaction_status(response dbservice.Upitransaction, s *Server) error {
	query := "Insert INTO upi.upi_transaction_status(id,add_info1,add_info2,add_info3,add_info4,amount,api_service_provider_id,api_user_id,api_wallet_id,approval_number,client_unique_id,created_date,credit_debit,is_auto_settled,is_sl,merchant_type,npci_trans_id,operation_performed,origin_identifier,param_a,param_b,param_c,rrn,service_provider_id,stan,status,status_desc,tnx_id_static,transaction_note,transaction_status_code,transaction_type,txn_auth_date,updated_date,upi_trans_ref_no,user_id,wallet_id) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	err := s.Session.Query(query, response.Id, response.Add_info1, response.Add_info2, response.Add_info3, response.Add_info4, response.Amount, response.Api_service_provider_id, response.Api_user_id, response.Api_wallet_id, response.Approval_number, response.Client_unique_id, response.Created_date, response.Credit_debit, response.Is_auto_settled, response.Is_sl, response.Merchant_type, response.Npci_trans_id, response.Operation_performed, response.Origin_identifier, response.Param_a, response.Param_b, response.Param_c, response.Rrn, response.Service_provider_id, response.Stan, response.Status, response.Status_desc, response.Tnx_id_static, response.Transaction_note, response.Transaction_status_code, response.Transaction_type, response.Txn_auth_date, response.Updated_date, response.Upi_trans_ref_no, response.User_id, response.Wallet_id)
	if err != nil {
		fmt.Println("Data Insertion error in initiation part")
	} else {
		fmt.Println("Init success in db query")
	}
	return nil
}
