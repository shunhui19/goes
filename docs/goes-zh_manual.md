#Goes
goes是一个通用的，简洁，灵活的Socket框架。

* [简介](#简介)
* [原理](#原理)
* [特性](#特性)
* [简单开发实例](#简单开发实例)
* [安装](#安装)
* [使用](#使用)
* [协议](#协议)
* [connection接口](#connection接口)
    - 方法
        - Send
        - Close
        - GetRemoteAddress
        - GetRemoteIP
        - GetRemotePort
        - GetLocalAddress
        - GetLocalIP
        - GetLocalPort
* [CStore接口](#cstore接口)
    - 方法
        - Set
        - Get
        - Del
        - Range
        - Len
* [Goer](#goer)
    - 属性
        - Transport
        - Protocol
        - Daemon
        - StdoutFile
        - PidFile
        - Connections
    - 回调属性
        - OnConnect
        - OnMessage
        - OnClose
        - OnError
        - OnBufferFull
        - OnBufferDrain
        - OnStop
        - OnReload
    - 方法
        - NewGoer
        - RunAll

* [TCPConnection](#tcpconnection)
    - 属性
        - ID
        - Protocol
        - MaxSendBufferSize
        - MaxPackageSize
        - Connections
    - 回调属性
        - OnMessage
        - OnClose
        - OnError
        - OnBufferFull
        - OnBufferDrain
        - Send
    - 方法
        - Close
        - GetRemoteAddress]
        - GetRemoteIP
        - GetRemotePort
        - GetLocalAddress
        - GetLocalIP
        - GetLocalPort
***

## 简介
   Goer 是一个通用的，简洁，灵活的Socket框架, 其设计思想主要是参考workerman。
## 原理
   **基于go语言协程便利，开一个goroutine(协程)只需关键字go func()，并且消耗资源也非常小**

   整体架构采用的是：主协程(main-goroutine)+子协程(child-goroutine)模式
   - 主协程主要负责：
       1. 指定协议并监听端口, 协议包括传输层协议(TCP或UDP)和应用层协议(自定义)
       2. 注册信号并监听信号, 信号包括停止服务，平滑启动信号等等
       3. 接收客户端连接，有新的连接过来时，根据协议类型并实例化对应的连接对象，开启一个goroutine提供服务
       4. 所有连接管理

   - 子协程负责: 每个子协程负责处理一个连接, 连接的数据发送与接收


## 特性
- 支持TCP, UDP

    同时支持TCP，UDP两种传输层协议，只需要更改传输协议，业务代码无需改动。
- 支持自定义应用层协议

    可根据实际情况定义符合业务的协议，应用层协议定义了一个接口，只需要实现这个接口即可。
    目前已经内置协议有text文本协议，简单的http协议，后续会增加websocket协议。
- 平滑重启

    当需要更新业务或发布新版本，旧的连接不会断开，直到客户端断开连接或超过设置超时重启时间;
    执行重启命令，会创建一个子进程来接收新的连接，父进程不会接收新的连接，直到所有旧连接处理完，然后退出。
- 定时器

    定时器采用时间轮算法实现，支持各种时间级别定时任务。
- 多种运行模式

    支持调试(debug)模式，后台(daemon)模式，只需配置一个字段即可。
- 日志重定向

    支持输出与错误重定向到文件，方便调试。
- 协程模型简单

    采用经典的main-goroutine+child-goroutine来处理用户请求。


## 简单开发实例
- 使用TCP协议对外提供服务, 文件: goer_tcp.go

    ```
    package main

    import (
        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
    )

    func main() {
        // 直接使用传输层TCP协议，监听8080端口
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")

        // 当客户端有消息时，发送'hello, goer'给客户端
        goer.OnMessage = func(connection connections.Connection, data []byte) {
            connection.Send("hello, goer", false)
        }

        // 启动服务
        goer.RunAll()
    }
    ````
  **启动命令**
  ```
  go goer_tcp.go start
  ```
  **测试**
  ```
  frank:~ frank$ nc localhost 8080
  hello         // 客户端发送'hello'
  hello, goer   // 收到服务端回复'hello, goer'
  ```
  *****

- 使用内置Http协议对外提供服务, 文件: goer_http.go
    ```
    package main

    import (
        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
        "github.com/shunhui19/goes/protocols/http"
    )

    func main() {
        // 使用http协议对外提供服务，监听8080端口
        goer := goes.NewGoer("127.0.0.1:8080", http.NewHttpProtocol(), "tcp")

        // 当客户端有消息时，发送'hello, goer'给客户端
        goer.OnMessage = func(connection connections.Connection, data []byte) {
            connection.Send("hello, goer", false)
        }

        // 启动服务
        goer.RunAll()
    }
    ````
  **启动命令**
  ```
  go goer_http.go start
  ```
  **测试**
  ```
  直接打开浏览器，输入地址：localhost:8080
  ```

## 安装
```
go get github.com/shunhui19/goes
```

## 使用
假设启动脚本文件为: server.go
- 启动, 可分为开发(debug)模式启动和守护(daemon)模式启动

    **开发模式启动**
    ```
    go server.go start
    ```
    **守护模式启动**
    ```
    go server.go start -d
    ```

- 停止, 停止服务
    ```
    go server.go stop
    ```
- 平滑重启, 当发布新版本或更新业务时, 可使服务不中断完成升级, 提升用户体验
   ```
    go server.go reload
    ```

## 协议
传输层有两种协议，TCP协议和UDP协议,
TCP协议相对于UDP协议，主要特点是：**面向连接的，字节流和可靠传输**, 所谓面向连接,就是通信双方必须先三次握手建立连接才能通信,
字节流就是数据像水一样从一端(比如服务端或客户端)流向另一端, 所以为了区分每次发送的数据，就必须定一个规则从数据流中取出每次发送的数据，
这个规则就是**应用层自定义协议**

 goes中定义了自定义协议**Protocol接口**,该接口封装三个方法，只要实现这三个方法即可完成新的协议制定，详情查看对应文档

## connection接口
该接口针对传输层协议的连接接口，该接口包含的方法用于对连接相关操作
* #### Send(data string, raw bool) interface{}
        说明:用当前连接发送数据
        参数:
        [data string] 待发送的数据。
        [raw bool]    是否发送原始数据，即不经过应用层协议编码，直接发送数据。
        返回值:
        interface{}     消息发送状态
        当返回nil时表示数据已经发送到应用层发送缓冲区，等待发送到系统内核发送缓冲区;
        当返回false时，发送失败，返回true表示已经发送到系统内核发送缓冲区。
    实例:
    ```
    package main

    import (
        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
        "github.com/shunhui19/goes/protocols"
    )

    func main() {
        // 这里用应用层自定义Text文本协议通信
        goer := goes.NewGoer("127.0.0.1:8080", protocols.NewTextProtocol(), "tcp")

        goer.OnMessage = func(connection connections.Connection, data []byte) {
            // 这里调用接口(connects.Connection)中的Send(data string, raw bool)方法发送数据,
            // 数据会在发送之前自动调用TextProtocol协议，按text协议把数据编码之后发给客户端
            connection.Send("hello, goes", false) // 发送的数据为 "hello, goes\n"
            // 数据不编码，直接发送原始数据
            connection.Send("hello, goes", true) // 发送的数据为 "hello, goes"
        }

        goer.RunAll()
    }
    ```
* #### Close(data string)
        说明:服务端关闭客户端连接
        参数:
        [data string]   关闭时待发送的数据, 可用于通知客户端，服务端准备关闭连接。
        返回值:
        无
    实例:
    ```
    package main

    import (
        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")

        goer.OnMessage = func(connection connections.Connection, data []byte) {
            // 收到客户端消息后，关闭客户端连接
            connection.Close("server is reload...")
            // 直接关闭，不发送数据
            // connection.Close("")
        }

        goer.RunAll()
    }
    ```
* #### GetRemoteAddress() string
        说明:获取连接客户端地址
        参数:
        无
        返回值:
        [string]    客户端地址，形式: "192.168.1.1:12345"
* #### GetRemoteIP() string
        说明:获取连接客户端IP
        参数:
        无
        返回值:
        [string]    客户端IP，形式: "192.168.1.1"
* #### GetRemotePort() int
        说明:获取客户端端口
        参数:
        无
        返回值:
        [int]    客户端IP，形式: 12345
* #### GetLocalAddress() string
        说明:获取连接服务端地址
        参数:
        无
        返回值:
        [string]    服务端地址，形式: "127.0.0.1:8080"
* #### GetLocalIP() string
        说明:获取连接服务端IP
        参数:
        无
        返回值:
        [string]    服务端IP，形式: "127.0.0.1"
* #### GetLocalPort() int
        说明:获取服务端端口
        参数:
        无
        返回值:
        [int]   服务端IP，形式: 8080

## cstore接口
该接口对TCPConnection连接对象的存储，查询，遍历，获取
* #### Set(conn *TCPConnection)
        说明:存储一个*TCPConnection, 该方法主要是goes内部使用，业务上基本用不到
        参数:
        *TCPConnection  tcp连接对象指针
        返回值:
        无
* #### Get(connectionID int)
        说明:根据connectionID连接唯一ID，获取对应的连接对象指针, 该方法主要是goes内部使用，业务上基本用不到
        参数:
        [int]   connectionID
        返回值:
        *TCPConnection  tcp连接对象指针
* #### Del(connectionID int)
        说明:根据connectionID连接唯一ID，删除连接对象, 该方法主要是goes内部使用，业务上基本用不到
        参数:
        [int]   connectionID
        返回值:
        无
* #### Range(f func(key, value interface{}) bool)
        说明:参数是一个函数，通过该函数遍历存储连接,业务中主要用到这个方法
        参数:
        func(key, value interface{}) bool
        返回值:
        无

    实例:
    ```
    package main

    import (
        "fmt"

        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")

        // 当有新的客户端连接过来时，通知其它在线客户端有新的连接上线了
        goer.OnConnect = func(connection connections.Connection) {
            // 遍历所有连接
            goer.Connections.Range(func(key, value interface{}) bool {
                if otherConnection := value.(*connections.TCPConnection); otherConnection.GetRemoteAddress() != connection.GetRemoteAddress() {
                    otherConnection.Send(fmt.Sprintf("a new client[%v] is online", connection.GetRemoteAddress()), false)
                }
                return true
            })
        }

        goer.RunAll()
    }
    ```
* #### Len()
        说明:获取连接总数,业务中基础用不到这个方法
        参数:
        无
        返回值:
        int     连接总数

## goer
Goer结构体是goes的核心对象，它接收连接，处理连接上的数据,并通过一系列回调函数实现业务处理
### 属性
- #### Transport
  说明:
  ```
  string    传输层协议, 包括tcp4, tcp, tcp6, unix unixpacket, ssl, udp4, udp, udp6, unixgram
  ```
目前只支持tcp, tcp4, udp, udp4, 如果为空, 则使用默认值为"tcp"

    实例:
    使用TCP协议
    ```
    package main

    import (
        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "")
        // 使用tcp协议
        goer.Transport = "tcp"
        goer.OnMessage = func(connection connections.Connection, data []byte) {
            connection.Send("hello, goes", false)
        }

        goer.RunAll()
    }
    ```
- #### Protocol
  说明:
  ```
  string    应用层自定义协议
  ```
可以为nil，表示不使用应用层协议，使用传输层协议

    实例:
    使用text文本协议
    ```
    package main

    import (
        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
        "github.com/shunhui19/goes/protocols"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")
        // 使用text文本协议
        goer.Protocol = protocols.NewTextProtocol()
        goer.OnMessage = func(connection connections.Connection, data []byte) {
            connection.Send("hello, goes", false)
        }

        goer.RunAll()
    }
    ```
- #### Daemon
  说明:
  ```
  bool  服务启动是否以守护进程运行, 此属性与启动时命令行执行 -d 参数效果相同
  ```
默认为 false, 当以daemon模式运行时，停止服务执行 go ./executeFile(可执行文件) stop 命令

    实例:
    ```
    package main

    import (
        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")
        // 以守护进程模式启动服务
        goer.Daemon = true
        goer.OnMessage = func(connection connections.Connection, data []byte) {
            connection.Send("hello, goes", false)
        }

        goer.RunAll()
    }
    ```
- #### StdoutFile
  说明:
  ```
  string  输出重定向文件, 所有输出和错误信息都将写入此文件，此属性只支持在Daemon模式下运行
  ```
当以Daemon模式运行，不设置StdoutFile属性，则输出到 /dev/null

    实例:
    ```
    package main

    import (
        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")

        goer.Daemon = true
        // 指定重定向文件
        goer.StdoutFile = "./out.log"
        goer.OnMessage = func(connection connections.Connection, data []byte) {
            connection.Send("hello, goes", false)
        }

        goer.RunAll()
    }
    ```
- #### PidFile
  说明:
  ```
  string    服务进程Pid值存储文件, 该文件用于在执行平滑重启命令 go ./exectFile reload
  ```
默认存储在 goes 项目根目录内

    实例:
    ```
    package main

    import (
        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")

        // 指定PidFile文件
        goer.PidFile = "./myPid.pid"
        goer.OnMessage = func(connection connections.Connection, data []byte) {
            connection.Send("hello, goes", false)
        }

        goer.RunAll()
    }
    ```
- #### Connections
  说明:
  ```
  connections.CStore    接口类型,该接口定义了存储连接的增删改查方法，主要用于给所有连接客户端发送消息等等
  ```
该接口目前只有一个结构体 ConnStore 实现

    实例
    ```
    package main

    import (
        "fmt"

        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")

        // 当有新的客户端连接过来时，通知其它在线客户端有新的连接上线了
        goer.OnConnect = func(connection connections.Connection) {
            goer.Connections.Range(func(key, value interface{}) bool {
                if otherConnection := value.(*connections.TCPConnection); otherConnection.GetRemoteAddress() != connection.GetRemoteAddress() {
                    otherConnection.Send(fmt.Sprintf("a new client[%v] is online", connection.GetRemoteAddress()), false)
                }
                return true
            })
        }

        goer.RunAll()
    }
    ```

## tcpconnection
TCPConnection结构体是goes的核心对象，每个TCPConnection实例表示一个TCP连接
### 属性
