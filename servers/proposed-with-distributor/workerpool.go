package proposedWithDistributor

import (
	"sync"
)

type WorkerAddress struct {
	Address string
	Port string
}

type WorkerPool struct {

	idleQueue chan WorkerAddress
	idleMux sync.Mutex
}


func (workerPool *WorkerPool) popIdle() WorkerAddress {
	workerPool.idleMux.Lock()
	var worker WorkerAddress
	if len(workerPool.idleQueue) > 0 {
		worker = <- workerPool.idleQueue
	} else {
		worker.Address = "-1"
	}
	workerPool.idleMux.Unlock()
	return worker
}


func (workerPool *WorkerPool) pushIdle(workerAddress WorkerAddress) {
	workerPool.idleMux.Lock()
	workerPool.idleQueue <- workerAddress
	workerPool.idleMux.Unlock()
}


func CreateWorkerPool(workerList [] Worker) *WorkerPool {
	workerPool := WorkerPool{
		idleQueue:   make(chan WorkerAddress, len(workerList)),
		idleMux:     sync.Mutex{},
	}
	for _, worker := range workerList {
		workerPool.idleQueue <- WorkerAddress{
			Address: worker.Address,
			Port:    worker.Port,
		}
	}
	return &workerPool
}
