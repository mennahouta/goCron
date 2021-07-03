// Package cron is responsible for jobs scheduling and execution.
package cron

import (
	"goCron/job"
	"log"
	"sort"
	"time"
)

type Cron struct {
	isRunning   bool
	jobs        []*job.Job
	newJobAdded chan *job.Job
	stop        chan bool
}

func New() *Cron {
	return &Cron{
		isRunning:   false,
		jobs:        []*job.Job{},
		newJobAdded: make(chan *job.Job),
		stop:        make(chan bool),
	}
}

func (cron *Cron) AddJob(id string, runInterval, frequency time.Duration, toRun func()) {
	if !cron.validId(id) {
		log.Printf("Job ID %s already exists, please retry with a new ID.", id)
		return
	}
	j := &job.Job{
		UniqueID:            id,
		ExpectedRunInterval: runInterval,
		SchedulingFrequency: frequency,
		ProcessToRun:        toRun,
	}
	if !cron.isRunning {
		cron.jobs = append(cron.jobs, j)
	} else {
		cron.newJobAdded <- j
	}
}

func (cron *Cron) RemoveJob(id string) {
	for idx, j := range cron.jobs {
		if j.UniqueID == id {
			cron.jobs = append(cron.jobs[:idx], cron.jobs[idx+1:]...)
		}
	}
}

func (cron *Cron) Start() {
	if cron.isRunning {
		return
	}
	cron.isRunning = true
	go cron.run()
}

func (cron *Cron) Stop() {
	if cron.isRunning {
		cron.stop <- true
	}
}

func (cron *Cron) run() {
	now := time.Now()
	cron.setRunTimes(now)
	for {
		sort.SliceStable(cron.jobs, func(i, j int) bool {
			return cron.jobs[i].NextRunTime.Before(cron.jobs[j].NextRunTime)
		})
		var cronTimer *time.Timer
		if len(cron.jobs) == 0 || cron.jobs[0].NextRunTime.IsZero() {
			cronTimer = time.NewTimer(30 * time.Minute)
		} else {
			cronTimer = time.NewTimer(cron.jobs[0].NextRunTime.Sub(now))
		}
		select {
		case now = <-cronTimer.C:
			now = time.Now()
			for _, j := range cron.jobs {
				if j.NextRunTime.After(now) || j.NextRunTime.IsZero() {
					break
				}
				go cron.executeJob(now, j)
				j.SetNextRunTime()
			}
		case j := <-cron.newJobAdded:
			cronTimer.Stop()
			now = time.Now()
			j.NextRunTime = now.Add(j.SchedulingFrequency)
			cron.jobs = append(cron.jobs, j)
		case <-cron.stop:
			cron.isRunning = false
			cronTimer.Stop()
			return
		}
	}
}

func (cron *Cron) executeJob(now time.Time, j *job.Job) {
	defer func() {
		log.Printf("Job '%s' execution time: %s", j.UniqueID, time.Since(now).String())
	}()
	j.ProcessToRun()
}

func (cron *Cron) setRunTimes(now time.Time) {
	for _, j := range cron.jobs {
		j.NextRunTime = now.Add(j.SchedulingFrequency)
	}
}

func (cron *Cron) validId(id string) bool {
	for _, j := range cron.jobs {
		if j.UniqueID == id {
			return false
		}
	}
	return true
}
