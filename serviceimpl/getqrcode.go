package serviceimpl

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"

	"github.com/skip2/go-qrcode"
	//"github.com/makiuchi-d/gozxing"
	//"github.com/makiuchi-d/gozxing/oned"
	"main.go/utils"
)

func GetQrdata(qr utils.Getfrmuser) (string, error) {
	u, err := url.Parse("upi://pay?")
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()

	q.Add("pa", qr.Pa)
	q.Add("pn", qr.Pn)
	u.RawQuery = q.Encode()

	v, err := url.PathUnescape(u.String())
	fmt.Println(v)
	if err != nil {
		log.Fatal(err)
	}
	var png []byte

	png, err = qrcode.Encode("v", qrcode.Medium, 64)

	qrcode_data := base64.StdEncoding.EncodeToString([]byte(png))
	fmt.Println(qrcode_data)

	return qrcode_data, err
}
