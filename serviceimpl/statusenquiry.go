package serviceimpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"main.go/utils"
)

func Statusenquiry(data utils.Statusen, s *Server) utils.Response2 {

	resp := utils.Response2{}

	if data.MerchantType == "AGGREGATE" {

		req := utils.Statusen{}
		req.TxId = data.TxId
		req.PgMerchantId = "INDB000000003150"
		req.CustRefNo = strings.ToLower(data.CustRefNo)
		req.PspRefNo = data.PspRefNo
		req.NpciTranId = "INDBAA4F5A16C75A4C5F8988CBEA8022FDB"
		out, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		//apirespoanse, resp := Bankcall(req)
		encryptedRequest, err := encrypt(string(out), "f06e573cbbaeb2757bcc953fe7fd7933")
		if err == nil {
			//hit bank
			fmt.Println("call to bank")
			requestingJson := utils.BankRequest{}
			requestingJson.RequestMsg = encryptedRequest
			requestingJson.PgMerchantId = data.PgMerchantId
			apirespoanse, resp := Bankcall(data, requestingJson)
			if apirespoanse.Status != " " {
				//decrypt the response
				fmt.Println("decrypt")
				//convert the bank response txn inquiery pojo
				//Trasactionpojo

			} else {
				resp.Status = "3"
				resp.StatusDesc = "Getting no response While checking for Transaction status,Please try again"
				return resp

			}
		} else {
			resp.Status = "2"
			resp.StatusDesc = "Failed to encrypt Request"
			return resp
		}

	} else if data.MerchantType == "DIRECT" {

		out, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		encryptedRequest, err := encrypt(string(out), "b27e4ddb1fcb9e3ae394b13a25155ac1")

		if err == nil {
			//hit bank
			fmt.Println("call to bank")
			requestingJson := utils.BankRequest{}
			requestingJson.RequestMsg = encryptedRequest
			requestingJson.PgMerchantId = data.PgMerchantId
			apirespoanse, resp := Bankcall(data, requestingJson)
			if apirespoanse.Status != " " {
				fmt.Println("decrypt")
				//decrypt the response
				//convert the bank response txn inquiery pojo

			} else {
				resp.Status = "3"
				resp.StatusDesc = "Getting no response While checking for Transaction status,Please try again"
				return resp

			}
		} else {
			resp.Status = "2"
			resp.StatusDesc = "Failed to encrypt Request"
			return resp
		}

	} else {
		resp.Status = "-1"
		resp.StatusDesc = "Please provide a valid merchanrtype for upi status for status check"
		log.Printf("Please provide a valid merchanrtype for upi status for status check,")
		return resp

	}
	return resp
}
func Bankcall(data utils.Statusen, data2 utils.BankRequest) (utils.SubMerchantonBoardApiResponse, utils.Response2) {

	baseUrl := " https://apig.indusind.com/ibl/prod/api/onBoardSubMerchant"
	jsonValue, _ := json.Marshal(data)
	resp := utils.Response2{}
	apiresponse := utils.SubMerchantonBoardApiResponse{}
	u := bytes.NewReader(jsonValue)

	req, err := http.NewRequest("POST", baseUrl, u)
	if err == nil {
		fmt.Println("Error is req: ", err)

		req.Header.Set("Content-Type", "application/json")
		apiresponse.Status = "ok"
	} else {
		resp.Status = "3"
		resp.StatusDesc = "Exception while cheacking trasaction"
		log.Panicln("Exception while cheacking trasaction:", err)
		return apiresponse, resp
	}
	return apiresponse, resp
}
