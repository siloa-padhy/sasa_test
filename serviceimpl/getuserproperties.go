package serviceimpl

import (
	"errors"
	"fmt"
	"time"

	"main.go/dbservice"
)

func Getuserproperties(userid int64, s *Server) (dbservice.Userproperties, error) {
	var userproperties dbservice.Userproperties
	userproperties, err = Getuserpropertiesdata(userid, s)
	if err != nil {
		return userproperties, err
	}
	fmt.Println(userproperties)
	return userproperties, nil

}
func Getuserpropertiesdata(userid int64, s *Server) (dbservice.Userproperties, error) {

	var userproperties []dbservice.Userproperties
	m := map[string]interface{}{}
	query := "select * from user_properties where user_id= ?"
	plandata := s.Session.Query(query, userid).Iter()
	for plandata.MapScan(m) {
		userproperties = append(userproperties, dbservice.Userproperties{
			Admin:             m["admin"].(int64),
			Is_upi_settle:     m["is_upi_settle"].(string),
			Master:            m["master"].(int64),
			Upi_settle_update: m["upi_settle_update"].(time.Time),
			User_id:           m["user_id"].(int64),
			User_name:         m["user_name"].(string),
		})
	}
	fmt.Println("userproperties", userproperties)
	if len(userproperties) == 0 {
		var emptyplan dbservice.Userproperties
		return emptyplan, errors.New("no Data Available for Userproperties")
	}
	return userproperties[0], nil
}
