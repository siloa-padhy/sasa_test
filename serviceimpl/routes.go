package serviceimpl

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gocql/gocql"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"main.go/utils"
)

type Server struct {
	cluster *gocql.ClusterConfig
	Session *gocql.Session
}

const panpattern = "[A-Z]{5}[0-9]{4}[A-Z]{1}"
const mobilepattern = "(0|91)?[1-9][0-9]{9}"

const initamount = "[0-9.]+"

func NewRouter(conn *Server) *mux.Router {
	r := mux.NewRouter()
	fmt.Println("router")
	r.HandleFunc("/", checkupi).Methods("GET")
	r.HandleFunc("/upi/submerchant/selfonboarding", conn.selfonboard).Methods("POST")
	r.HandleFunc("/upi/submerchant/onboarding", conn.bulkonboard).Methods("POST")
	r.HandleFunc("/upi/init", conn.upiinit).Methods("POST")
	r.HandleFunc("/upi/qr/init", conn.generateqr).Methods("POST") //tested
	r.HandleFunc("/upi/submerchant/generateqrcode", conn.generateqrcode).Methods("POST")
	r.HandleFunc("/upi/list", conn.createlist).Methods("POST")
	r.HandleFunc("/upi/submerchant/generateqrcodeUrl", conn.getqrfromuser).Methods("POST")
	r.HandleFunc("/upi/submerchant/checkMobileNumber", conn.checkmobilenumber).Methods("POST") //tested
	r.HandleFunc("/upi/validate/vpa", conn.validatevpa).Methods("POST")                        //unit testing
	r.HandleFunc("/upi/submerchant/checkvpa", conn.checkvpa).Methods("POST")                   //tested
	r.HandleFunc("/upi/txnStatus", conn.statusenquiry).Methods("POST")
	r.HandleFunc("/indus/direct/commitResponse", conn.dynamical).Methods("POST")
	r.HandleFunc("/indus/commitResponse", conn.staticcal).Methods("POST")
	return r
}

func checkupi(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Commission appliaction start(GOLANG)")
}
func (s *Server) staticcal(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	reqBody, _ := ioutil.ReadAll(req.Body)
	var req6 utils.Static
	resp := utils.Bank_response{}
	json.Unmarshal(reqBody, &req6)
	resp = Static_call(req6, s)
	json.NewEncoder(res).Encode(resp)

	//fmt.Println(n)
	//fmt.Println(resp)
}

func (s *Server) checkvpa(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	v := validator.New()
	reqBody, _ := ioutil.ReadAll(req.Body)
	var vpa utils.Vparequest
	json.Unmarshal(reqBody, &vpa)
	resp := utils.Checkvpa_response{}

	if vpa.VAReqType == "R" || vpa.VAReqType == "T" {
		valid := utils.Vparequest{
			VirtualAddress: vpa.VirtualAddress,
			VAReqType:      vpa.VAReqType,
		}
		err := v.Struct(valid)
		if err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				fmt.Println("validation error", e)
				resp.Status = "-1"
				resp.StatusDesc = err.Error()
				json.NewEncoder(res).Encode(resp)
				return
			}
		} else {
			vpa.Authorization = req.Header.Get("Authorization")
			resp := Checkvpa_upi(vpa, s)
			json.NewEncoder(res).Encode(resp)
		}
	} else {
		resp.Status = "2"
		resp.StatusDesc = "Vpa Check Failed,Vpa Request Type must be R (Registration) 'or' T (Transaction) "
		json.NewEncoder(res).Encode(resp)
	}

}
func (s *Server) validatevpa(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	v := validator.New()
	reqBody, _ := ioutil.ReadAll(req.Body)
	var vpa utils.Vparequest
	json.Unmarshal(reqBody, &vpa)
	resp := utils.Vparesponse{}
	if vpa.VAReqType == "R" || vpa.VAReqType == "T" {
		valid := utils.Vparequest{
			VirtualAddress: vpa.VirtualAddress,
			VAReqType:      vpa.VAReqType,
		}
		err := v.Struct(valid)
		if err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				fmt.Println("validation error", e)
				log.Println(err.Error())
				resp.Status = "-1"
				resp.StatusDesc = err.Error()
				json.NewEncoder(res).Encode(resp)
				return
			}
		} else {
			resp := Check_vpa_in_bank(vpa, s)
			json.NewEncoder(res).Encode(resp)
		}
	} else {
		resp.Status = "2"
		resp.StatusDesc = "Vpa Check Failed,Vpa Request Type must be R (Registration) 'or' T (Transaction) "
		json.NewEncoder(res).Encode(resp)
	}
}

