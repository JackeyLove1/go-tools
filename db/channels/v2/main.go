package main

import (
    "math/rand"
    "time"

    "golang.org/x/sync/errgroup"
)

type Job struct {
    Id         int
    RandNumber int
}

func (j Job) Work() int {
    sum := 0
    for j.RandNumber > 0 {
        sum += j.RandNumber % 10
        j.RandNumber /= 10
    }
    return sum
}

type Result struct {
    job *Job
    Sum int
}

//TODO: use ErrGroup to control it
func CreateConsumers(jobs chan *Job, results chan *Result, nums int) {
    g := new(errgroup.Group)
    g.SetLimit(nums)
    for i := 0; i < nums; i++ {
        g.Go(func() error {
            for job := range jobs {
                result := &Result{
                    job: job,
                    Sum: job.Work(),
                }
                results <- result
            }
            return nil
        })
    }
}

func CreateProducts(jobs chan *Job) {
    idx := 0
    for {
        job := &Job{
            Id:         idx,
            RandNumber: rand.Int(),
        }
        jobs <- job
        idx++
        time.Sleep(time.Millisecond * 50)
    }
}

func main() {
    const workers = 10
    const jobs = 100
    job := make(chan *Job, jobs)
    result := make(chan *Result, jobs)
    go CreateConsumers(job, result, workers)
    go CreateProducts(job)
    for r := range result {
        println("Idx: ", r.job.Id, " Rand:", r.job.RandNumber, " Sum:", r.Sum)
    }

}
