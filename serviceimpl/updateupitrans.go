package serviceimpl

import (
	"fmt"

	"main.go/dbservice"
)

//add_info1,add_info2,add_info3,add_info4,//amount,api_service_provider_id,//api_user_id,api_wallet_id,approval_number,client_unique_id,//created_date,credit_debit,is_auto_settled,is_sl,merchant_type,npci_trans_id,operation_performed,origin_identifier,param_a,param_b,param_c,rrn,service_provider_id,//stan,status,status_desc,tnx_id_static,transaction_note,transaction_status_code,//transaction_type,txn_auth_date,updated_date,upi_trans_ref_no,user_id,wallet_id
func Update_upi_transaction_status(response dbservice.Upitransaction, s *Server) error {

	err := s.Session.Query("UPDATE upi_transaction_status SET add_info1=?,add_info2=?,add_info3=?,add_info4=?,amount=?,api_service_provider_id=?,api_user_id=?,api_wallet_id=?,approval_number=?,client_unique_id=?,created_date=?,credit_debit=?,is_auto_settled=?,is_sl=?,merchant_type=?,npci_trans_id=?,operation_performed=?,origin_identifier=?,param_a=?,param_b=?,param_c=?,rrn=?,service_provider_id=?,stan=?,status=?,status_desc=?,tnx_id_static=?,transaction_note=?,transaction_status_code=?,transaction_type=?,txn_auth_date=?,updated_date=?,upi_trans_ref_no=?,user_id=?,wallet_id=? WHERE id=?", response.Add_info1, response.Add_info2, response.Add_info3, response.Add_info4, response.Amount, response.Api_service_provider_id, response.Api_user_id, response.Api_wallet_id, response.Approval_number, response.Client_unique_id, response.Created_date, response.Credit_debit, response.Is_auto_settled, response.Is_sl, response.Merchant_type, response.Npci_trans_id, response.Operation_performed, response.Origin_identifier, response.Param_a, response.Param_b, response.Param_c, response.Rrn, response.Service_provider_id, response.Stan, response.Status, response.Status_desc, response.Tnx_id_static, response.Transaction_note, response.Transaction_status_code, response.Transaction_type, response.Txn_auth_date, response.Updated_date, response.Upi_trans_ref_no, response.User_id, response.Wallet_id, response.Id).Exec()
	if err != nil {
		fmt.Println("Error in updation response data in upi_transaction_status table :", err.Error())
		return err
	} else {
		fmt.Println("Successfully Updateed in upi_transaction_status table ")
	}
	return err
}
func Updateupiproperties(response dbservice.Upitransaction, s *Server) error {

	err := s.Session.Query("UPDATE upi_transaction_status SET status=?, transaction_status_code=? WHERE id=?", response.Status, response.Transaction_status_code, response.Id).Exec()
	if err != nil {
		fmt.Println("Error", err.Error())
		return err
	}
	return err
}

//Below two for static callback
func Updateupiproperties_for_static(response dbservice.Upitransaction, s *Server) error {
	if response.Is_sl {
		err := s.Session.Query("UPDATE upi_transaction_status SET status=?, transaction_status_code=?,status_desc=?,is_sl=? WHERE id=?", response.Status, response.Transaction_status_code, response.Status_desc, response.Is_sl, response.Id).Exec()
		if err != nil {
			fmt.Println("Error", err.Error())
			return err
		}
		return err
	} else {
		err := s.Session.Query("UPDATE upi_transaction_status SET status=?, transaction_status_code=?,status_desc=? WHERE id=?", response.Status, response.Transaction_status_code, response.Status_desc, response.Id).Exec()
		if err != nil {
			fmt.Println("Error", err.Error())
			return err
		}
		return err
	}
}
func Do_wallet_operation() error {
	//All Wallet OPeration by calling wallet api are do here
	return err
}
