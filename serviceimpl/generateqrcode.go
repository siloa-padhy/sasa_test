package serviceimpl

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"main.go/dbservice"
	"main.go/utils"
)

type Upivalue struct {
	Upi_adm string
	Upi_api string
}

var servicetype Upivalue

// else part should be present
func Generateqrcode(qrvalue utils.Initreq, s *Server) utils.Qr_response {
	resp := utils.Qr_response{}
	user := dbservice.Userapimap{}
	service := utils.Serviceproviderdetail{}
	amount := utils.Amount{}
	user.User_name = utils.Getusername(qrvalue.Authorization)
	//fmt.Println("Username is :", user.User_name)
	if user.User_name != "" {
		err = s.Session.Query(`select count(*) from user_api_mapping where user_name = ? allow filtering`, user.User_name).Scan(&user.User_count)
		//fmt.Println("User count is :", user.User_count)
		if err != nil {
			fmt.Println(err.Error())
			resp.Status = "-1"
			resp.StatusDesc = "Error in Getting username count from DB"
			return resp
		}
		if user.User_count != 0 {
			err = s.Session.Query(`SELECT user_id,api_user_id,wallet_id,wallet2_id,api_wallet_id,api_wallet2_id FROM user_api_mapping WHERE user_name=? allow filtering`, user.User_name).Scan(&user.User_id, &user.Api_user_id, &user.Wallet_id, &user.Wallet2_id, &user.Api_wallet_id, &user.Api_wallet2_id)
			if err != nil {
				fmt.Println(err.Error())
				resp.Status = "-1"
				resp.StatusDesc = "Error in getting Wallet Data from DB as Wallets are not properly mapped for UPI transaction "
				return resp
			} else {
				var err1 error
				amount.Amountindouble, _ = strconv.ParseFloat(qrvalue.Amount, 64)
				servicetype.Upi_adm = "UPI_ADM"
				servicetype.Upi_api = "UPI_API"
				Query_for_upi_api := ("Select id,service_provider_name from service_provider_details where service_provider_type = ? and min<=? and max>=?  ALLOW FILTERING")
				Query_for_upi_adm := ("Select id,service_provider_name from service_provider_details where service_provider_type = ? and min<=? and max>=?  ALLOW FILTERING")
				err = s.Session.Query(Query_for_upi_api, servicetype.Upi_api, amount.Amountindouble, amount.Amountindouble).Consistency(gocql.One).Scan(&service.Detail, &service.Servicename)
				err1 = s.Session.Query(Query_for_upi_adm, servicetype.Upi_adm, amount.Amountindouble, amount.Amountindouble).Consistency(gocql.One).Scan(&service.Detail_adm, &service.Servicename_adm)
				if err != nil || err1 != nil {
					fmt.Println(err.Error()) ////check for both upi api and upi adm
					resp.Status = "-1"
					resp.StatusDesc = "Service Provider Not Mapped "
					return resp
				} else {
					// Transaction initiation code here in upi_transaction_status table
					uniqueid := Genereteid()
					// To check the uniqueid already present in db or not
					//fmt.Println(uniqueid)
					//fmt.Println("The virtual address ", qrvalue.Virtualaddress)

					transaction := dbservice.Upitransaction{}
					transaction.Amount = amount.Amountindouble
					transaction.Status = "INITIATED"
					transaction.Transaction_status_code = "1"
					transaction.Operation_performed = "QR_COLLECT"
					transaction.Transaction_type = "UPI"
					transaction.Merchant_type = qrvalue.Merchanttype
					transaction.Credit_debit = true
					transaction.User_id = user.User_id
					transaction.Api_user_id = user.Api_user_id
					transaction.Wallet_id = user.Wallet2_id
					transaction.Api_wallet_id = user.Api_wallet2_id
					transaction.Service_provider_id = service.Detail_adm // upi adm id should be stored here
					transaction.Api_service_provider_id = service.Detail //upi api id should be stored here
					transaction.Origin_identifier = qrvalue.Virtualaddress
					transaction.Transaction_note = qrvalue.Message
					transaction.Created_date = time.Now()
					transaction.Updated_date = time.Now()
					transaction.Id = uniqueid

					if qrvalue.Merchanttype == "DIRECT" {
						transaction.Client_unique_id = "INDB000000384512"
						transaction.Param_a = "iserveupvtltd@indus"
						transaction.Param_b = "iServeU"
						transaction.Param_c = "6012"
					}
					if qrvalue.Merchanttype == "AGGREGATE" {
						transaction.Client_unique_id = "INDB000000349936"
						transaction.Param_a = "iserveubiz@indus"
						transaction.Param_b = "iServeU"
						transaction.Param_c = "6012"
					}
					if qrvalue.Issl == "1" {
						transaction.Is_sl = true
					} else {
						transaction.Is_sl = false
					}
					//fmt.Println("Insert start")
					//fmt.Println(user.Wallet2_id)
					resp, err = Upitransaction(transaction, s)
					if err != nil {
						resp.Status = "-1"
						resp.StatusDesc = "Exception in insert data"
						return resp
					} else {
						// Publishing the STATUS=INITIATED, TxnStatusCode=1 to PubSub

						// const projectID = "iserveustaging"
						// const topicID = "upi_topic_staging"
						// err = Publishstatusmessage(projectID, topicID, transaction, user)
						v, err := Geturlinit(transaction)
						if err != nil {
							fmt.Println(err.Error())
							resp.Status = "-1"
							resp.StatusDesc = "Error in generating url from the data"
							return resp
						} else {
							qrcode_data, err := Generateqrcode_data(v)
							if err != nil {
								fmt.Println(err.Error())
								resp.Status = "-1"
								resp.StatusDesc = "Error in generating QR Code Base64 from the url"
								return resp
							} else {
								//transid := utils.Upitransactionidcount{}
								fmt.Println("Base64 Value of Qr_Code is :  ", qrcode_data)
								resp.Status = transaction.Transaction_status_code
								resp.StatusDesc = "Qr code generation done"
								resp.TxnId = strconv.FormatInt(transaction.Id, 10)
								resp.MerchantId = transaction.Client_unique_id
								resp.Amount = strconv.FormatFloat(transaction.Amount, 'f', 1, 32)
								resp.Paymentstate = transaction.Status
								resp.PayeeVPA = transaction.Param_a
								//resp.PayerVPA = transid.Origin_idntf //if empty string in request comes it should show null
								resp.Qrdata = qrcode_data
								return resp
							}

						}

					}

				}
			}

		} else {
			resp.Status = "-1"
			resp.StatusDesc = "user Not availeble In Table"
			return resp
		}
	} else {
		resp.Status = "-1"
		resp.StatusDesc = "user Not availeble With This Names"
		return resp
	}
}
