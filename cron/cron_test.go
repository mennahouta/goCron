package cron

import (
	"testing"
	"time"
)

func TestAddJobBeforeStart(t *testing.T) {
	done := make(chan bool)
	cron := New()
	cron.AddJob("job1", time.Duration(0), time.Duration(1*time.Second), func() {
		done <- true
	})
	cron.Start()
	defer cron.Stop()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("job1 did not run")
	}
}

func TestAddJobAfterStart(t *testing.T) {
	done := make(chan bool)
	cron := New()
	cron.Start()
	defer cron.Stop()
	cron.AddJob("job2", time.Duration(0), time.Duration(1*time.Second), func() {
		done <- true
	})

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("job2 did not run")
	}
}

func TestNumbersOfCalls_AddJobBeforeStart(t *testing.T) {
	numOfCalls := 0
	cron := New()
	cron.AddJob("job3", time.Duration(0), time.Duration(1*time.Second), func() {
		numOfCalls += 1
	})
	cron.Start()
	defer cron.Stop()

	time.Sleep(5 * time.Second)
	if numOfCalls != 4 {
		t.Errorf("want number of calls for job3: 4, got: %d", numOfCalls)
	}
}

func TestNumberOfCalls_AddJobAfterStart(t *testing.T) {
	numOfCalls := 0
	cron := New()
	cron.Start()
	defer cron.Stop()
	cron.AddJob("job4", time.Duration(0), time.Duration(5*time.Second), func() {
		numOfCalls += 1
	})

	time.Sleep(11 * time.Second)
	if numOfCalls != 2 {
		t.Errorf("want number of calls for job4: 2, got: %d", numOfCalls)
	}
}

func TestRunTimeDifference(t *testing.T) {
	numOfCalls := 0
	cron := New()
	cron.Start()
	defer cron.Stop()
	cron.AddJob("job5", time.Duration(0), time.Duration(3*time.Second), func() {
		numOfCalls += 1
	})

	time.Sleep(4 * time.Second)
	if diff := cron.jobs[0].NextRunTime.Sub(cron.jobs[0].PrevRunTime); diff != 3*time.Second {
		t.Errorf("want time difference between job5 runs: 3s, got: %s", diff.String())
	}
}
