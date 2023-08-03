package queue

import (
	"k8s.io/klog/v2"
)

// Job is an interface for job
type Job interface {
	Do() error
	// 	用于返回job的信息
	Message() string
}

// JobFunc is a function for job
type JobFunc func() error

// Do implements Job
func (j JobFunc) Do() error {
	return j()
}

type Worker struct {
	JobQueue chan Job  //任务队列
	Quit     chan bool //停止当前任务
}

// NewWorker 新建一个 worker 通道实例   新建一个工人
func NewWorker() Worker {
	return Worker{
		JobQueue: make(chan Job), //初始化工作队列为null
		Quit:     make(chan bool),
	}
}

/*
整个过程中 每个Worker(工人)都会被运行在一个协程中，
在整个WorkerPool(领导)中就会有num个可空闲的Worker(工人)，
当来一条数据的时候，领导就会小组中取一个空闲的Worker(工人)去执行该Job，
当工作池中没有可用的worker(工人)时，就会阻塞等待一个空闲的worker(工人)。
每读到一个通道参数 运行一个 worker
*/

func (w Worker) Run(wq chan chan Job) {
	go func() {
		for {
			//将当前的worker注册到worker队列中
			wq <- w.JobQueue
			select {
			case job := <-w.JobQueue: //从当前的worker队列中取出一个job
				err := job.Do() //执行job
				if err != nil {
					// 记录日志
					klog.Errorf("job[%v] do error: %v", job.Message(), err)
				}
			case <-w.Quit:
				return
			}
		}
	}()
}

//WorkerPool 领导
type WorkerPool struct {
	workerLen   int      //线程池中  worker(工人) 的数量
	JobQueue    chan Job //线程池的  job 通道
	WorkerQueue chan chan Job
}

func NewWorkerPool(workerLen int) *WorkerPool {
	return &WorkerPool{
		workerLen:   workerLen,                      //开始建立 workerLen 个worker(工人)协程
		JobQueue:    make(chan Job),                 //工作队列 通道
		WorkerQueue: make(chan chan Job, workerLen), //最大通道参数设为 最大协程数 workerLen 工人的数量最大值
	}
}

// Run runs the worker pool
func (wp *WorkerPool) Run() {
	//初始化时会按照传入的num，启动num个后台协程，然后循环读取Job通道里面的数据，
	//读到一个数据时，再获取一个可用的Worker，并将Job对象传递到该Worker的chan通道
	klog.V(7).Infof("Starting %d workers", wp.workerLen)
	for i := 0; i < wp.workerLen; i++ {
		//新建 workerLen 20万 个 worker(工人) 协程(并发执行)，每个协程可处理一个请求
		worker := NewWorker() //运行一个协程 将线程池 通道的参数  传递到 worker协程的通道中 进而处理这个请求
		worker.Run(wp.WorkerQueue)
	}

	// 循环获取可用的worker,往worker中写job
	go func() { //这是一个单独的协程 只负责保证 不断获取可用的worker
		for {
			select {
			case job := <-wp.JobQueue: //读取任务
				//尝试获取一个可用的worker作业通道。
				//这将阻塞，直到一个worker空闲
				worker := <-wp.WorkerQueue
				worker <- job //将任务 分配给该工人
			}
		}
	}()
}
