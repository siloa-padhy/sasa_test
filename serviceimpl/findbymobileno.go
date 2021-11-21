package serviceimpl

import (
	"log"

	"main.go/utils"
)

func Findby_mobilenumber(mobile utils.Check_mobile_number, s *Server) utils.Mobile_Response {
	resp := utils.Mobile_Response{}
	mob := utils.Mobiledata{}
	err := s.Session.Query(`select count(*) from indus_upi_properties where mobile_number = ? allow filtering`, mobile.StrCntMobile).Scan(&mob.Mobile_count)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println(err.Error())
		resp.Status = "-1"
		resp.StatusDesc = "Error in getting Mobile Number in DB"
		return resp
	}
	if mob.Mobile_count == 0 {
		resp.Status = "0"
		resp.StatusDesc = "Mobile Number Available "
		return resp
	} else {
		resp.Status = "-1"
		resp.StatusDesc = "Mobile Number Already Taken "
		return resp
	}
}