func (s *Server) generateqr(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	v := validator.New()
	reqBody, _ := ioutil.ReadAll(req.Body)
	var qrvalue utils.Initreq
	json.Unmarshal(reqBody, &qrvalue)
	resp := utils.Qr_response{}
	qrvalue.Authorization = req.Header.Get("Authorization")
	if qrvalue.Merchanttype == "DIRECT" || qrvalue.Merchanttype == "AGGREGATE" {
		if qrvalue.Issl == "1" || qrvalue.Issl == "0" {
			qrvalid := utils.Initreq{
				Amount:        qrvalue.Amount,
				Merchanttype:  qrvalue.Merchanttype,
				Issl:          qrvalue.Issl,
				Authorization: qrvalue.Authorization,
			}
			err := v.Struct(qrvalid)
			if err != nil {
				for _, e := range err.(validator.ValidationErrors) {
					fmt.Println("validation error in generating QR Code", e)
					log.Println(err.Error())
					resp.Status = "-1"
					resp.StatusDesc = err.Error()
					json.NewEncoder(res).Encode(resp)
					return
				}
			} else {
				resp = Generateqrcode(qrvalue, s)
				json.NewEncoder(res).Encode(resp)
			}
		} else {
			resp.Status = "-1"
			resp.StatusDesc = "Please provide issl as 0 or 1"
			json.NewEncoder(res).Encode(resp)
		}
	} else {
		resp.Status = "-1"
		resp.StatusDesc = "Please provide merchantType as DIRECT or AGGREGATE "
		json.NewEncoder(res).Encode(resp)
	}
}
func (s *Server) selfonboard(res http.ResponseWriter, req *http.Request) {
	//log.Print("This is my log msg")
	fmt.Println("Self Onboarding")
	res.Header().Set("Content-Type", "application/json")
	v := validator.New()
	reqBody, _ := ioutil.ReadAll(req.Body)
	var value utils.Selfonboard
	json.Unmarshal(reqBody, &value)
	resp := utils.Response{}
	isStringAlphabet := regexp.MustCompile(panpattern).MatchString
	if !isStringAlphabet(value.PanNo) {
		resp.Status = -1
		resp.StatusDesc = "onBoarding Failed,Pan Card Not valid"
		json.NewEncoder(res).Encode(resp)
		return
	}
	value.Authorization = req.Header.Get("Authorization")
	valid := utils.Selfonboard{
		StrCntMobile:  value.StrCntMobile,
		FirstName:     value.FirstName,
		LastName:      value.LastName,
		AccNo:         value.AccNo,
		Ifsc:          value.Ifsc,
		PanNo:         value.PanNo,
		MerVirtualAdd: value.MerVirtualAdd,
		LegalStrName:  value.LegalStrName,
		BankName:      value.BankName,
		Mebussname:    value.Mebussname,
		Authorization: value.Authorization,
	}
	err := v.Struct(valid)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Print("validation error", e)
			log.Println(err.Error())
			resp.Status = -1
			resp.StatusDesc = err.Error()
			json.NewEncoder(res).Encode(resp)
			return
		}
	} else {
		resp = Selfonboard(value, s)
		json.NewEncoder(res).Encode(resp)
		return
	}
}

func (s *Server) bulkonboard(res http.ResponseWriter, req *http.Request) {
	var bulkreqdata utils.Bulkreq
	v := validator.New()
	err := json.NewDecoder(req.Body).Decode(&bulkreqdata)
	res.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("Error in parsing", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	resp := utils.Response{}
	isStringAlphabet := regexp.MustCompile(panpattern).MatchString
	if !isStringAlphabet(bulkreqdata.PanNo) {
		resp.Status = -1
		resp.StatusDesc = "onBoarding Failed,Pan Card Not valid"
		json.NewEncoder(res).Encode(resp)
		return
	}
	bulkreqdata.Authorization = req.Header.Get("Authorization")
	a := utils.Bulkreq{
		AccNo:         bulkreqdata.AccNo,
		LegalStrName:  bulkreqdata.LegalStrName,
		Mebussname:    bulkreqdata.Mebussname,
		MerVirtualAdd: bulkreqdata.MerVirtualAdd,
		PanNo:         bulkreqdata.PanNo,
		StrCntMobile:  bulkreqdata.StrCntMobile,
		FirstName:     bulkreqdata.FirstName,
		LastName:      bulkreqdata.LastName,
		Ifsc:          bulkreqdata.Ifsc,
		UserName:      bulkreqdata.UserName,
		BankName:      bulkreqdata.BankName,
		Authorization: bulkreqdata.Authorization,
	}
	err = v.Struct(a)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println("validation error", e)
			resp.Status = -1
			resp.StatusDesc = err.Error()
			json.NewEncoder(res).Encode(resp)
			return
		}
	} else {
		resp := Checkvirtualaddress(bulkreqdata, s)
		json.NewEncoder(res).Encode(resp)
		return
	}
}

