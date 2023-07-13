package main

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/executors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DailyTask struct {
	insertExecutor *executors.BulkExecutor
}

type TaskInfo struct {
	ID       int64
	TaskName string
	SqlQuery string
}

func (dts *DailyTask) Init() {
	// insertIntoCk() 是真正insert执行函数【需要开发者自己编写具体业务逻辑】
	dts.insertExecutor = executors.NewBulkExecutor(
		dts.insertIntoCk,
		executors.WithBulkInterval(time.Second*3), // 3s会自动刷一次container中task去执行
		executors.WithBulkTasks(32),            // container最大task数。一般设为2的幂次
	)
}

func (dts *DailyTask) insertIntoCk(tasks []any) {
	for _, v := range tasks {
		fmt.Printf("%+v\n", v)
	}
}

func (dts *DailyTask) insertNewData(ch chan interface{}) (err error) {
	for item := range ch {
		r, vok := item.(*TaskInfo)
		if !vok {
			continue
		}
		err := dts.insertExecutor.Add(r)
		if err != nil {
			r.ID = 1
			r.TaskName = "default"
			r.SqlQuery = ""
			// 1. Add Task
			err := dts.insertExecutor.Add(r)
			if err != nil {
				logx.Error(err)
			}
		}
	}
	// 2. Flush Task container
	dts.insertExecutor.Flush()
	// 3. Wait All Task Finish
	dts.insertExecutor.Wait()
	return err
}

func main() {
	DailyTask := DailyTask{}
	DailyTask.Init()
	ch := make(chan interface{}, 10)
	// stop := make(chan struct{})
	go func() {
		id := 1
		for {
			taskInfo := TaskInfo{
				ID:       int64(id),
				TaskName: "inset",
				SqlQuery: "xxx",
			}
			ch <- &taskInfo
			time.Sleep(time.Millisecond * 200)
			id++
			if id == 101{
				close(ch)
				return
			}
		}
	}()
	_ = DailyTask.insertNewData(ch)
}
