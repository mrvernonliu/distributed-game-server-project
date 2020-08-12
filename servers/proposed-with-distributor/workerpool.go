package proposedWithDistributor

import (
	"sync"
)

type WorkerAddress struct {
	Address string
	Port string
}

type WorkerPool struct {
	count int
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

func CreateWorkerAddressPool(workerList []WorkerAddress) *WorkerPool {
	workerPool := WorkerPool{
		count :		 0,
		idleQueue:   make(chan WorkerAddress, len(workerList)),
		idleMux:     sync.Mutex{},
	}
	for _, worker := range workerList {
		workerPool.count++
		workerPool.idleQueue <- worker
	}
	return &workerPool
}


func CreateWorkerPool(workerList [] Worker) *WorkerPool {
	workerPool := WorkerPool{
		count :		 0,
		idleQueue:   make(chan WorkerAddress, len(workerList)),
		idleMux:     sync.Mutex{},
	}
	for _, worker := range workerList {
		workerPool.count++
		workerPool.idleQueue <- WorkerAddress{
			Address: worker.Address,
			Port:    worker.Port,
		}
	}
	return &workerPool
}
