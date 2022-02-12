package service

import (
	"runtime"
)

type Job struct {
	Args []interface{}
	Run  func()
}

type WokerPool struct {
	jobChan chan Job
}

const MAX_JOBS = 20

func NewWorkPool() *WokerPool {
	pool := &WokerPool{make(chan Job, MAX_JOBS)}
	return pool
}

func (wp *WokerPool) Start() {
	maxWorkers := runtime.GOMAXPROCS(0) - 1
	if maxWorkers > 2 {
		maxWorkers = 2
	} else if maxWorkers <= 0 {
		maxWorkers = 1
	}

	for w := 0; w < maxWorkers; w++ {
		go worker(wp)
	}
}

func (wp *WokerPool) Stop() {
	close(wp.jobChan)
}

func (wp *WokerPool) QueueJob(job *Job) bool {
	select {
	case wp.jobChan <- *job:
		return true
	default:
		return false
	}
}

func worker(pool *WokerPool) {
	for job := range pool.jobChan {
		job.Run()
	}
}
