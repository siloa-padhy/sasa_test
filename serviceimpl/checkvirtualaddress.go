package serviceimpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"main.go/dbservice"
	"main.go/utils"
)

func Checkvirtualaddress(bulkreqdata utils.Bulkreq, s *Server) utils.Response {

	resp := utils.Response{}

	fmt.Println(bulkreqdata.Authorization)
	username := utils.Getusername(bulkreqdata.Authorization)
	fmt.Println(username)
	var usercount int
	var induscount int
	//check username access or not
	if username == "" {
		fmt.Println("error is getting username")
		log.Println("error is getting username", err.Error())
		resp.Status = -1
		resp.StatusDesc = "user not exist or some other erroe kindly contact admin"
		return resp
	} else {

		Txn := Genereteid()

		//Getuser(retcomm, s)
		//etuserproperties(retcomm, s)
		fmt.Println("txnid is", Txn)
		//Check transaction id is 0 or not
		if Txn != 0 {
			txn := strconv.FormatInt(Txn, 10)
			req := utils.SubMerchantOnBoardApiRequest{}
			req.PspRefNo = txn
			req.PgMerchantId = "IND3000000384512"
			virtualAddress2 := strings.ToLower(bulkreqdata.VirtualAddress)
			req.VirtualAddress = virtualAddress2
			req.VAReqType = "R"
			apiResponsedata := UpiCheckVirtualAdressCheck2(req)
			//if apiResponsedata==" "{}
			fmt.Println("res", apiResponsedata)
			//apiResponsedata.Status = "VN"
			if apiResponsedata.Status == "VN" {
				//var s *Server
				fmt.Println("usename", username)
				userapimapping, err := Getuser(username, s)
				fmt.Println("enter")
				if err == nil {
					fmt.Println(userapimapping.User_id)
					//check payeevpa availability
					err = s.Session.Query(`select count(*) from indus_upi_properties where payee_vpa = ? allow filtering`, bulkreqdata.MerVirtualAdd).Scan(&induscount)
					fmt.Println("count", induscount)
					if err != nil {
						fmt.Println(err)
						log.Println(err)
						resp.Status = -1
						resp.StatusDesc = "user not exist or some other erroe kindly contact admin" + err.Error()
						return resp
					}
					if induscount == 0 {
						mobile_number := utils.Mobile{}
						//check mobile number exist or not
						err = s.Session.Query(`select count(*) from indus_upi_properties where mobile_number = ? allow filtering`, bulkreqdata.StrCntMobile).Scan(&mobile_number.Mobileno)
						if err != nil {
							fmt.Println(err)
							log.Println(err)
							resp.Status = -1
							resp.StatusDesc = "Error in Getting Data from DB" + err.Error()
							return resp
						}
						if mobile_number.Mobileno == 0 {

							// SimpleDateFormat formatter = new SimpleDateFormat("yyyy-M-dd 05:30:00");
							// String strDate = formatter.format(new Date());
							// Date dateBefore = new SimpleDateFormat("yyyy-M-dd hh:mm:ss").parse(strDate);
							// Date date1 = new Date(dateBefore.getTime() - 1 * 24 * 3600 * 1000);
							// Date date2 = new SimpleDateFormat("yyyy-M-dd hh:mm:ss").parse(strDate);
							//check username already exsit or not

							err = s.Session.Query(`select count(*) from indus_upi_properties where user_name = ? allow filtering`, username).Scan(&usercount)
							if err != nil {
								fmt.Println(err.Error())
								log.Println(err.Error())
								resp.Status = -1
								resp.StatusDesc = "Error in Getting Data from DB" + err.Error()
								return resp
							}
							if usercount == 0 {
								userproperties, err := Getuserproperties(userapimapping.User_id, s)
								fmt.Println(userproperties)
								if err == nil {
									bulkreqdata.PgMerchantId = "INDB000000349936"
									bulkreqdata.Awlmcc = "6012"
									bulkreqdata.RequestUrl1 = "https://indusindtest.iserveu.online/indus/commitResponse"
									bulkreqdata.RequestUrl2 = "https://indusindtest.iserveu.online/indus/commitResponse"
									bulkreqdata.MerchantType = "AGGMER"
									bulkreqdata.IntegrationType = "WEBAPI"
									bulkreqdata.SettleType = "NET"
									PanNo2 := strings.ToUpper(bulkreqdata.PanNo)
									bulkreqdata.PanNo = PanNo2
									responseInit := dbservice.Indusupiproporties{}
									indusPropertyKey := dbservice.IndusUpiPropertiedKey{}
									indusPropertyKey.PayeeVPA = bulkreqdata.MerVirtualAdd
									indusPropertyKey.UserName = username
									responseInit.UserPropertyKey = indusPropertyKey
									responseInit.Merchant_code = "6012"
									responseInit.Merchant_key = "f06e573cbbaeb2757bcc953fe7fd7933"

									responseInit.Store_name = bulkreqdata.LegalStrName
									responseInit.User_id = userapimapping.User_id
									responseInit.Created_date = time.Now() //chnage
									responseInit.Account_number = bulkreqdata.AccNo
									responseInit.Mobile_number = bulkreqdata.StrCntMobile
									responseInit.Bene_name = bulkreqdata.FirstName + " " + bulkreqdata.LastName
									responseInit.Ifsc_code = bulkreqdata.Ifsc
									responseInit.Bank_name = bulkreqdata.BankName
									responseInit.Status = "INITIATED"
									responseInit.Sub_status = "0"
									responseInit.User_name = username
									//responseInit.Payee_vpa = bulkreqdata.Payeevpa
									err = Insertindusupi(responseInit, s)
									if err != nil {
										fmt.Println("error in insert upi")
										responseInit.Status = "INITIATED"
										responseInit.Sub_status = "2"
										err = Insertindusupi(responseInit, s)
										if err != nil {
											resp.Status = -1
											resp.StatusDesc = "Exception in insert data"
											log.Println("exception", err)
											return resp
										}

									}

									apiResponse := AggregateMerchantOnBoarding2(bulkreqdata)
									fmt.Println("API Response", apiResponse) //function for bank call
									//apiResponse := utils.SubMerchantonBoardApiResponse{}

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
										responseInit.Account_number = bulkreqdata.AccNo
										responseInit.Mobile_number = bulkreqdata.StrCntMobile
										responseInit.Bene_name = bulkreqdata.FirstName + " " + bulkreqdata.LastName
										responseInit.Ifsc_code = bulkreqdata.Ifsc
										responseInit.Bank_name = bulkreqdata.BankName
										responseInit.Status = "SUCCESS"
										responseInit.Sub_status = "0"
										//responseInit.User_name = username
										err = Insertindusupi(responseInit, s)

										if err == nil {
											fmt.Println("Enter for update")
											userproperties.Is_upi_settle = "0"
											userproperties.Upi_settle_update = time.Now().AddDate(0, 0, -1)
											// responseInit2 := dbservice.Userproperties{}
											// responseInit2.User_id = userapimapping.User_id

											err = Updateuserproperties(userproperties, s)
											if err == nil {
												resp.Status = 0
												resp.StatusDesc = "User on board sucess"
												resp.IntegrationType = bulkreqdata.IntegrationType
												resp.Mebussname = bulkreqdata.Mebussname
												resp.MerVirtualAdd = bulkreqdata.MerVirtualAdd
												//timedata := time.Now()
												//timedata := time.Now().Format("02-JUN-2006")
												resp.LegalStrName = bulkreqdata.LegalStrName
												resp.PgMerchantID = apiResponse.PgMerchantID
												resp.CrtDate = time.Now().Format("02-JUN-2006")
												//resp.MasterName=
												//resp.AdminName
												//resp.CrtDate = time.Now()

												return resp
											} else {
												resp.Status = -1
												fmt.Println("User not saved due to db error")
												resp.StatusDesc = "User not saved due to db error"
												log.Println("Exception while calling onboard merchant plaease try again", err.Error())
												return resp
											}

										} else {
											resp.Status = -1
											fmt.Println("User not saved due to db error1")
											resp.StatusDesc = "User not saved due to db error1"
											log.Println("User not saved due to db error1", err)
											return resp

										}
									} else {
										responseInit.Status = "INITIATED"
										responseInit.Sub_status = "3"
										responseInit.Status_desc = apiResponse.StatusDesc
										err = Insertindusupi(responseInit, s)
										if err != nil {
											log.Println(err.Error())
											resp.Status = -1
											resp.StatusDesc = apiResponse.StatusDesc
											return resp
										}

										resp.Status = -1
										resp.StatusDesc = "Exception while calling onboard merchant plaease try again"
										fmt.Println("Exception while calling onboard merchant plaease try again")
										log.Println("Exception while calling onboard merchant plaease try again", err)
										//log.Println(err.Error())
										return resp

									}

								} else {

									resp.Status = -1
									fmt.Println("user Not availeble In Table")
									resp.StatusDesc = "user Not availeble In Table"
									log.Println("user  already availeble In table", err)
									//log.Println(err.Error())
									return resp
								}
							} else {
								//log.Println(err.Error())
								resp.Status = -1
								fmt.Println("user  already availeble In table")
								resp.StatusDesc = "user  already availeble In table"
								log.Println("user  already availeble In table", err)

								return resp
							}

						} else {

							fmt.Println("Mobile Number Already Registerd")
							resp.Status = -1
							resp.StatusDesc = "Onboard fail mobile number already taken "
							log.Println("Onboard fail mobile number already taken", err)

							return resp
						}

					} else {
						//log.Println(err.Error())
						resp.Status = -1
						resp.StatusDesc = "VPA already map to another user"
						fmt.Println("VPA already map to another user")
						log.Println("VPA already map to another user", err)

						return resp

					}

				} else {
					//log.Println(err.Error())
					resp.Status = -1
					resp.StatusDesc = "user Not availeble With This Names"
					fmt.Println("user Not availeble With This Names")
					log.Println("user Not availeble With This Names", err)

					return resp
				}

			} else if apiResponsedata.Status == "VE" {
				//log.Println(err.Error())

				resp.Status = -1
				resp.StatusDesc = "VPA unavailable from bank"
				fmt.Println("VPA unavailable from bank")
				log.Print("VPA unavailable from bank", err)
				return resp

			} else {

				resp.Status = -1
				resp.StatusDesc = "Some issue in fetching VPA kindly contact admin"
				fmt.Println("Some issue in fetching VPA kindly contact admin")
				log.Println("Some issue in fetching VPA kindly contact admin", err)
				return resp
			}
		} else {
			//log.Println(err.Error())
			resp.Status = -1
			resp.StatusDesc = "Error in generating unique ID"
			fmt.Println("Error in generating unique ID")
			log.Println("Error in generating unique ID", err)
			return resp
		}

	}
}
func AggregateMerchantOnBoarding2(request utils.Bulkreq) utils.SubMerchantonBoardApiResponse {
	apiresponse := utils.SubMerchantonBoardApiResponse{}
	out, err := json.Marshal(request)
	fmt.Println(err)

	request.MerVirtualAdd = strings.ToLower(request.MerVirtualAdd)
	encryptedRequest, err := encrypt(string(out), "f06e573cbbaeb2757bcc953fe7fd7933")
	if err == nil {
		requestingJson := utils.BankRequest{}
		requestingJson.RequestMsg = encryptedRequest
		requestingJson.PgMerchantId = request.PgMerchantId
		baseUrl := " https://apig.indusind.com/ibl/prod/api/onBoardSubMerchant"
		jsonValue, _ := json.Marshal(request)

		u := bytes.NewReader(jsonValue)

		req, err := http.NewRequest("POST", baseUrl, u)
		if err != nil {
			fmt.Println("Error is req: ", err.Error())
			log.Println("error in send req", err)
			apiresponse.Status = "-1"
			apiresponse.StatusDesc = err.Error()
			return apiresponse

		}
		fmt.Println(req)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		// Do sends an HTTP request and
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error in send req: ", err.Error())
			log.Println("error in send req", err)
			apiresponse.Status = "-1"
			apiresponse.StatusDesc = err.Error()
			return apiresponse
		}
		defer resp.Body.Close()
		fmt.Println(resp)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err.Error())
			apiresponse.Status = "-1"
			apiresponse.StatusDesc = err.Error()
			return apiresponse
		}
		decryptResponse, err := decrypt(body, "f06e573cbbaeb2757bcc953fe7fd7933")
		if err != nil {
			fmt.Println("error in decrypt")
			log.Println(err.Error())
			apiresponse.Status = "-1"
			apiresponse.StatusDesc = err.Error()
			return apiresponse
		}
		json.Unmarshal([]byte(decryptResponse), &apiresponse)
		apiresponse.Status = "0"
		return apiresponse

	} else {
		fmt.Println("error", err.Error())
		apiresponse.Status = "-1"
		apiresponse.StatusDesc = "Exception "
		return apiresponse
	}

}
func UpiCheckVirtualAdressCheck2(req utils.SubMerchantOnBoardApiRequest) utils.CheckVirtualAdressApiResponse {
	fmt.Println("error")
	response := utils.CheckVirtualAdressApiResponse{}
	requestInfoMap := utils.RequestInfoMap{}
	requestInfoMap.PgMerchantId = req.PgMerchantId
	requestInfoMap.PspRefNo = req.PspRefNo
	payeeTypeMap := utils.PayeeTypeMap{}
	payeeTypeMap.VirtualAddress = req.VirtualAddress
	checkvirtualAdressinstance := utils.CheckvirtualAdressinstance{}
	checkvirtualAdressinstance.PayeeType = payeeTypeMap
	checkvirtualAdressinstance.RequestInfo = requestInfoMap
	checkvirtualAdressinstance.VAReqType = req.VAReqType
	out, err := json.Marshal(checkvirtualAdressinstance)
	fmt.Println(err)
	//log.Println(err.Error())
	encryptedRequest, err := encrypt(string(out), "b27e4ddb1fcb9e3ae394b13a25155ac1")
	fmt.Println("Encrypt")
	if err != nil {
		fmt.Println("Error ", err)
		response.Status = "-1"
		response.StatusDesc = err.Error()
		log.Println("error:", err.Error())
		return response
	}
	if encryptedRequest != "" {
		requestingJson := utils.BankRequest{}
		requestingJson.PgMerchantId = req.PgMerchantId
		requestingJson.RequestMsg = encryptedRequest
		baseUrl := "https://apig.indusind.com/ibl/prod/upijson/validateVPAWeb"
		jsonValue, _ := json.Marshal(requestingJson)

		u := bytes.NewReader(jsonValue)
		fmt.Println("")
		req, err := http.NewRequest("POST", baseUrl, u)
		if err != nil {
			fmt.Println("Error is req: ", err)
			log.Println(err.Error())
			response.Status = "-1"
			response.StatusDesc = err.Error()
			return response
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		// Do sends an HTTP request and
		resp, err := client.Do(req)
		if err != nil {
			//fmt.Println(err.)
			fmt.Println("error in send req: ", err.Error())
			response.Status = "-1"
			response.StatusDesc = err.Error()
			log.Println("error in sending req", err.Error())
			return response
		}
		defer resp.Body.Close()
		fmt.Println(resp)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
			log.Println("error:", err.Error())
		}
		decryptResponse, err := decrypt(body, "b27e4ddb1fcb9e3ae394b13a25155ac1")
		if err != nil {
			fmt.Println(err)
			log.Println("error:", err.Error())
		}
		json.Unmarshal([]byte(decryptResponse), &response)
		if response.Status != "" && response.Status == "VE" {
			response.Status = "VN"
			return response
		} else {
			return response
		}
	} else {
		//fmt.Println(err)
		response.Status = "-1"
		response.StatusDesc = err.Error()
		log.Println("error:", err.Error())
		return response
	}

}
