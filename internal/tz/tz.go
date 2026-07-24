// Package tz provides time zone information for the Asia/Dhaka time zone.
package tz

import (
	_ "embed"
	"fmt"
	"time"
)

//go:embed zoneinfo/Asia/Dhaka
var dhakaTZ []byte

var Dhaka *time.Location

func init() {
	var err error
	Dhaka, err = time.LoadLocationFromTZData("Asia/Dhaka", dhakaTZ)
	if err != nil {
		panic(fmt.Errorf("loading Asia/Dhaka time zone: %w", err))
	}
}
