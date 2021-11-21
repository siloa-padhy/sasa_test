package serviceimpl

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"main.go/dbservice"
	"main.go/utils"
)

const projectID = "iserveustaging"
const topicID = "upi_topic_staging"

func Static_call(data utils.Static, s *Server) utils.Bank_response {
	// Check meres from request param as not null
	// then check pgMerchantId and resp as not null
	//status=1,statusDesc=Given information is not correct. Errors:// for both

	// decrypt, err := decrypt([]byte(data.Resp), "b27e4ddb1fcb9e3ae394b13a25155ac1")
	// if err != nil {
	// 	fmt.Println("DEcrypt error ", err)
	// }
	// fmt.Println("decrypt", decrypt)
	//next steps
	//var err error
	fmt.Println("Static call api ")
	userdata := dbservice.Userapimap{}
	resp := utils.Bank_response{}
	count := utils.Bankcount{}
	bank := utils.Bankdata{}
	user := utils.Data{}

	bank.CustRefNo = "11056921323" //"110569213605" //"11056921323" //Customer Reference No.
	// bank.OrderNo = "82700956932"
	bank.OrderNo = "832518461920903168" //"832519632106559" //Order No.
	bank.Status_callback = "Success"
	bank.Payeevpa = "abcxyzshop@indus"
	bank.Amount = "10"
	bank.MerchantType = "DIRECT"
	bank.Status_desc = "Problem in transaction" //that got from callback request
	myquery := `select count(*) from upi_transaction_status where rrn = ? allow filtering`
	err = s.Session.Query(myquery, bank.CustRefNo).Scan(&count.Rrn_count) //this count can be convert to string
	if err != nil {
		resp.Status = "-1"
		resp.StatusDesc = "Error in Getting Rrn from upi_transaction_status table" + err.Error()
		return resp
	}
	fmt.Println("Rrn Count", count.Rrn_count)
	if count.Rrn_count == 0 { //if rrn not found
		//convert to long
		// float_value := utils.Amount{}
		// float_value.OrderNo_in_float, err = strconv.ParseInt(bank.OrderNo, 16, 64) //float_value.OrderNo_in_float
		// if err != nil {
		// 	fmt.Println("conversion error in order no")
		// }
		err = s.Session.Query(`select count(*) from upi_transaction_status where id = ? allow filtering`, bank.OrderNo).Scan(&count.Id_count)
		if err != nil {
			resp.Status = "-1"
			resp.StatusDesc = "Error in Getting id from upi_transaction_status table " + err.Error()
			return resp
		}
		fmt.Println("Id found ", count.Id_count)
		if count.Id_count == 0 {
			err = s.Session.Query(`select count(*) from indus_upi_properties where payee_vpa = ? allow filtering`, bank.Payeevpa).Scan(&count.Payee_vpa_count)
			if err != nil {
				resp.Status = "-1"
				resp.StatusDesc = "Error in Getting payee_vpa from indus_upi_proporties table " + err.Error()
				return resp
			}
			fmt.Println(count.Payee_vpa_count)
			if count.Payee_vpa_count != 0 {
				err = s.Session.Query(`select user_name,user_id,pg_merchant_id,payee_vpa,store_name,merchant_code from indus_upi_properties where payee_vpa = ? allow filtering`, bank.Payeevpa).Scan(&user.Username, &user.User_id, &user.PgMerchantID, &user.PayeeVPA, &user.Store_name, &user.Merchant_code)
				if err != nil {
					resp.Status = "-1"
					resp.StatusDesc = "Error in Getting username and userid from payee_vpa of indus_upi_proporties table " + err.Error()
					return resp
				}
				fmt.Println(user.Username, user.User_id)
				//next in user_feature table by user_id
				//select * from featureid =53
				//select count(*) from user_feature where user_id = ? and feature_id = 53 allow filtering
				//if found is_active
				//if not , statusdesc=Service is not available to this user for payeeVpa
				err = s.Session.Query(`select count(*) from user_feature where user_id = ? and feature_id = 53 allow filtering`, user.User_id).Scan(&count.Feature_count)
				if err != nil {
					resp.Status = "-1"
					resp.StatusDesc = "Error in checking active status from user_id in  user_feature table " + err.Error()
					return resp
				}
				if count.Feature_count != 0 {
					err = s.Session.Query(`select is_active from user_feature where user_id = ? and feature_id = 53 allow filtering`, user.User_id).Scan(&user.Feature_active)
					if err != nil {
						resp.Status = "-1"
						resp.StatusDesc = "Error in checking active status from user_id in  user_feature table " + err.Error()
						return resp
					}
					fmt.Println(user.Feature_active) //if data found and is active = false ,//

					// if user.Feature_active !=nil && user.Feature_active=true
					//if false or nil , then //"Service is not available to this user for payeeVpa i.e blocked "
					a := strconv.FormatBool(user.Feature_active)
					if user.Feature_active && a != "" { //check user details in userapi_mapping//only user id
						err = s.Session.Query(`select count(*) from user_api_mapping where user_name = ?  allow filtering`, user.Username).Scan(&count.Usermap_count)
						if err != nil {
							resp.Status = "-1"
							resp.StatusDesc = "Error in Getting details  from user_api_mapping table " + err.Error()
							return resp
						}
						fmt.Println(count.Usermap_count)
						if count.Usermap_count != 0 {
							err = s.Session.Query(`select count(*) from user_properties where user_id = ? allow filtering`, user.User_id).Scan(&count.User_prop_count)
							if err != nil {
								resp.Status = "-1"
								resp.StatusDesc = "Error in Getting details  from user_properties table " + err.Error()
								return resp
							}
							fmt.Println(count.User_prop_count)
							if count.User_prop_count != 0 {
								err = s.Session.Query(`select admin from user_properties where user_id = ? allow filtering`, user.User_id).Scan(&user.Admin_proporties)
								if err != nil {
									resp.Status = "-1"
									resp.StatusDesc = "Error in Getting details  from user_api_mapping table " + err.Error()
									return resp
								}
								fmt.Println(user.Admin_proporties)
								err = s.Session.Query(`select user_name from user_properties where user_id = ? allow filtering`, user.Admin_proporties).Scan(&user.User_admin_prop)
								if err != nil {
									resp.Status = "-1"
									resp.StatusDesc = "Error in Getting details  from user_proporties table " + err.Error()
									return resp
								}
								err = s.Session.Query(`SELECT api_user_id,wallet_id,wallet2_id,api_wallet2_id FROM user_api_mapping WHERE user_name = ? allow filtering`, user.Username).Scan(&userdata.Api_user_id, &userdata.Wallet_id, &userdata.Wallet2_id, &userdata.Api_wallet2_id)
								if err != nil { //not null check
									fmt.Println(err.Error())
									resp.Status = "-1"
									resp.StatusDesc = "Error in getting Wallet Data from DB as Wallets are not properly mapped for UPI transaction "
									return resp
								}
								//if any one of them
								fmt.Println("No error in above")
								if err == nil {
									amount := utils.Amount{}
									service := utils.Serviceproviderdetail{}
									if bank.Amount != "" { //or amount can not conver to double
										var err1 error
										amount.Amountindouble, err = strconv.ParseFloat(bank.Amount, 64)
										if err != nil {
											fmt.Println("Amount can not conver to float", err)
											resp.Status = "-1"
											resp.StatusDesc = "Error in convering amount to float "
											return resp
										}
										servicetype.Upi_adm = "UPI_ADM"
										servicetype.Upi_api = "UPI_API"
										Query_for_upi_api := ("Select id,service_provider_name from service_provider_details where service_provider_type = ? and min<=? and max>=?  ALLOW FILTERING")
										Query_for_upi_adm := ("Select id,service_provider_name from service_provider_details where service_provider_type = ? and min<=? and max>=?  ALLOW FILTERING")
										err = s.Session.Query(Query_for_upi_api, servicetype.Upi_api, amount.Amountindouble, amount.Amountindouble).Consistency(gocql.One).Scan(&service.Detail, &service.Servicename)
										err1 = s.Session.Query(Query_for_upi_adm, servicetype.Upi_adm, amount.Amountindouble, amount.Amountindouble).Consistency(gocql.One).Scan(&service.Detail_adm, &service.Servicename_adm)
										if err != nil || err1 != nil {
											response := dbservice.Upitransaction{}
											response.Amount = amount.Amountindouble
											response.Origin_identifier = bank.PayerVPA
											response.Transaction_note = bank.TxnNote
											response.Npci_trans_id = bank.NpciTransId
											response.Upi_trans_ref_no = bank.UpiTransRefNo
											response.Rrn = bank.CustRefNo
											response.Txn_auth_date = bank.TxnAuthDate
											response.Tnx_id_static = bank.OrderNo
											response.Approval_number = bank.ApprovalNumber
											response.Add_info1 = bank.AddInfo.AddInfo1
											response.Add_info2 = bank.AddInfo.AddInfo2
											response.Add_info3 = bank.AddInfo.AddInfo3
											response.Add_info4 = bank.AddInfo.AddInfo4
											response.Client_unique_id = user.PgMerchantID
											response.Param_a = user.PayeeVPA
											response.Param_b = user.Store_name
											response.Param_c = user.Merchant_code
											response.Status = "INITIATED"
											response.Transaction_status_code = "3"
											response.Operation_performed = "QR_STATIC"
											response.Transaction_type = "UPI"
											response.Merchant_type = "AGGREGATE"
											response.Status_desc = "Service Provder Not Mapped"
											response.Credit_debit = true
											response.User_id = user.User_id
											response.Api_user_id = userdata.Api_user_id
											response.Wallet_id = userdata.Wallet2_id
											response.Api_wallet_id = userdata.Api_wallet2_id
											if user.User_admin_prop == "iserveu" || user.User_admin_prop == "iserveu2" {
												response.Is_sl = true
											} else {
												response.Is_sl = false
											}
											TxnId := Genereteid()
											fmt.Println(TxnId)
											if TxnId == 0 {
												resp.Status = "1"
												resp.StatusDesc = "TxnId unable to generate"
												return resp
											} else {
												response.Id = TxnId
											}
											response.Created_date = time.Now().UTC() //2020-08-05 16:13:31
											response.Updated_date = time.Now().UTC()
											err = Insert_upi_transaction_status(response, s) //insert
											if err != nil {
												fmt.Println("Data Initiation error ", err)
											} else {
												resp.Status = "-1"
												resp.StatusDesc = "Service Provder not Mapped"
												return resp
											}
											//Publish to pubsub
											// err = Publishstatusmessage(projectID, topicID,response)

										} else {
											response := dbservice.Upitransaction{}
											response.Amount = amount.Amountindouble
											response.Origin_identifier = bank.PayerVPA
											response.Transaction_note = bank.TxnNote
											response.Npci_trans_id = bank.NpciTransId
											response.Upi_trans_ref_no = bank.UpiTransRefNo
											response.Rrn = bank.CustRefNo
											response.Txn_auth_date = bank.TxnAuthDate
											response.Tnx_id_static = bank.OrderNo //Txn
											response.Approval_number = bank.ApprovalNumber
											response.Add_info1 = bank.AddInfo.AddInfo1
											response.Add_info2 = bank.AddInfo.AddInfo2
											response.Add_info3 = bank.AddInfo.AddInfo3
											response.Add_info4 = bank.AddInfo.AddInfo4
											response.Client_unique_id = user.PgMerchantID
											response.Param_a = user.PayeeVPA
											response.Param_b = user.Store_name
											response.Param_c = user.Merchant_code
											response.Status = "INITIATED"
											response.Transaction_status_code = "1"
											response.Operation_performed = "QR_STATIC"
											response.Transaction_type = "UPI"
											response.Merchant_type = "AGGREGATE"
											response.Status_desc = "Initiated successfully"
											response.Credit_debit = true
											response.User_id = user.User_id
											response.Api_user_id = userdata.Api_user_id
											response.Wallet_id = userdata.Wallet2_id
											response.Api_wallet_id = userdata.Api_wallet2_id
											if user.User_admin_prop == "iserveu" || user.User_admin_prop == "iserveu2" {
												response.Is_sl = true
											} else {
												response.Is_sl = false
											}
											response.Service_provider_id = service.Detail_adm
											response.Api_service_provider_id = service.Detail
											TxnId := Genereteid()
											fmt.Println(TxnId)
											if TxnId == 0 {
												resp.Status = "1"
												resp.StatusDesc = "TxnId unable to generate"
												return resp
											} else {
												response.Id = TxnId
											}
											response.Created_date = time.Now().UTC() //2020-08-05 16:13:31
											response.Updated_date = time.Now().UTC()
											//Publish to pubsub
											// err = Publishstatusmessage(projectID, topicID,response)
											err = Insert_upi_transaction_status(response, s) //insert
											if err != nil {
												//fmt.Println("updated in if id not found part")
												resp.Status = "-1"
												resp.StatusDesc = "Failed to initiate upi transaction due to db error"
												return resp
											} else {
												if bank.Status_callback == "Success" || bank.Status_callback == "SUCCESS" || bank.Status_callback == "S" {
													err = Do_wallet_operation()
													if err == nil {
														//save to db
														response.Status = "SUCCESS"
														response.Transaction_status_code = "0"
														response.Status_desc = "Wallet Credited Successfully"
														response.Updated_date = time.Now()
														//Check for is_sl in func

														// Publishing the STATUS,TxnStatusCode and status_desc to PubSub
														// err = Publishstatusmessage(projectID, topicID,response)

														err = Update_upi_transaction_status(response, s) //update
														if err == nil {
															fmt.Println("Status messages for SUCCESS for wallet operation are Saved in upi_transaction_status")
															// Publishing the STATUS,TxnStatusCode and status_desc to PubSub
															// err = Publishstatusmessage(projectID, topicID,response)
														} else {
															fmt.Println("Update Error in Wallet Credited Successfully and status code 0 part", err)
															resp.Status = "-1"
															resp.StatusDesc = "Failed to Update upi transaction due to db error in Wallet Credited Successfully and status code 0 part"
															return resp
														}

													} else {
														//save to db
														response.Status = "FAILED"
														response.Transaction_status_code = "0"
														response.Status_desc = "Wallet Credit Failed"
														response.Updated_date = time.Now()
														//Check for is_sl in func
														err = Update_upi_transaction_status(response, s) //update
														if err == nil {
															fmt.Println("Status messages for FAILED for wallet operation are Saved in upi_transaction_status")
															// Publishing the STATUS,TxnStatusCode and status_desc to PubSub
															// err = Publishstatusmessage(projectID, topicID,response)
														} else {
															fmt.Println("Update Error in Wallet Credited failed and status code 0 part", err)
															resp.Status = "-1"
															resp.StatusDesc = "Failed to Update upi transaction due to db error in Wallet Credit Failed and status code 0 part"
															return resp
														}
													}

												} else if bank.Status_callback == "FAILED" || bank.Status_callback == "FAIL" || bank.Status_callback == "F" || bank.Status_callback == "Rejected" {
													response.Status = "FAILED"
													response.Transaction_status_code = "6"
													response.Status_desc = bank.Status_desc
													response.Updated_date = time.Now()
													// Publishing the STATUS,TxnStatusCode and status_desc to PubSub
													// err = Publishstatusmessage(projectID, topicID,response)
													err = Update_upi_transaction_status(response, s) //update
													if err == nil {
														fmt.Println("Status messages for FAILED are Saved in upi_transaction_status")
													} else {
														fmt.Println("Update Error in bank.Status_desc and status code 6 part", err)
														resp.Status = "-1"
														resp.StatusDesc = "Failed to Update upi transaction due to db error in bank.Status_desc and status code 6 part"
														return resp
													}
												} else {
													response.Status = "INPROGRESS"
													response.Transaction_status_code = "9"
													response.Status_desc = bank.Status_desc
													response.Updated_date = time.Now()
													// Publishing the STATUS,TxnStatusCode and status_desc to PubSub
													// err = Publishstatusmessage(projectID, topicID,response)
													err = Update_upi_transaction_status(response, s) //update
													if err == nil {
														fmt.Println("Status messages for INPROGRESS are Saved in upi_transaction_status")
													} else {
														fmt.Println("Update Error in bank.Status_desc and status code 6 part", err)
														resp.Status = "-1"
														resp.StatusDesc = "Failed to Update upi transaction in INPROGRESS error in bank.Status_desc and status code 6 part"
														return resp
													}
												}

											}

										}
									} else {
										resp.Status = "-1"
										resp.StatusDesc = "Failed to fetch amount from bank " //+ err.Error()
										return resp
									}
								} else {
									amount := utils.Amount{}
									service := utils.Serviceproviderdetail{}
									var err1 error
									amount.Amountindouble, _ = strconv.ParseFloat(bank.Amount, 64)
									servicetype.Upi_adm = "UPI_ADM"
									servicetype.Upi_api = "UPI_API"
									Query_for_upi_api := ("Select id,service_provider_name from service_provider_details where service_provider_type = ? and min<=? and max>=?  ALLOW FILTERING")
									Query_for_upi_adm := ("Select id,service_provider_name from service_provider_details where service_provider_type = ? and min<=? and max>=?  ALLOW FILTERING")
									err = s.Session.Query(Query_for_upi_api, servicetype.Upi_api, amount.Amountindouble, amount.Amountindouble).Consistency(gocql.One).Scan(&service.Detail, &service.Servicename)
									err1 = s.Session.Query(Query_for_upi_adm, servicetype.Upi_adm, amount.Amountindouble, amount.Amountindouble).Consistency(gocql.One).Scan(&service.Detail_adm, &service.Servicename_adm)
									if err != nil || err1 != nil {
										response := dbservice.Upitransaction{}
										response.Amount = amount.Amountindouble
										response.Origin_identifier = bank.PayerVPA
										response.Transaction_note = bank.TxnNote
										response.Npci_trans_id = bank.NpciTransId
										response.Upi_trans_ref_no = bank.UpiTransRefNo
										response.Rrn = bank.CustRefNo
										response.Txn_auth_date = bank.TxnAuthDate
										response.Tnx_id_static = bank.OrderNo //Txn
										response.Approval_number = bank.ApprovalNumber
										response.Add_info1 = bank.AddInfo.AddInfo1
										response.Add_info2 = bank.AddInfo.AddInfo2
										response.Add_info3 = bank.AddInfo.AddInfo3
										response.Add_info4 = bank.AddInfo.AddInfo4
										response.Client_unique_id = user.PgMerchantID
										response.Param_a = user.PayeeVPA
										response.Param_b = user.Store_name
										response.Param_c = user.Merchant_code
										response.Status = "INITIATED"
										response.Transaction_status_code = "3"
										response.Operation_performed = "QR_STATIC"
										response.Transaction_type = "UPI"
										response.Merchant_type = "AGGREGATE"
										response.Status_desc = "Wallet not properly mapped for upi transaction"
										response.Credit_debit = true
										response.User_id = user.User_id
										response.Api_user_id = userdata.Api_user_id
										response.Wallet_id = userdata.Wallet2_id
										response.Api_wallet_id = userdata.Api_wallet2_id
										if user.User_admin_prop == "iserveu" || user.User_admin_prop == "iserveu2" {
											response.Is_sl = true
										} else {
											response.Is_sl = false
										}
										TxnId := Genereteid()
										fmt.Println(TxnId)
										if TxnId == 0 {
											resp.Status = "1"
											resp.StatusDesc = "TxnId unable to generate"
											return resp
										} else {
											response.Id = TxnId
										}
										response.Created_date = time.Now() //2020-08-05 16:13:31
										response.Updated_date = time.Now()

										//Generate a transaction id , save db , created date and updated date== time.now()
										//Publish to pubsub

										// intid, _ := strconv.ParseInt(bank.OrderNo, 10, 64)
										// response.Id = intid
										// err = Publishstatusmessage(projectID, topicID,response)
										err = Insert_upi_transaction_status(response, s) //insert
										if err != nil {
											fmt.Println("updated in if id not found part")
										} else {
											resp.Status = "-1"
											resp.StatusDesc = "service Provider not Mapped and wallet not mapped properly "
											return resp
										}
									} else {
										fmt.Println("else part")
										// Check service provider != null , then only go to below part
										response := dbservice.Upitransaction{}
										response.Amount = amount.Amountindouble
										response.Origin_identifier = bank.PayerVPA
										response.Transaction_note = bank.TxnNote
										response.Npci_trans_id = bank.NpciTransId
										response.Upi_trans_ref_no = bank.UpiTransRefNo
										response.Rrn = bank.CustRefNo
										response.Txn_auth_date = bank.TxnAuthDate
										response.Tnx_id_static = bank.OrderNo //Txn
										response.Approval_number = bank.ApprovalNumber
										response.Add_info1 = bank.AddInfo.AddInfo1
										response.Add_info2 = bank.AddInfo.AddInfo2
										response.Add_info3 = bank.AddInfo.AddInfo3
										response.Add_info4 = bank.AddInfo.AddInfo4
										response.Client_unique_id = user.PgMerchantID
										response.Param_a = user.PayeeVPA
										response.Param_b = user.Store_name
										response.Param_c = user.Merchant_code
										response.Status = "INITIATED"
										response.Transaction_status_code = "3"
										response.Operation_performed = "QR_STATIC"
										response.Transaction_type = "UPI"
										response.Merchant_type = "AGGREGATE"
										response.Status_desc = "Wallet not properly mapped for upi transaction"
										response.Credit_debit = true
										response.User_id = user.User_id
										response.Api_user_id = userdata.Api_user_id
										response.Wallet_id = userdata.Wallet2_id
										response.Api_wallet_id = userdata.Api_wallet2_id
										if user.User_admin_prop == "iserveu" || user.User_admin_prop == "iserveu2" {
											response.Is_sl = true
										} else {
											response.Is_sl = false
										}
										response.Service_provider_id = service.Detail_adm
										response.Api_service_provider_id = service.Detail
										TxnId := Genereteid()
										fmt.Println(TxnId)
										if TxnId == 0 {
											resp.Status = "1"
											resp.StatusDesc = "TxnId unable to generate"
											return resp
										} else {
											response.Id = TxnId
										}
										response.Created_date = time.Now() //2020-08-05 16:13:31
										response.Updated_date = time.Now()

										err = Insert_upi_transaction_status(response, s) //insert
										if err == nil {
											fmt.Println("updated in if id not found part")
											//Publish to pubsub
											// err = Publishstatusmessage(projectID, topicID,response)
										} else {
											resp.Status = "-1"
											resp.StatusDesc = "Wallet not properly mapped for upi transaction"
											return resp
										}
									}
								}

							} else {
								resp.Status = "-1"
								resp.StatusDesc = "User not mapped with User proporties table"
								return resp
							}
						} else {
							resp.Status = "-1"
							resp.StatusDesc = "User not mapped with UserApiMappingEntity Table"
							return resp
						}
					} else {
						resp.Status = "-1"
						resp.StatusDesc = "Service is not available to this user for payeeVpa i.e blocked"
						return resp
					}
				} else {
					resp.Status = "-1"
					resp.StatusDesc = "Service is not available to this user for payeeVpa "
				}
			} else {
				resp.Status = "-1"
				resp.StatusDesc = "payeeVPA not Exists in  IndusUpiProperties Table " //+ err.Error()
				return resp
			}
		} else {
			//code of if id found in upi_transaction_status
			//Orderid in float64 but not working , string is working

			//select user_id,Txn_Statu_code,status,status_code,status_desc,approvalNo,Npci_txn_id,rrn,addinfo 1,2,3,4 from upi_transaction_status where id=float_value.OrderNo_in_float//bank.OrderNo
			response := dbservice.Upitransaction{}
			response, err = Select_upi_transaction_status(bank, s) //select * from upi_transaction_status table
			if err != nil {
				fmt.Println("select error from upi_transaction_status")
				resp.Status = "-1"
				resp.StatusDesc = "Error in Getting id from upi_transaction_status table " + err.Error()
				return resp
			}
			fmt.Println("Upi transaction status data", response)
			// Select user_name from user_proporties where user_id=user_id of upi_transaction_status //user_name for pubsub
			var userproperties dbservice.Userproperties
			userproperties, err = Getuserpropertiesdata(response.User_id, s) //select * from user_properties table for updating in table
			if err != nil {
				fmt.Println("select error from user_properties")
				resp.Status = "-1"
				resp.StatusDesc = "Error in Getting all data(user_name for use) from user_properties table " + err.Error()
				return resp
			}
			//fmt.Println("Username", userproperties.User_name) //use it in Publish in Pubsub
			fmt.Println("User proporties data", userproperties)
			if bank.Status_callback == "Success" || bank.Status_callback == "SUCCESS" || bank.Status_callback == "S" {
				if response.Status != "SUCCESS" {
					response.Status = "SUCCESS"
					response.Transaction_status_code = "5"
					response.Status_desc = bank.Status_desc
					response.Approval_number = bank.ApprovalNumber
					response.Npci_trans_id = bank.NpciTransId
					response.Rrn = bank.CustRefNo
					response.Add_info1 = bank.AddInfo.AddInfo1
					response.Add_info2 = bank.AddInfo.AddInfo2
					response.Add_info3 = bank.AddInfo.AddInfo3
					response.Add_info4 = bank.AddInfo.AddInfo4
					response.Updated_date = time.Now()
					err = Update_upi_transaction_status(response, s) //update
					if err == nil {
						fmt.Println("updated in if id not found in success part")
					}
					// err = Publishstatusmessage(projectID, topicID,response)
					if bank.MerchantType == "DIRECT" {
						//Do wallet Operation
						err = Do_wallet_operation()
						if err != nil {
							fmt.Println("wallet failed")
							resp.Status = "-1"
							resp.StatusDesc = "Error in Wallet api operation " + err.Error()
							return resp

						} else {
							fmt.Println("Wallet api called success")
						}
					}
				} else {
					stringValue := strconv.FormatInt(response.Id, 10)
					a := "AggregateMerchant CallBackResponse  Already Commited or insert into Queue for id" + stringValue + "  and next commit amount  " + bank.Amount
					fmt.Println(a)
				}
			} else if bank.Status_callback == "FAILED" || bank.Status_callback == "FAIL" || bank.Status_callback == "F" || bank.Status_callback == "Rejected" {
				if response.Status != "SUCCESS" {
					response.Status = "FAILED"
					response.Transaction_status_code = "6"
					response.Status_desc = bank.Status_desc
					response.Rrn = bank.CustRefNo
					response.Add_info1 = bank.AddInfo.AddInfo1
					response.Add_info2 = bank.AddInfo.AddInfo2
					response.Add_info3 = bank.AddInfo.AddInfo3
					response.Add_info4 = bank.AddInfo.AddInfo4
					response.Updated_date = time.Now()
					err = Update_upi_transaction_status(response, s) //update
					if err == nil {
						fmt.Println("updated in if id not found in failed part")
					}
					// err = Publishstatusmessage(projectID, topicID,response)
				} else {
					stringValue := strconv.FormatInt(response.Id, 10)
					a := "AggregateMerchant CallBackResponse  Already Commited or insert into Queue for id" + stringValue + "and next commit amount" + bank.Amount
					fmt.Println(a)
				}
			} else {
				if response.Status != "SUCCESS" {
					response.Status = "INPROGRESS"
					response.Transaction_status_code = "9"
					response.Status_desc = bank.Status_desc
					response.Rrn = bank.CustRefNo
					response.Add_info1 = bank.AddInfo.AddInfo1
					response.Add_info2 = bank.AddInfo.AddInfo2
					response.Add_info3 = bank.AddInfo.AddInfo3
					response.Add_info4 = bank.AddInfo.AddInfo4
					response.Updated_date = time.Now()
					err = Update_upi_transaction_status(response, s) //update
					if err == nil {
						fmt.Println("updated in if id not found in inprogress part")
						// err = Publishstatusmessage(projectID, topicID,response)
					} else {
						resp.Status = "-1"
						resp.StatusDesc = "Error in Update upi_transaction_status table in INPROGRESS and 9 and if id not found part"
						return resp
					}

				} else {
					stringValue := strconv.FormatInt(response.Id, 10)
					a := "AggregateMerchant CallBackResponse  Already Commited or insert into Queue for id" + stringValue + "  and next commit amount  " + bank.Amount
					fmt.Println(a)
				}
			}
		}
	} else {
		resp.Status = "-1"
		resp.StatusDesc = "RRN already exist  in Upi transaction status Table"
		return resp
	}
	return My_response(resp, bank)
}
