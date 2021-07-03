package main

import "C"
import (
	"goCron/cron"
	"time"
)

func main() {
	c := cron.New()
	c.AddJob("menma", time.Duration(5), time.Duration(5*time.Second), func() {
		time.Sleep(3 * time.Second)
	})
	c.Start()
	c.AddJob("menma", time.Duration(5), time.Duration(3*time.Second), func() {
		time.Sleep(3 * time.Second)
	})
	time.Sleep(50 * time.Second)
}
