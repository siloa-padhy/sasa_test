package serviceimpl

import (
	"fmt"

	"main.go/dbservice"
)

func Insertindusupi(responseInit dbservice.Indusupiproporties, s *Server) error {
	// fmt.Println("Before save gateway transaction id", transdata.Gateway_transaction_id)
	err := s.Session.Query("Insert INTO upi.indus_upi_properties(account_number,bank_name,bene_name,beni_status,created_date,direct_payee_vpa1,dynamic_wallet1,ifsc_code,merchant_code,merchant_key,mobile_number,onboarding_status,payee_vpa,pg_merchant_id,static_wallet1,status,status_desc,store_name,sub_status,upi_update_status,upi_updated_date,user_id,user_name) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		responseInit.Account_number, responseInit.Bank_name, responseInit.Bene_name, responseInit.Beni_status, responseInit.Created_date, responseInit.Direct_payee_vpa1, responseInit.Dynamic_wallet1, responseInit.Ifsc_code, responseInit.Merchant_code, responseInit.Merchant_key, responseInit.Mobile_number, responseInit.Onboarding_status, responseInit.Payee_vpa, responseInit.Pg_merchant_id, responseInit.Static_wallet1, responseInit.Status, responseInit.Status_desc, responseInit.Store_name, responseInit.Sub_status, responseInit.Upi_update_status, responseInit.Upi_update_date, responseInit.User_id, responseInit.User_name).Exec()
	if err != nil {
		fmt.Println("Error", err.Error())
		return err
	}
	fmt.Println("successfully inserted")
	return err
}
