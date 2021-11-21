package serviceimpl

import (
	"fmt"
	"time"

	"main.go/dbservice"
)

func Getindusupi(username string, s *Server) (dbservice.Indusupiproporties, error) {
	var indusupiproperties dbservice.Indusupiproporties
	fmt.Println("username", username)
	indusupiproperties, err = Getindusupiproperties(username, s)

	if err != nil {

		return indusupiproperties, err
	}
	fmt.Println(indusupiproperties)
	return indusupiproperties, nil

}
func Getindusupiproperties(username string, s *Server) (dbservice.Indusupiproporties, error) {
	fmt.Println("username", username)

	var indusupiproperties []dbservice.Indusupiproporties
	m := map[string]interface{}{}
	fmt.Println("enter2")
	query := ("select * from indus_upi_properties where user_name = ? ALLOW FILTERING")
	plandata := s.Session.Query(query, username).Iter()
	fmt.Println(plandata.RowData())
	for plandata.MapScan(m) {
		indusupiproperties = append(indusupiproperties, dbservice.Indusupiproporties{
			Account_number:    m["account_number"].(string),
			Bank_name:         m["bank_name"].(string),
			Bene_name:         m["bene_name"].(string),
			Beni_status:       m["beni_status"].(string),
			Created_date:      m["created_date"].(time.Time),
			Direct_payee_vpa1: m["direct_payee_vpa1"].(string),
			Dynamic_wallet1:   m["dynamic_wallet1"].(int64),
			Ifsc_code:         m["ifsc_code"].(string),
			IsSl:              m["isSl"].(bool),
			Merchant_code:     m["merchant_code"].(string),
			Merchant_key:      m["merchant_key"].(string),
			Mobile_number:     m["mobile_number"].(string),
			Onboarding_status: m["onboarding_status"].(string),
			Payee_vpa:         m["payee_vpa"].(string),
			Pg_merchant_id:    m["pg_merchant_id"].(string),
			Static_wallet1:    m["static_wallet1"].(int64),
			Status:            m["status"].(string),
			Status_desc:       m["status_desc"].(string),
			Store_name:        m["store_name"].(string),
			Sub_status:        m["sub_status"].(string),
			Upi_update_status: m["upi_update_status"].(string),
			Upi_update_date:   m["upi_updated_date"].(time.Time),
			User_id:           m["user_id"].(int64),
			User_name:         m["user_name"].(string),
		})
	}

	fmt.Println("indusupiproperties", indusupiproperties)
	if len(indusupiproperties) == 0 {
		var emptyplan dbservice.Indusupiproporties
		return emptyplan, err
	}
	return indusupiproperties[0], nil
}
