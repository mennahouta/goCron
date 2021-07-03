# goCron

A cron scheduler written in go.

The user can:
- Start/stop the cron.
- Add a job to be executed periodically.
- Remove a job.
- Monitor actual run time duration of every executed job.

Usage snippets can be found in main.go

## Components
### Job
This is the job that needs to be scheduled and periodically run.

### Cron
This is the scheduler itself. It's responsible for receiving scheduling requests, running the jobs according to the schedule, and instrumenting each job run.

## Logic
- To implement the scheduling functionality, I thought of creating a goroutine for each job, the goroutine could hold a ticker channel and run the job on every tick, but then I thought that it doesn't make sense to have a background routine running for every job even when the job isn't running, it might not be a big deal for frequent jobs but imagine a job that has to be executed every 3 days or something? Not the best way to go.
- The cron should continuously check for jobs that need to be executed. Simplest way to do this is to sort the jobs by "next run time", look through them and send the ones that need to be executed for execution. That's the current implementation.

## Trade-offs
- Saving cron.jobs as map[string]*Job VS []*Job:
  - Since we have an id for each job, it's intuitive to use a map, but that didn't work well when it was time to sort the jobs according to their "next run time". It's doable even if we use a map[string]*Job, but it wouldn't be as straightforward. 
  - I thought of keeping the map as is, to use it for id uniqueness validation and job removal, and I could add a []*Job that I can use in sorting. But this redundancy of information didn't seem right to me, so I switched to only using a []*Job.
  - It's not a big tradeoff since the sorting will be done much more frequently than the id uniqueness validation and/or the job removal.
- Run Interval Duration validation:
  - I wasn't sure whether the cron should enforce the non-concurrency of instances of one job, but I decided not to enforce it, since the cron's job is to simply "schedule something to be executed every x interval"; I don't think it's supposed to be smart enough to make sure no per-job conflicts arise. That can be left to the user. 

## Future Work
- Run Interval Duration
  - The cron doesn't have to be smart but it can at least ask the user whether they're okay with the concurrency of instances of one job. 
  - The cron can warn the user if the actual run interval surpasses the expected run interval.
- Support dates and different timezones, not just durations.
- File logging.
- Accept requests from the user to see the current job-configurations.
- Instead of being completely in-memory, the cron can read/write its state from a file, so configured jobs won't go away in each run. 
  - This will also give the user the flexibility of writing a config file without having to use cron.AddJob for each job.
- Allow the user to edit an existing job.







