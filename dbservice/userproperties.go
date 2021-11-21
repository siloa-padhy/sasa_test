package dbservice

import "time"

type Userproperties struct {
	Admin             int64
	Is_upi_settle     string
	Master            int64
	Upi_settle_update time.Time
	User_id           int64
	User_name         string
}
