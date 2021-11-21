package serviceimpl

import (
	"encoding/base64"
	"fmt"

	"github.com/skip2/go-qrcode"
)

func Generateqrcode_data(qrdata string) (string, error) {
	var png []byte
	png, err = qrcode.Encode("v", qrcode.Medium, 256)
	fmt.Println()
	qrcode_data := base64.StdEncoding.EncodeToString([]byte(png))
	//fmt.Println(qrcode_data)
	return qrcode_data, err
}
