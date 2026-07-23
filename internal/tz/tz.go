// Package tz provides time zone information for the Asia/Dhaka time zone.
package tz

import (
	"time"
	_ "time/tzdata"
)

var Dhaka *time.Location

func init() {
	var err error
	Dhaka, err = time.LoadLocation("Asia/Dhaka")
	if err != nil {
		panic(err)
	}
}