func (s *Server) upiinit(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var initreqdata utils.Initreq

	v := validator.New()
	err := json.NewDecoder(req.Body).Decode(&initreqdata)

	if err != nil {
		fmt.Println("Error in parsing", err)
		res.WriteHeader(http.StatusBadRequest)
	}

	resp := utils.Qr_response{}

	isStringAlphabet := regexp.MustCompile(initamount).MatchString
	if !isStringAlphabet(initreqdata.Amount) {
		resp.Status = "-1"
		resp.StatusDesc = "Invalid input amount"
		json.NewEncoder(res).Encode(resp)
		return
	}
	a := utils.Initreq{
		Amount:       initreqdata.Amount,
		Merchanttype: initreqdata.Merchanttype,
		Issl:         initreqdata.Issl,
	}
	err = v.Struct(a)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println("validation error", e)
			resp.Status = "-1"
			resp.StatusDesc = err.Error()
			json.NewEncoder(res).Encode(resp)
			fmt.Println(resp)
			return
		}

	} else {
		initreqdata.Authorization = req.Header.Get("Authorization")
		resp := Inititedfunc(initreqdata, s)
		json.NewEncoder(res).Encode(resp)
		return

	}

}

func (s *Server) createlist(res http.ResponseWriter, req *http.Request) {

	reqBody, _ := ioutil.ReadAll(req.Body)
	var req1 utils.Create
	resp := utils.Response2{}
	err := json.Unmarshal(reqBody, &req1)
	if err != nil {
		fmt.Println("Error")
		log.Println("error")
	}
	resp = Createlist(req1, s)
	json.NewEncoder(res).Encode(resp)
	fmt.Println(resp)
	return

}
func (s *Server) getqrfromuser(res http.ResponseWriter, req *http.Request) {

	v := validator.New()
	reqBody, _ := ioutil.ReadAll(req.Body)
	var req3 utils.Getfrmuser
	json.Unmarshal(reqBody, &req3)
	resp := utils.Response1{}
	valid := utils.Getfrmuser{
		UserName: req3.UserName,
	}
	err := v.Struct(valid)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println("name can not be empty", e)
			resp.Status = "-1"
			resp.Message = "Invalied Request " + err.Error()
			json.NewEncoder(res).Encode(resp)
			return
		}
	} else {
		resp = Generateqr(req3, s)
		json.NewEncoder(res).Encode(resp)
		return
	}

}
func (s *Server) generateqrcode(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	reqBody, _ := ioutil.ReadAll(req.Body)
	var req2 utils.Generate
	var data utils.Getfrmuser
	json.Unmarshal(reqBody, &req2)
	resp := utils.Response1{}
	req2.Authorization = req.Header.Get("Authorization")
	username := utils.Getusername(req2.Authorization)
	fmt.Println("Username is :", username)
	data.UserName = username
	data.Pa = req2.Pa
	data.Pn = req2.Pn
	resp = Generateqr(data, s)
	json.NewEncoder(res).Encode(resp)

}

func (s *Server) checkmobilenumber(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	v := validator.New()
	reqBody, _ := ioutil.ReadAll(req.Body)
	var mobile utils.Check_mobile_number
	json.Unmarshal(reqBody, &mobile)
	resp := utils.Mobile_Response{}
	isStringAlphabet := regexp.MustCompile(mobilepattern).MatchString
	if !isStringAlphabet(mobile.StrCntMobile) {
		resp.Status = "-1"
		resp.StatusDesc = "Invalied Mobile Number , Mobile pattern not match "
		json.NewEncoder(res).Encode(resp)
		return
	}
	valid := utils.Check_mobile_number{
		StrCntMobile: mobile.StrCntMobile,
	}
	err := v.Struct(valid)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println("validation error", e)
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Println(err.Error())
			resp.Status = "-1"
			resp.StatusDesc = "Incorrect Mobile Number" + err.Error()
			json.NewEncoder(res).Encode(resp)
			return
		}
	} else {
		resp := Findby_mobilenumber(mobile, s)
		json.NewEncoder(res).Encode(resp)
	}
}
func (s *Server) statusenquiry(res http.ResponseWriter, req *http.Request) {

	reqBody, _ := ioutil.ReadAll(req.Body)
	var req6 utils.Statusen
	v := validator.New()
	resp := utils.Response2{}
	err := json.Unmarshal(reqBody, &req6)
	if err != nil {
		fmt.Println("Error")
		log.Println("error")
	}
	a := utils.Statusen{

		PgMerchantId: req6.PgMerchantId,
	}
	err = v.Struct(a)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println("validation error", e)
			resp.Status = "-1"
			resp.StatusDesc = err.Error()
			json.NewEncoder(res).Encode(resp)
			fmt.Println(resp)
			return
		}
	} else {
		resp = Statusenquiry(req6, s)
		json.NewEncoder(res).Encode(resp)
		fmt.Println(resp)
		return
	}

}
func (s *Server) dynamical(res http.ResponseWriter, req *http.Request) {

	reqBody, _ := ioutil.ReadAll(req.Body)
	var req6 utils.Dynamic
	resp := utils.Response2{}
	json.Unmarshal(reqBody, &req6)
	resp = Dynamicall(req6, s)
	json.NewEncoder(res).Encode(resp)
	fmt.Println(resp)
	return

}
