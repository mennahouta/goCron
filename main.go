package main

import "C"
import (
	"goCron/cron"
	"time"
)

func main() {
	c := cron.New()
	c.AddJob("1", time.Duration(5), time.Duration(5*time.Second), func() {
		time.Sleep(3 * time.Second)
	})
	c.Start()
	c.AddJob("2", time.Duration(5), time.Duration(3*time.Second), func() {
		time.Sleep(6 * time.Second)
	})
	c.AddJob("3", time.Duration(5), time.Duration(7*time.Second), func() {
		time.Sleep(13 * time.Second)
	})
	time.Sleep(30 * time.Second)
	c.Stop()
}
