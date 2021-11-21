package serviceimpl

import (
	"fmt"
	"net/url"
	"strconv"

	"main.go/dbservice"
)

func Geturlinit(transaction dbservice.Upitransaction) (string, error) {
	a := strconv.FormatInt(int64(transaction.Id), 10)
	b := strconv.FormatInt(int64(transaction.Amount), 10)
	str := "upi://pay"
	u, _ := url.Parse(str)
	//fmt.Println("url:", u)
	values, _ := url.ParseQuery(u.RawQuery)
	values.Add("pa", transaction.Param_a)
	values.Add("pn", transaction.Param_b)
	values.Add("mc", transaction.Param_c)
	values.Add("tr", a)
	values.Add("tn", transaction.Transaction_note)
	values.Add("am", b)
	u.RawQuery = values.Encode()
	//fmt.Println("new url:", u)
	dataurl, _ := url.PathUnescape(u.String())
	fmt.Println("URL is :  ", dataurl)
	return dataurl, err
}

//pa=iserveupvtltd@indus&pn=iserveu&mc=6012&tr=34567893456789&mam=10&tn=WalletTopUp&am=100.0
//"upi://pay?pa=iserveupvtltd@indus&pn=iserveu&mc=6012&tr=34567893456789&mam=10&tn=WalletTopUp&am=100.0"
