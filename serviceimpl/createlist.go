package serviceimpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"main.go/utils"
)

func Createlist(data utils.Create, s *Server) utils.Response2 {
	//var data1 utils.Create
	resp := utils.Response2{}
	apiresponse, resp := AggregateMerchantOnBoarding3(data)
	fmt.Println("API Response", apiresponse)
	if apiresponse.Status != "" {
		//decrypt msg
		fmt.Println("decrypt")
		if apiresponse.Status == "s" {
			fmt.Println("return decrypt")
			//returned decrypted response
		} else {
			resp.Status = "1"
			return resp
		}

	} else {
		resp.Status = "-1"
		resp.StatusDesc = "failed to fetch response"
		return resp
	}
	//out, err := json.Marshal(data)
	//fmt.Println(err)
	//encryptedRequest := encrypt(string(out), "f06e573cbbaeb2757bcc953fe7fd7933")

	return resp
}
func AggregateMerchantOnBoarding3(data utils.Create) (utils.SubMerchantonBoardApiResponse, utils.Response2) {
	apiresponse := utils.SubMerchantonBoardApiResponse{}
	resp := utils.Response2{}
	out, err := json.Marshal(data)
	fmt.Println(err)
	encryptedRequest, err := encrypt(string(out), "f06e573cbbaeb2757bcc953fe7fd7933")
	if err == nil {
		requestingJson := utils.BankRequest{}
		requestingJson.RequestMsg = encryptedRequest
		requestingJson.PgMerchantId = data.PgMerchantId
		baseUrl := " https://apig.indusind.com/ibl/prod/api/onBoardSubMerchant"
		jsonValue, _ := json.Marshal(data)

		u := bytes.NewReader(jsonValue)

		req, err := http.NewRequest("POST", baseUrl, u)
		if err != nil {
			fmt.Println("Error is req: ", err)
			log.Println("Error is req: ", err)
		}
		fmt.Println(req)
		req.Header.Set("Content-Type", "application/json")
		apiresponse.Status = "s"

	} else {
		resp.Status = "-1"
		resp.StatusDesc = "failed to encrypt requestS"
		return apiresponse, resp
	}
	return apiresponse, resp

}
