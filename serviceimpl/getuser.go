package serviceimpl

import (
	"errors"
	"fmt"

	"main.go/dbservice"
)

func Getuser(username string, s *Server) (dbservice.Userapimap, error) {
	var userapimapping dbservice.Userapimap

	userapimapping, err = Getuserapimapping(username, s)
	if err != nil {
		return userapimapping, err

	}
	fmt.Println(userapimapping)

	return userapimapping, nil

}
func Getuserapimapping(userName string, s *Server) (dbservice.Userapimap, error) {

	var userapimapping []dbservice.Userapimap
	m := map[string]interface{}{}

	query := "select * from user_api_mapping where user_name = ?"
	plandata := s.Session.Query(query, userName).Iter()
	for plandata.MapScan(m) {

		userapimapping = append(userapimapping, dbservice.Userapimap{
			Api_user_id:    m["api_user_id"].(int64),
			Api_wallet2_id: m["api_wallet2_id"].(int64),
			Api_wallet_id:  m["api_wallet_id"].(int64),
			City:           m["city"].(string),
			State:          m["state"].(string),
			User_id:        m["user_id"].(int64),
			User_name:      m["user_name"].(string),
			Wallet2_id:     m["wallet2_id"].(int64),
			Wallet_id:      m["wallet_id"].(int64),
		})
	}

	fmt.Println("userapimapping details", userapimapping)
	if len(userapimapping) == 0 {
		var emptyplan dbservice.Userapimap
		return emptyplan, errors.New("no Data Available for User apimapping")
	}

	return userapimapping[0], nil

}
