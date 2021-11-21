package serviceimpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"main.go/dbservice"
	"main.go/utils"
)

func Check_vpa_in_bank(value utils.Vparequest, s *Server) utils.Vparesponse {
	resp := utils.Vparesponse{}
	TxnId := Genereteid()
	fmt.Println(TxnId)
	if TxnId == 0 {
		resp.Status = "1"
		resp.StatusDesc = "TxnId unable to generate"
		return resp
	} else {
		req := utils.Vparequest{}
		req.PspRefNo = fmt.Sprint(TxnId)
		req.PgMerchantId = "INDB000000384512"
		req.VirtualAddress = value.VirtualAddress
		req.VAReqType = value.VAReqType
		apiResponsedata := Upi_bank_request(req)
		fmt.Println(apiResponsedata)
		if apiResponsedata.Status != "" {

			//dectypt
			if apiResponsedata.Status == "0" {

				//return pojo objecta as it is
				//status 200
				return apiResponsedata

			} else {

				pojo := dbservice.Pojo2{}
				pojo.Status = "0"
				//Pojo object and return
				//bad request,how to send bad request
				return apiResponsedata

			}

		} else {
			resp.Status = "-1"
			resp.StatusDesc = "failed to fetch response"
			return resp
		}
		//return resp

	}

}
func Upi_bank_request(req utils.Vparequest) utils.Vparesponse {
	resp := utils.Vparesponse{}
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
	if err != nil {
		fmt.Println(err)
	}
	encryptedRequest, err := encrypt(string(out), "b27e4ddb1fcb9e3ae394b13a25155ac1")
	if err != nil {
		fmt.Println(err)
	}

	if encryptedRequest != "" {
		requestingJson := utils.BankRequest{}
		requestingJson.PgMerchantId = req.PgMerchantId
		requestingJson.RequestMsg = encryptedRequest
		baseUrl := "https://apig.indusind.com/ibl/prod/upijson/validateVPAWeb"
		jsonValue, _ := json.Marshal(requestingJson)

		u := bytes.NewReader(jsonValue)

		req, err := http.NewRequest("POST", baseUrl, u)
		if err != nil {
			fmt.Println("Error is req: ", err)
		}
		req.Header.Set("Content-Type", "application/json")
		//client := &http.Client{}

		// Do sends an HTTP request and
		//resp, err := client.Do(req)
		// if err != nil {
		// 	fmt.Println("error in send req: ", err.Error())

		// }

		// //defer resp.Body.Close()
		// fmt.Println(resp)
		// body, err := ioutil.ReadAll(resp.Body)
		// decryptResponse := decrypt("", "b27e4ddb1fcb9e3ae394b13a25155ac1")
		// fmt.Println(string(body))
		// fmt.Println(string(decryptResponse))
		// //Close

		resp.Status = "VN"
	}

	return resp
}
