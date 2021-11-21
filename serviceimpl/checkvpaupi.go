package serviceimpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"main.go/utils"
)

func Checkvpa_upi(value utils.Vparequest, s *Server) utils.Checkvpa_response {
	resp := utils.Checkvpa_response{}
	TxnId := Genereteid()
	//fmt.Println(TxnId)
	//fmt.Println(value.Authorization)
	username := utils.Getusername(value.Authorization)
	//fmt.Println(username)
	if TxnId == 0 || username == "" {
		resp.Status = "-1"
		resp.StatusDesc = "TxnId or Username unable to generate"
		return resp
	} else {
		count := utils.Vpatype{}
		req := utils.Vparequest{}
		req.PspRefNo = fmt.Sprint(TxnId)
		req.PgMerchantId = "INDB000000384512"
		req.VirtualAddress = strings.ToLower(value.VirtualAddress)
		req.VAReqType = value.VAReqType
		apiResponsedata := Bank_request(req)
		//fmt.Println(apiResponsedata)
		if apiResponsedata.Status != "" {
			if apiResponsedata.Status == "VE" {
				resp.Status = "-1"
				resp.StatusDesc = "vpa already taken"
				return resp
			}
			if apiResponsedata.Status == "VN" {
				//fmt.Println(value.VirtualAddress)
				err = s.Session.Query(`select count(*) from user_api_mapping where user_name = ? allow filtering`, username).Scan(&count.User_count)
				//fmt.Println("User count is :", count.User_count)
				if err != nil {
					log.Println(err.Error())
					resp.Status = "-1"
					resp.StatusDesc = "Error in Getting username count from DB"
					return resp
				}
				if count.User_count == 1 {
					err = s.Session.Query(`select count(*) from indus_upi_properties where payee_vpa = ? allow filtering`, value.VirtualAddress).Scan(&count.Vpa_count)
					//fmt.Println("Vpa count  is :", count.Vpa_count)
					if err != nil {
						log.Println(err.Error())
						resp.Status = "-1"
						resp.StatusDesc = "Error in Getting vpa count from indus  DB"
						return resp
					}
					if count.Vpa_count == 0 {
						resp.Status = "0"
						resp.StatusDesc = "VPA available"
						return resp
					} else {
						resp.Status = "-1"
						resp.StatusDesc = "VPA already taken"
						return resp
					}
				} else {
					resp.Status = "-1"
					resp.StatusDesc = "User not exist or some error . Kindly Try again or contact to admin"
					return resp
				}
			}
			if !(apiResponsedata.Status == "VE" || apiResponsedata.Status == "VN") {
				resp.Status = "-1"
				resp.StatusDesc = "Some issue in fetching vpa, kindly contact admin"
				return resp
			}
		} else {
			log.Println(err.Error())
			resp.Status = "-1"
			resp.StatusDesc = "Error in fetching response from bank"
			return resp
		}
	}
	return resp

}
func Bank_request(req utils.Vparequest) utils.Checkvpa_response {
	resp := utils.Checkvpa_response{}
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
	fmt.Println("String data", string(out))
	encryptedRequest, err := encrypt(string(out), "b27e4ddb1fcb9e3ae394b13a25155ac1")
	if err != nil {
		fmt.Println(err)
		resp.Status = "1"
		resp.StatusDesc = err.Error()
	}

	if encryptedRequest != "" {
		requestingJson := utils.BankRequest{}
		requestingJson.PgMerchantId = req.PgMerchantId
		//	requestingJson.RequestMsg = encryptedRequest
		requestingJson.RequestMsg = "8608325D38FD938622007BF5615DEF6C8D9D94349FB7090066836DBCD59B4F25AA24B3F2F994743F16EFC01F78E23822ED053EED6C30244CCBD3ABA8AD3E1DF40A9FE87FDFC57AD01D387FACF0C6226763A401E1604B23BE2F93D8EFCB395295240ED11A62F73B937EC8BE750FA47C2C977E54D648C8FD9DA2BB17309E488938F6A9D22EC24CA4997D096CD3AEA9623C2C51DEFE2B36EBA30553E8783B7571F7"
		baseUrl := "https://apig.indusind.com/ibl/prod/upijson/validateVPAWeb"
		fmt.Println("request data PG merchnat", requestingJson.PgMerchantId)
		fmt.Println("request data Encrypt", requestingJson.RequestMsg)

		jsonValue, _ := json.Marshal(requestingJson)
		u := bytes.NewReader(jsonValue)

		req, err := http.NewRequest("POST", baseUrl, u)
		if err != nil {
			fmt.Println("Error is req: ", err)
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}

		// Do sends an HTTP request and
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error in send req: ", err.Error())

		}

		defer resp.Body.Close()
		fmt.Println("Response is", resp)
		body, err := ioutil.ReadAll(resp.Body)
		decryptResponse, err := decrypt(body, "b27e4ddb1fcb9e3ae394b13a25155ac1")
		// fmt.Println(string(body))
		fmt.Println(string(decryptResponse))
		// //Close
		resp.Status = "VN"
	}
	return resp
}
