package serviceimpl

import (
	"fmt"
	"log"
	"time"

	"main.go/dbservice"
	"main.go/utils"
)

var err error

func Selfonboard(value utils.Selfonboard, s *Server) utils.Response {
	resp := utils.Response{}
	mobile_number := utils.Mobile{}
	//check mobile number in indus_upi_properties table
	err = s.Session.Query(`select count(*) from indus_upi_properties where mobile_number = ? allow filtering`, value.StrCntMobile).Scan(&mobile_number.Mobileno)
	if err != nil {
		fmt.Println(err.Error())
		log.Println("Error in check mobile number in indus_upi_properties ", err)
		resp.Status = -1
		resp.StatusDesc = "Error in Getting Mobile Number from indus_upi_proporties table " + err.Error() //change
		return resp
	}
	//check mobile number exist in indus_upi_properties or not
	if mobile_number.Mobileno == 0 {
		username := utils.Getusername(value.Authorization) //Finding username from Authorization key
		//fmt.Println(username)
		var usercount int
		if username != "" { //to check username exist or not
			err = s.Session.Query(`select count(*) from indus_upi_properties where user_name = ? allow filtering`, username).Scan(&usercount)
			if err != nil {
				fmt.Println(err.Error())
				log.Println("Error in to check username exist or not", err)
				//log.Println(err.Error())
				resp.Status = -1
				resp.StatusDesc = "Error in Getting username from indus_upi_properties table " + err.Error()
				return resp
			}
			if usercount == 1 {
				resp.Status = -1
				resp.StatusDesc = "user already available"
				return resp

			} else {
				//To generate unique id
				uniqueid := Genereteid()
				req := utils.SubMerchantOnBoardApiRequest{}
				req.PspRefNo = fmt.Sprint(uniqueid)
				req.PgMerchantId = "INDB000000384512"
				req.VirtualAddress = value.VirtualAddress
				req.VAReqType = "R"
				apiResponsedata := UpiCheckVirtualAdressCheck2(req) //
				fmt.Println(apiResponsedata)                        //Need to change for call function encryption
				if apiResponsedata.Status == "VN" {
					//Code
					userapimapping, err := Getuser(username, s)
					if err == nil {
						fmt.Println(userapimapping.User_id)

						userproperties, err := Getuserproperties(userapimapping.User_id, s)
						fmt.Println(userproperties) //Print coz use nowhere
						if err == nil {
							value.PgMerchantId = "INDB000000349936"
							value.Awlmcc = "6012"
							value.RequestUrl1 = "https://indusindtest.iserveu.online/indus/commitResponse"
							value.RequestUrl2 = "https://indusindtest.iserveu.online/indus/commitResponse"
							value.MerchantType = "AGGMER"
							value.IntegrationType = "WEBAPI"
							value.SettleType = "NET"
							responseInit := dbservice.Indusupiproporties{}
							indusPropertyKeyInit := dbservice.IndusUpiPropertiedKey{}
							indusPropertyKeyInit.PayeeVPA = value.MerVirtualAdd
							indusPropertyKeyInit.UserName = username
							responseInit.UserPropertyKey = indusPropertyKeyInit
							responseInit.Merchant_code = "6012"
							responseInit.Merchant_key = "f06e573cbbaeb2757bcc953fe7fd7933"
							responseInit.Store_name = value.LegalStrName
							responseInit.User_id = userapimapping.User_id
							responseInit.Created_date = time.Now() //Changes as per format
							responseInit.Account_number = value.AccNo
							responseInit.Mobile_number = value.StrCntMobile
							responseInit.Bene_name = value.FirstName + " " + value.LastName
							responseInit.Ifsc_code = value.Ifsc
							responseInit.Bank_name = value.BankName //Need to null check
							responseInit.Status = "INITIATED"
							responseInit.Sub_status = "0"
							responseInit.User_name = username

							err = Insertindusupi(responseInit, s)
							if err != nil {
								responseInit.Status = "INITIATED"
								responseInit.Sub_status = "2"
								err = Insertindusupi(responseInit, s)
								if err != nil {
									fmt.Println(err.Error())
									log.Println("Error in Insertindusupi(responseInit, s)", err)
									//log.Println(err.Error())
									resp.Status = -1
									resp.StatusDesc = "Exception in insert data"
									return resp
								}

							}

							//FUnction for Bank call
							value2 := utils.Bulkreq{}

							apiResponse := AggregateMerchantOnBoarding2(value2)

							fmt.Println("API Response", apiResponse)
							if apiResponse.MerVirtualAdd != "" {
								indusPropertyKey := dbservice.IndusUpiPropertiedKey{}
								indusPropertyKey.PayeeVPA = apiResponse.MerVirtualAdd
								indusPropertyKey.UserName = username
								responseInit.UserPropertyKey = indusPropertyKey
								responseInit.Merchant_code = "6012"
								responseInit.Merchant_key = "f06e573cbbaeb2757bcc953fe7fd7933"
								responseInit.Pg_merchant_id = apiResponse.PgMerchantID
								responseInit.Store_name = apiResponse.LegalStrName
								responseInit.User_id = userapimapping.User_id
								responseInit.Created_date = time.Now() //chnage
								responseInit.Account_number = value.AccNo
								responseInit.Mobile_number = value.StrCntMobile
								responseInit.Bene_name = value.FirstName + " " + value.LastName
								responseInit.Ifsc_code = value.Ifsc
								responseInit.Bank_name = value.BankName
								responseInit.Status = "SUCCESS"
								responseInit.Sub_status = "0"

								//Save in db

								userproperties.Is_upi_settle = "0"
								userproperties.Upi_settle_update = time.Now().AddDate(0, 0, -1)

								updateerr := Updateuserproperties(userproperties, s)
								if updateerr != nil {
									log.Println("Error in Updateuserproperties(userproperties, s) ", err)
									resp.Status = -1
									resp.StatusDesc = "Exception in Update data in user properties"
									return resp
								}

								err = Insertindusupi(responseInit, s)
								if err != nil {
									fmt.Println(err.Error())
									log.Println("Error in Insertindusupi(responseInit, s) ", err)
									resp.Status = -1
									resp.StatusDesc = "Exception in insert data"
									return resp
								}
								// timedata := time.Now()
								// timedata.Format("02-JUN-21")
								//timedata := time.Now().Format("02-JUN-2006")
								resp.Status = 0
								resp.StatusDesc = "user onBoard sucess"
								resp.Mebussname = value.Mebussname
								resp.PgMerchantID = value.PgMerchantId //not from request , but from bank side//apiresponse
								resp.CrtDate = time.Now().Format("02-JUN-2006")
								resp.IntegrationType = value.IntegrationType
								resp.MerVirtualAdd = value.MerVirtualAdd
								resp.LegalStrName = value.LegalStrName

								return resp

							} else {
								//Else code
								//fmt.Println("Not save for success")

								responseInit.Status = "INITIATED"
								responseInit.Sub_status = "3"
								responseInit.Status_desc = apiResponse.StatusDesc
								err = Insertindusupi(responseInit, s)
								if err != nil {
									fmt.Println(err.Error())
									log.Println("Error in Insertindusupi(responseInit, s) in err", err)
									resp.Status = -1
									resp.StatusDesc = apiResponse.StatusDesc
									return resp
								}
								resp.Status = -1
								resp.StatusDesc = apiResponse.StatusDesc
								return resp

							}

						} else {
							fmt.Println(err.Error())
							log.Println("UUser not available in Proporty table", err)
							resp.Status = -1
							resp.StatusDesc = "User not available in Proporty table" + err.Error()
							return resp
						}

					} else {
						fmt.Println(err.Error())
						log.Println("User not saved due to db error1", err)
						//	log.Println(err.Error())
						resp.Status = -1
						resp.StatusDesc = err.Error()
						return resp
					}

				} else if apiResponsedata.Status == "0" {
					resp.Status = -1
					resp.StatusDesc = "Vpa already taken"
					return resp
				} else {
					//fmt.Println("some issue in fetching ")
					//log.Println(err.Error())
					resp.Status = -1
					resp.StatusDesc = "Some issue in fetching VPA kindly contact admin"
					return resp
				}

			}

		} else {
			//fmt.Println(err.Error())
			log.Println("user Not availeble With This Names", err)
			//log.Println(err.Error())
			resp.Status = -1
			resp.StatusDesc = "user Not availeble With This Names"
			return resp
		}
	} else {
		//fmt.Println("Mobile Number Already Registerd")
		log.Println("Error in check mobile number exist in indus_upi_properties", err)
		resp.Status = -1
		resp.StatusDesc = "Onboarding Failed Mobile Number Already Taken " //change
		return resp
	}
	//return resp
}
