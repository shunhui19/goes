Goer 是一个通用的，灵活的Socket框架, 其设计思想主要是参考workerman。

#### 特性
* 支持TCP, UDP  
    同时支持TCP，UDP两种传输层协议，只需要更改传输协议，业务代码无需改动。
* 支持自定义应用层协议  
    可根据实际情况定义符合业务的协议，应用层协议定义了一个接口，只需要实现这个接口即可。
    目前已经内置协议有text文本协议，简单的http协议，后续会增加websocket协议。
* 平滑重启  
    当需要更新业务或发布新版本，旧的连接不会断开，直到客户端断开连接或超过设置超时重启时间; 
    执行重启命令，会创建一个子进程来接收新的连接，父进程不会接收新的连接，直到所有旧连接处理完，然后退出。
* 定时器  
    定时器采用时间轮算法实现，支持各种时间级别定时任务。
* 多种运行模式  
    支持调试(debug)模式，后台(daemon)模式，只需配置一个字段即可。
* 协程模型简单  
    采用经典的main-goroutine+child-goroutine来处理用户请求，后续会增加协程池来处理百万协程数量复用。
    
#### 性能测试  
    **测试环境**
    ```
    CPU         Inter(R) Xeon(R) Platinum 8255C CPU @ 2.50GHz
    OS          Ubuntu Server 16.04.1 LTS 64位
    Memery      8G
    TestSoft    ab
    ```
    **测试脚本**
    ```
    package main
    
    import (
    	"goes"
    	"goes/connections"
    )
    
    func main() {
    	goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")
    	goer.OnMessage = func(connection connections.Connection, data []byte) {
    		connection.Send("HTTP/1.1 200 OK\r\nConnection: keep-alive\r\nServer: goes\\0.1\r\nContent-Length: 5\r\n\r\nhello", false)
    	}
    
    	goer.RunAll()
    }
    ```
    **测试报告**
    ```
    ab -n1000000 -c100 -k http://127.0.0.1:8080/
    This is ApacheBench, Version 2.3 <$Revision: 1706008 $>
    Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
    Licensed to The Apache Software Foundation, http://www.apache.org/
    
    Benchmarking 127.0.0.1 (be patient)
    
    
    Server Software:        goes\0.1
    Server Hostname:        127.0.0.1
    Server Port:            8080
    
    Document Path:          /
    Document Length:        5 bytes
    
    Concurrency Level:      100
    Time taken for tests:   8.272 seconds
    Complete requests:      1000000
    Failed requests:        0
    Keep-Alive requests:    1000000
    Total transferred:      85000000 bytes
    HTML transferred:       5000000 bytes
    Requests per second:    120886.34 [#/sec] (mean)
    Time per request:       0.827 [ms] (mean)
    Time per request:       0.008 [ms] (mean, across all concurrent requests)
    Transfer rate:          10034.51 [Kbytes/sec] received
    
    Connection Times (ms)
                  min  mean[+/-sd] median   max
    Connect:        0    0   0.0      0       2
    Processing:     0    1   0.2      1      12
    Waiting:        0    1   0.2      1      12
    Total:          0    1   0.2      1      12
    
    Percentage of the requests served within a certain time (ms)
      50%      1
      66%      1
      75%      1
      80%      1
      90%      1
      95%      1
      98%      1
      99%      2
     100%     12 (longest request)
    ```
#### 使用实例 
