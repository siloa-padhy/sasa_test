package serviceimpl

import (
	"fmt"
	"strconv"

	"main.go/dbservice"
	"main.go/utils"
)

func Dynamicall(data utils.Dynamic, s *Server) utils.Response2 {

	//var res string
	//decrypt := decrypt(data.Resp, "b27e4ddb1fcb9e3ae394b13a25155ac1")
	//fmt.Println("decrypt", decrypt)
	//fmt.Println("decrypt", decrypt)
	// a, _ := json.Marshal(decrypt)
	// fmt.Println(a)
	var rrn int64
	var id int64
	var resp2 utils.Response2
	var resp utils.Dynamicres
	resp.CustRefNo = "110569213605"
	resp.OrderNo = "741449330982928"
	resp.Status = "SUCCESS"
	err := s.Session.Query(`select count(*) from upi_transaction_status where rrn = ? allow filtering`, resp.CustRefNo).Scan(&rrn)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("error in db")
	}
	if rrn == 0 {

		err := s.Session.Query(`select count(*) from upi_transaction_status where id = ? allow filtering`, resp.OrderNo).Scan(&id)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error in db")
		}
		if id != 0 {
			if resp.Status == "SUCCESS" {
				//update
				response := dbservice.Upitransaction{}
				response.Status = "SUCCESS"
				response.Transaction_status_code = "5"
				intid, _ := strconv.ParseInt(resp.OrderNo, 10, 64)
				response.Id = intid

				err = Updateupiproperties(response, s)
				if err == nil {
					fmt.Println("updated")
				}
				//publish in pubsub

			} else if resp.Status == "FAILED" {
				response := dbservice.Upitransaction{}
				response.Status = "FAILED"
				response.Transaction_status_code = "6"
				intid, _ := strconv.ParseInt(resp.OrderNo, 10, 64)
				response.Id = intid

				err = Updateupiproperties(response, s)
				if err == nil {
					fmt.Println("updated")
				}
				//publish in pubsub

			} else {
				response := dbservice.Upitransaction{}
				response.Status = "INPROGRESS"
				response.Transaction_status_code = "9"
				intid, _ := strconv.ParseInt(resp.OrderNo, 10, 64)
				response.Id = intid

				err = Updateupiproperties(response, s)
				if err == nil {
					fmt.Println("updated")
				}

				//publish in pubsub

			}
		} else {
			fmt.Println("order number not found")
		}

	} else {
		fmt.Println("RRN already present it should be unique")
	}
	return resp2
}
