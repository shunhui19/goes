package main

import (
	"fmt"
	"goes/lib"
	"time"
)

func main() {
	timeFormat := "2006-01-02 15:04:05.9999"

	// 秒级定时
	t := lib.NewTimer(10, 1*time.Second)
	// 毫秒级定时
	//t := lib.NewTimer(10, 100*time.Millisecond)

	// 5秒后执行一次定时任务
	t.Add(5*time.Second, func(v ...interface{}) {
		// 这里是具体回调函数要执行的内容
		fmt.Printf("[%v]: %v\n", time.Now().Format(timeFormat), v)
	}, "after 5 second to run", false)

	// 2秒后执行，一直循环定时任务
	timerID := t.Add(2*time.Second, func(v ...interface{}) {
		fmt.Printf("[%v]: %v\n", time.Now().Format(timeFormat), v)
	}, "2秒循环定时任务", true)

	// 10秒钟后删除循环定时任务
	t.Add(10*time.Second, func(v ...interface{}) {
		fmt.Printf("[%v]: 开始执行删除操作\n", time.Now().Format(timeFormat))
		if ok := t.Del(timerID); ok {
			fmt.Println("删除成功")
		} else {
			fmt.Println("删除失败")
		}
	}, timerID, false)

	// 启动运行
	fmt.Printf("[%v]: start...\n", time.Now().Format(timeFormat))
	t.Start()

	select {}
}
