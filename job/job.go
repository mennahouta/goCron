// Package job includes the job definition and all its related methods.
package job

import (
	"time"
)

type Job struct {
	UniqueID            string
	ExpectedRunInterval time.Duration
	LastActualInterval  time.Duration
	SchedulingFrequency time.Duration
	ProcessToRun        func()
	NextRunTime         time.Time
	PrevRunTime         time.Time
}

func (j *Job) SetNextRunTime() {
	j.PrevRunTime = j.NextRunTime
	j.NextRunTime = j.PrevRunTime.Add(j.SchedulingFrequency)
}
