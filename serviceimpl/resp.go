package serviceimpl

import "main.go/utils"

func My_response(resp utils.Bank_response, bank utils.Bankdata) utils.Bank_response {
	if bank.ResponseCode == "0" || bank.ResponseCode == "00" {
		resp.ApiComment = bank.Status_desc
		resp.Status = "0" //string(0)
		resp.StatusDesc = "SUCCESS"
		resp.CustRefNo = bank.CustRefNo
		resp.OrderNo = bank.OrderNo
	} else {
		resp.ApiComment = bank.Status_desc
		resp.Status = "3" //int
		resp.StatusDesc = "FAILED"
		resp.CustRefNo = bank.CustRefNo
		resp.OrderNo = bank.OrderNo
	}
	return resp
}
