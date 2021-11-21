package serviceimpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"main.go/dbservice"
	"main.go/utils"
)

func Inititedfunc(initreqdata utils.Initreq, s *Server) utils.Qr_response {
	resp := utils.Qr_response{}
	//user := dbservice.Userapimap{}
	fmt.Println("enter")
	resp, err := Initiated(initreqdata, s)
	if err != nil {
		fmt.Println(err.Error())
	}
	//resp1.User_id
	fmt.Println("data", resp)
	apiresponse := AggregateMerchantOnBoarding4(initreqdata)
	fmt.Println(apiresponse)
	if apiresponse.Status == "s" {
		response := dbservice.Upitransaction{}
		response.Status = "INITIATED"
		response.Transaction_status_code = "10"
		response.Id = resp.Userid
		fmt.Println(resp.Userid)

		err = Updateupiproperties(response, s)
		if err == nil {
			fmt.Println("updated")
		}
		//publish to pub/sub topic upi topic

	} else if apiresponse.Status == "2" {

		response := dbservice.Upitransaction{}
		response.Status = "INITIATED"
		response.Transaction_status_code = "2"
		response.Id = resp.Userid
		fmt.Println(resp.Userid)

		err = Updateupiproperties(response, s)
		if err == nil {
			fmt.Println("updated")
		}
		//publish to pub/sub topic upi topic
	} else if apiresponse.Status == "4" {

		response := dbservice.Upitransaction{}
		response.Status = "PENDING"
		response.Transaction_status_code = "4"
		response.Id = resp.Userid
		fmt.Println(resp.Userid)

		err = Updateupiproperties(response, s)
		if err == nil {
			fmt.Println("updated")
		}
		//publish to pub/sub topic upi topic

	} else {

		response := dbservice.Upitransaction{}
		response.Status = "3"
		response.Transaction_status_code = "4"
		response.Id = resp.Userid
		fmt.Println(resp.Userid)

		err = Updateupiproperties(response, s)
		if err == nil {
			fmt.Println("updated")
		}
		//publish to pub/sub topic upi topic
	}
	return resp
}
func AggregateMerchantOnBoarding4(initreqdata utils.Initreq) utils.SubMerchantonBoardApiResponse {
	apiresponse := utils.SubMerchantonBoardApiResponse{}

	out, err := json.Marshal(initreqdata)
	fmt.Println(err)
	encryptedRequest, err := encrypt(string(out), "f06e573cbbaeb2757bcc953fe7fd7933")
	if err == nil {
		requestingJson := utils.BankRequest{}
		requestingJson.RequestMsg = encryptedRequest
		requestingJson.PgMerchantId = initreqdata.PgMerchantId
		baseUrl := "https://indusupiuat.indusind.com:9043/upi/web/onBoardSubMerchant"
		jsonValue, _ := json.Marshal(initreqdata)

		u := bytes.NewReader(jsonValue)

		req, err := http.NewRequest("POST", baseUrl, u)
		if err != nil {
			fmt.Println("Error is req: ", err)
		}
		fmt.Println(req)
		req.Header.Set("Content-Type", "application/json")
		apiresponse.Status = "4"

	}
	return apiresponse

}
