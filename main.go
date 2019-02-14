package main

import (
	"fmt"
	"github.com/jakecoffman/cron"
	"time"
)

func main() {
	fmt.Println("cron start")

	DoMyJob()

}

func DoMyJob() {
	fmt.Println("corn 初始化")
	c := cron.New()
	c.Start()
	defer c.Stop()

	fmt.Println("添加定时启动任务")
	// 定时启动任务  11:05启动
	specStart := "0 19 11 * * ?"
	jobStart := NewAddMyJob(c, StartDynamicJobs)
	c.AddJob(specStart, jobStart, "jobStart")

	fmt.Println("添加定时停止任务")
	//	定时停止任务 11:20停止
	specStop := "0 20 11 * * ?"
	jobStop := NewAddMyJob(c, StopDynamicJobs)
	c.AddJob(specStop, jobStop, "jobStop")

	// 阻塞
	select {}
	fmt.Println("任务结束了")
}

// 添加动态任务
func StartDynamicJobs(c *cron.Cron) {

	go func() {
		// 模拟从数据库中获取数据
		d := makeData()
		for _, v := range d {
			// 添加动态任务
			c.AddJob(v.Spec, NewAddDoJob(v.Id, ShowInfo), v.Name)
			fmt.Printf("任务 %s 添加了\n", v.Name)
		}
	}()
}

// 移除动态任务
func StopDynamicJobs(c *cron.Cron) {
	go func() {
		d := makeData()
		for _, v := range d {
			//	移除动态任务
			c.RemoveJob(v.Name)
			fmt.Printf("任务 %s 移除了\n", v.Name)
		}
	}()
}

// 动态任务具体执行方法
func ShowInfo(n int) {
	fmt.Printf("任务执行时间：%s，收到参数：%d \n", time.Now().Format("2006-01-02 15:04:05"), n)
}

// 模拟从数据库中获取待执行的任务信息
func makeData() []TmpModel {
	dlist := []TmpModel{}
	d1 := TmpModel{Name: "1", Spec: "*/5 * * * * ?", Id: 1}
	d2 := TmpModel{Name: "2", Spec: "*/7 * * * * ?", Id: 2}
	dlist = append(dlist, d1)
	dlist = append(dlist, d2)
	return dlist
}

type TmpModel struct {
	Name string // 任务名称
	Spec string // 任务执行crontab
	Id   int    // 执行时为方法传入的参数
}

// **实现带参数的扩展方法***

// 实现一个带有参数为 cron.Cron 的扩展方法
type FuncMyJob func(*cron.Cron)

type MyJob struct {
	c        *cron.Cron
	function FuncMyJob
}

// 实现Run方法
func (t *MyJob) Run() {
	if t.function != nil {
		t.function(t.c)
	}
}

func NewAddMyJob(c *cron.Cron, job FuncMyJob) *MyJob {
	instance := &MyJob{
		c,
		job,
	}
	return instance
}

//******************

// 实现一个带有参数为 int 的扩展方法
type FuncDoJob func(n int)

type DoJob struct {
	n        int
	function FuncDoJob
}

func (t *DoJob) Run() {
	if t.function != nil {
		t.function(t.n)
	}
}

func NewAddDoJob(n int, job FuncDoJob) *DoJob {
	instance := &DoJob{
		n,
		job,
	}
	return instance
}

//******************
