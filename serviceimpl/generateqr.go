package serviceimpl

import (
	"fmt"
	"log"

	"main.go/utils"
)

func Generateqr(req1 utils.Getfrmuser, s *Server) utils.Response1 {
	resp := utils.Response1{}
	if req1.UserName != "" {
		userapimapping, err := Getuser(req1.UserName, s)
		fmt.Println(userapimapping)
		if err == nil {
			userproperties, err := Getuserproperties(userapimapping.User_id, s)
			if err == nil {
				fmt.Println(userproperties)
				var adminname string
				var mastername string

				if userproperties.Admin != 0 {
					userproperties.User_id = userproperties.Admin
					err = s.Session.Query(`SELECT user_name FROM user_properties WHERE user_id=? allow filtering`, userproperties.Admin).Scan(&adminname)

					if err != nil {
						fmt.Println(err.Error())
						resp.Status = "-1"
						resp.Message = "Error in Getting Data from DB is " + err.Error()
						return resp
					}

				}
				if userproperties.Master != 0 {

					userproperties.User_id = userproperties.Master
					err = s.Session.Query(`SELECT user_name FROM user_properties WHERE user_id=? allow filtering`, userproperties.Admin).Scan(&mastername)
					fmt.Println(mastername)

					if err != nil {
						fmt.Println(err.Error())
						resp.Status = "-1"
						resp.Message = "Error in Getting Data from DB is " + err.Error()
						return resp
					}

				}
				indusproperties, err := Getindusupi(req1.UserName, s)

				if err == nil {
					fmt.Println(indusproperties.Status)
					fmt.Println(indusproperties.Sub_status)
					if indusproperties.Status == "SUCCESS" && indusproperties.Sub_status == "0" {
						//check feature 53 is active ot not
						var active bool
						err1 := s.Session.Query(`SELECT is_active FROM user_feature WHERE user_id=? allow filtering`, userproperties.User_id).Scan(&active)
						fmt.Println(active)
						fmt.Println(err1)
						if active == true {
							fmt.Println("Generate Qr code")
							//var qr utils.Generate
							Qrcode, err := GetQrdata(req1)
							if err == nil {
								fmt.Println("Qrcode:", Qrcode)
								resp.Qrcode = Qrcode
								resp.Status = "0"
								resp.Message = "QR code generate successfully"
								resp.Vpaid = req1.Pa
								resp.MasterName = mastername
								resp.AdminName = adminname
								resp.Storename = indusproperties.Store_name
								return resp
							} else {
								resp.Qrcode = "Null"
								resp.Status = "1"
								resp.Message = "QR code generation fail"
								log.Println("QR code generation fail")
								return resp
							}

						} else {
							resp.Qrcode = "Null"
							resp.Status = "2"
							resp.Message = "kindly contact admin for feature activation to show Qr code"
							log.Println("kindly contact admin for feature activation to show Qr code")
							return resp
						}

					} else if indusproperties.Status == "INITIATED" && (indusproperties.Sub_status == "0" || indusproperties.Sub_status == "2" || indusproperties.Sub_status == "3") {
						resp.Qrcode = "Null"
						resp.Status = "-1"
						resp.Message = "Kindly register again"
						log.Println("Kindly register again")
						return resp
					}

				} else {
					resp.Qrcode = "Null"
					resp.Status = "-1"
					resp.Message = "user not available"
					log.Println("user not available")

					return resp

				}

			} else {
				resp.Qrcode = "Null"
				resp.Status = "-1"
				resp.Message = "user not exist"
				log.Println("user not exist")
				return resp
			}
		} else {
			resp.Qrcode = "Null"
			resp.Status = "-1"
			resp.Message = "user not exist"
			log.Println("user not exist")
			return resp
		}
	} else {
		resp.Qrcode = "Null"
		resp.Status = "-1"
		resp.Message = "username Not available"
		log.Println("user not available")
		return resp
	}
	return resp
}
