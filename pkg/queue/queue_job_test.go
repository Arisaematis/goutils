package queue

import (
	"fmt"
	"testing"
	"time"
)

type DoSomeThing struct {
	Num int
}

func (d *DoSomeThing) Do() error {
	fmt.Println("开启线程数：", d.Num)
	time.Sleep(1 * 1 * time.Second)
	return nil
}

func (d *DoSomeThing) Message() string {
	return fmt.Sprintf("当前num[%v]", d.Num)
}

func TestQueue(t *testing.T) {
	//设置最大线程数
	num := 100 * 100 * 20

	// 注册工作池，传入任务
	// 参数1 初始化worker(工人)并发个数 20万个
	p := NewWorkerPool(num)
	p.Run() //有任务就去做，没有就阻塞，任务做不过来也阻塞

	//datanum := 100 * 100 * 100 * 100    //模拟百万请求
	datanum := 100 * 100
	go func() { //这是一个独立的协程 保证可以接受到每个用户的请求
		for i := 1; i <= datanum; i++ {
			sc := &DoSomeThing{Num: i}
			p.JobQueue <- sc //往线程池 的通道中 写参数   每个参数相当于一个请求  来了100万个请求
		}
	}()

	time.Sleep(10 * time.Second)
}
