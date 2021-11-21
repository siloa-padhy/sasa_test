package serviceimpl

import (
	"fmt"

	"main.go/dbservice"
)

// func Insertuserproperties(responseInit2 dbservice.Userproperties, s *Server) error {
// 	// fmt.Println("Before save gateway transaction id", transdata.Gateway_transaction_id)
// 	err := s.Session.Query("Insert INTO upi.indus_upi_properties(admin,is_upi_settle,master,upi_settle_update,user_id,user_name) VALUES(?,?,?,?,?,?)",

// 		responseInit2.Admin, responseInit2.Is_upi_settle, responseInit2.Master, responseInit2.Upi_settle_update, responseInit2.User_id, responseInit2.User_name).Exec()
// 	if err != nil {
// 		fmt.Println("Error", err.Error())
// 		return err
// 	}
// 	return err
// }
func Updateuserproperties(responseInit2 dbservice.Userproperties, s *Server) error {

	err := s.Session.Query("UPDATE upi.user_properties SET is_upi_settle=?, upi_settle_update=? WHERE user_id=?", responseInit2.Is_upi_settle, responseInit2.Upi_settle_update, responseInit2.User_id).Exec()
	if err != nil {
		fmt.Println("Error", err.Error())

		return err
	}

	return err
}
