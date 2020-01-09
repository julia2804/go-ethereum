package ebtree_v2

type Task struct {
	Id      int
	Err     error
	Prepool *WorkerPool
	f       func(id int, prepool *WorkerPool) *TaskR
}

func (task *Task) Do() *TaskR {
	return task.f(task.Id, task.Prepool)
}

type WorkerPool struct {
	PoolSize    int
	tasksSize   int
	tasksChan   chan Task
	resultsChan chan TaskR
	Results     func(e int) *[]TaskR
}

func NewWorkerPool(tasks *[]Task, poolsize int) *WorkerPool {
	tasksChan := make(chan Task, len(*tasks))
	resultsChan := make(chan TaskR, len(*tasks))
	for _, task := range *tasks {
		tasksChan <- task
	}
	close(tasksChan)
	pool := &WorkerPool{PoolSize: poolsize, tasksSize: len(*tasks), tasksChan: tasksChan, resultsChan: resultsChan}
	pool.Results = pool.results
	return pool
}

func (pool *WorkerPool) Start() {
	for i := 0; i < pool.PoolSize; i++ {
		go pool.worker()
	}
}

func (pool *WorkerPool) worker() {
	for task := range pool.tasksChan {
		re := task.Do()
		pool.resultsChan <- *re
	}
}

func (pool *WorkerPool) results(e int) *[]TaskR {
	results := make([]TaskR, e)
	for i := 0; i < e; i++ {
		results[i] = <-pool.resultsChan
	}
	return &results
}

func ToChannel(id int, prepool *WorkerPool) *TaskR {
	//single task repsonse
	var strps TaskR
	var tmprps []ResultD

	for i := 1; i <= interval; i++ {
		block := bc.GetBlockByNumber(uint64((id * interval) + i))
		if block != nil {
			trans := block.Transactions()
			for j := 0; j < trans.Len(); j++ {
				var tmprd ResultD
				tmprd.Value = trans[j].Value().Bytes()
				var tmptd TD
				tmptd.IdentifierData = Convert2IdentifierData(i, j)
				tmprd.ResultData = append(tmprd.ResultData, tmptd)
				tmprps = append(tmprps, tmprd)
			}
		}
	}
	strps.TaskResult = tmprps
	return &strps
}

func FromChannel(id int, prepool *WorkerPool) *TaskR {
	var rps TaskR
	trps := prepool.Results(takenum)
	//把这些排序
	rps.TaskResult = *mergeSortAndMergeSame(trps)
	return &rps
}

func AssembleTaskAndStart(tasknum int, threadnum int, f func(id int, prepool *WorkerPool) *TaskR, prepool *WorkerPool) *WorkerPool {
	tasks := make([]Task, tasknum)
	for i := 0; i < tasknum; i++ {
		tasks[i] = *new(Task)
		tasks[i].Id = i
		tasks[i].Prepool = prepool
		tasks[i].f = f
	}

	pool := NewWorkerPool(&tasks, threadnum)
	pool.Start()
	return pool
}
