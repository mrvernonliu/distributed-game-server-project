package proposed

import "sync"

type WorkerPool struct {
	activeQueue chan Worker
	activeMux sync.Mutex

	idleQueue chan Worker
	idleMux sync.Mutex
}

func (workerPool *WorkerPool) popActive() Worker {
	workerPool.activeMux.Lock()
	worker := <- workerPool.activeQueue
	workerPool.activeMux.Unlock()
	return worker
}

func (workerPool *WorkerPool) popIdle() Worker {
	workerPool.idleMux.Lock()
	worker := <- workerPool.idleQueue
	workerPool.idleMux.Unlock()
	return worker
}

func (workerPool *WorkerPool) pushActive(worker Worker) {
	workerPool.activeMux.Lock()
	workerPool.activeQueue <- worker
	workerPool.activeMux.Unlock()
}

func (workerPool *WorkerPool) pushIdle(worker Worker) {
	workerPool.idleMux.Lock()
	workerPool.idleQueue <- worker
	workerPool.idleMux.Unlock()
}


func CreateWorkerPool(workerList [] Worker) *WorkerPool {
	workerPool := WorkerPool{
		activeQueue: make(chan Worker, len(workerList)),
		activeMux:   sync.Mutex{},
		idleQueue:   make(chan Worker, len(workerList)),
		idleMux:     sync.Mutex{},
	}
	for _, worker := range workerList {
		workerPool.idleQueue <- worker
	}
	return &workerPool
}
