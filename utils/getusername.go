package utils

import (
	b64 "encoding/base64"
	"encoding/json"
	"strings"
)

type Getuser struct {
	Sub      string `json:"sub"`
	Audience string `json:"audience"`
	Created  int64  `json:"web"`
	Exp      int64  `json:"exp"`
}

func Getusername(authorization string) string {
	var strArr []string
	strArr = strings.Split(authorization, ".")
	uDec, _ := b64.URLEncoding.DecodeString(strArr[1])
	var userdetails Getuser
	json.Unmarshal([]byte(string(uDec)), &userdetails)
	return userdetails.Sub
}
