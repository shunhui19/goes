---

---

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
* [Protocol接口](#protocol接口)
    - 方法
        - Input
        - Decode
        - Encode
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
    - 回调属性
        - OnMessage
        - OnClose
        - OnError
        - OnBufferFull
        - OnBufferDrain
    - 方法
        - Send
        - Close
        - GetRemoteAddress
        - GetRemoteIP
        - GetRemotePort
        - GetLocalAddress
        - GetLocalIP
        - GetLocalPort

* [timer定时器](#timer定时器)
    - 方法
        - NewTimer
        - Add
        - Del
        - Start
        - Stop
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
这个规则就是**应用层自定义协议**,在实际编程中就两点：
- 设置数据的边界
- 数据格式定义

以内置的text文本协议为例,以一个换行符作为包的边界：
- 边界设置："\n"就是一个边界标识
- 数据格式定义: data+"\n",这就是一个完整的包

以内置的http协议为例, http协议为协议头header + \r\n\r\n + 协议体body, body可以为空:
- 边界设置; 根据http协议的定义，一个完整的http请求，至少要包括header + \r\n\r\n , 所以边界标识就是 \r\n\r\n
- 数据格式定义: header + \r\n\r\n + body

goes中定义了自定义协议**Protocol接口**,该接口封装三个方法，只要实现这三个方法即可完成新的协议制定,这三个方法为：
```
Input(data []byte, maxPackageSize int) interface{} // 判断当前数据中data是否包含一个边界标识符
Decode(data[]byte) interface{} // 解析数据data格式，并发送给OnMessage回调函数
Encode(data[]byte) []byte // 打包数据data为指定格式，并发给客户端
```

## connection接口
该接口针对传输层协议的连接接口，该接口包含的方法用于对连接相关操作
* #### Send(data string, raw bool) interface{}
        说明:用当前连接发送数据
        参数:
        [data string]       待发送的数据。
        [raw bool]          是否发送原始数据，即不经过应用层协议编码，直接发送数据。
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
                if otherConnection := value.(connections.Connection); otherConnection.GetRemoteAddress() != connection.GetRemoteAddress() {
                    otherConnection.Send(fmt.Sprintf("a new client[%v] is online", connection.GetRemoteAddress()), false)
                }
                return true
            })
        }

        goer.RunAll()
    }
    ```
* #### Len()
        说明:获取连接总数,业务中基本用不到这个方法
        参数:
        无
        返回值:
        int     连接总数

## protocol接口
该接口定义了自定义协议的规则，只需实现三个方法即可完成新协议制定，在启动服务时配置协议字段 Goer.Protocol
* #### Input(data []byte) interface{}
        说明:判断接收的数据是否为一个完整的包
        参数:
        data    []byte  客户端发送的数据
        返回值:
        interface{}		有两种具体类型，bool类型和int类型，分别代表不同含意
        如果包大小超过了最大包长度限制，返回false
        如果不是一个完整的包，则继续接收数据，返回0
        如果收到一个完整的包数据，返回包的长度
* #### Decode(data []byte) []byte
        说明:从data数据中解析一个包
        参数:
        data    []byte  客户端发送的数据
        返回值:
        []byte  解包后的数据
* #### Encode(data []byte) interface{}
        说明:数据发送给客户端之前，按协议规则打包指定格式
        参数:
        data    []byte  发送客户端的数据
        返回值:
        interface{}  打包之后的数据,有多种具体类型，目前只支持 []byte类型和string类型

## goer
Goer结构体是goes的核心对象，它接收连接，处理连接上的数据,并通过一系列回调函数实现业务处理
### 属性
- #### Transport
  说明:
  ```
  string    传输层协议, 包括tcp4, tcp, tcp6, unix unixpacket, ssl, udp4, udp, udp6, unixgram
  ```
目前只支持tcp, tcp4, udp, udp4, 如果为空, 则使用默认值为"tcp"

  实例(使用TCP协议):
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

  实例(使用text文本协议):

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
  			if otherConnection := value.(connections.Connection); otherConnection.GetRemoteAddress() != connection.GetRemoteAddress() {
  				otherConnection.Send(fmt.Sprintf("a new client[%v] is online", connection.GetRemoteAddress()), false)
  			}
  			return true
  		})
  	}

        goer.RunAll()
    }
    ```

### 回调方法
- #### OnConnect(conneciton connections.Connection)
        说明:连接客户端与服务端完成TCP三次握手之后触发回调方法, udp协议无此回调方法
        参数:
        conneciton  connections.Connection 连接接口, TCPConnection实现了这个接口
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

        // OnConnect 回调方法
        goer.OnConnect = func(connection connections.Connection) {
            fmt.Println("a new client is coming, the client address: %v", connection.GetRemoteAddress())
        }

        goer.RunAll()
    }
    ```
- #### OnMessage(conneciton connections.Connection, data []byte)
        说明: 当服务端收到客户端消息时触发回调方法
        参数:
        conneciton  connections.Connection 连接接口
        data        []byte  客户端发送的消息，如果指定应用层协议，则data是解码后的数据
        返回值:
        无

    实例:
    ```
    package main

    import (
        "github.com/shunhui19/goes/lib"

        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")

        // OnMessage 回调方法
        goer.OnMessage = func(connection connections.Connection, data []byte) {
            lib.Info("receive client data: %v", string(data))
        }

        goer.RunAll()
    }
    ```
- #### OnClose(conneciton connections.Connection)
        说明:收到客户端发送的FIN包时触发
        参数:
        conneciton  connections.Connection 连接接口
        返回值:
        无

    实例:
    ```
    package main

    import (
        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
        "github.com/shunhui19/goes/lib"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")

        // OnClose 回调方法
        goer.OnClose = func(connection connections.Connection) {
            lib.Info("client is closing")
        }

        goer.RunAll()
    }
    ```
- #### OnError(conneciton connections.Connection, code int, message string)
        说明:与客户端的连接发生错误时触发,目前主要两个地方会触发
        1.执行connection.Send()方法时，与客户端连接已经断开
        2.发送缓冲区sendBuffer已经满了
        参数:
        conneciton  connections.Connection 连接接口
        code        int     错误码
        message     string  错误信息
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

        // OnMessage 回调函数
        goer.OnError = func(connection connections.Connection, code int, message string) {
            fmt.Printf("error %d, reason %s", code, message)
        }

        goer.RunAll()
    }
    ```
- #### OnBufferFull(conneciton connections.Connection)
        说明:当发送缓冲区满时触发
        参数:
        conneciton  connections.Connection 连接接口
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

        // OnBufferFull 回调函数
        goer.OnBufferFull = func(connection connections.Connection) {
            fmt.Println("the send buff is full")
        }

        goer.RunAll()
    }
    ```
- #### OnBufferDrain(conneciton connections.Connection)
        说明:当发送缓冲区为空时触发
        参数:
        conneciton  connections.Connection 连接接口
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

        // OnBufferDrain 回调函数
        goer.OnBufferDrain = func(connection connections.Connection) {
            fmt.Println("the send buff is empty")
        }

        goer.RunAll()
    }
    ```
- #### OnStop()
        说明:当goes服务停止时触发
        参数:
        无
        返回值:
        无

    实例:
    ```
    package main

    import (
        "fmt"

        "github.com/shunhui19/goes"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")

        // OnStop 回调函数
        goer.OnStop = func() {
            fmt.Println("the server is stop")
        }

        goer.RunAll()
    }
    ```
- #### OnReload()
        说明:当goes执行reload平滑重启命令时触发
        参数:
        无
        返回值:
        无

    实例:
    ```
    package main

    import (
        "fmt"

        "github.com/shunhui19/goes"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")
        goer.StdoutFile = "./reload.log"

        // OnReload 回调函数
        goer.OnReload = func() {
            fmt.Println("the server is reload")
        }

        goer.RunAll()
    }
    ```

### 方法
- #### NewGoer(socketName string, applicationProtocol protocol.Protocol, transportProtocol string) *Goer
        说明:获取Goer实例
        参数:
        socketName  string  监听地址，形式为ip+port, eg: 127.0.0.1:8080
        applicationProtocol protocotl.Protocol  协议接口,自定义应用层协议都需要实现此接口,如果为nil表示直接使用传输层协议，不使用应用层协议
        transportProtocol   string  传输层协议字符串,为空则使用默认"tcp"协议
        返回值:
        *Goer

    实例:
    ```
    package main

    import (
        "github.com/shunhui19/goes"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")

        goer.RunAll()
    }
    ```
- ### RunAll()
        说明:启动Goes服务
        参数:
        无
        返回值:
        无

## tcpconnection
TCPConnection结构体是goes的核心对象，每个TCPConnection实例表示一个TCP连接

### 属性
- #### ID
  说明:
  ```
  int   当前连接对象唯一ID
  ```

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

        goer.OnConnect = func(connection connections.Connection) {
            fmt.Println("the client ID: ", connection.(*connections.TCPConnection).ID)
        }

        goer.RunAll()
    }
    ```
- #### Protocol
  说明:
  ```
  protocols.Protocol    接口类型    当前连接的协议实例，只要实现了protocols.Protocol接口,一般直接在启动服务的时候启动,这针对当前连接有效
  ```

    实例:
    ```
    package main

    import (
        "github.com/shunhui19/goes/protocols"

        "github.com/shunhui19/goes"
        "github.com/shunhui19/goes/connections"
    )

    func main() {
        goer := goes.NewGoer("127.0.0.1:8080", nil, "tcp")

        // 有连接过来时，设置应用层协议为text文本协议,只针对当前连接效
        goer.OnConnect = func(connection connections.Connection) {
            connection.(*connections.TCPConnection).Protocol = protocols.NewTextProtocol()
        }

        goer.RunAll()
    }
    ```
- #### MaxSendBufferSize
  说明:
  ```
  int   当前连接发送缓冲区最大长度, 不设置默认值为 1M，此属性会影响到OnBufferFull回调
  ```
- #### MaxPackageSize
  说明:
  ```
  int   当前连接收发包最大长度, 不设置默认值为 10M
  ```

### 回调函数
回调函数和Goer结构体中的回调函数类似，只是Goer中的回调函数是全局的，即对所有连接生效，而TCPConnection中的回调函数只针对当前连接有效
### 方法
方法connections.Connection接口中的方法, 具体查看[Connection接口](#connection接口)文档

## timer定时器
定时器是基于时间轮算法实现的，原理 [查看](https://www.ibm.com/developerworks/cn/linux/l-cn-timers/)

### 方法
- #### NewTimer(slotNumber int, si time.Duration)
    ```
    说明:获取一个定时器实例
    参数：
    slotNumber  int             转动一圈的格子数
    si          time.Duration   在时间轮转动一格的时间, 理解为转一格的速度
    ```
    实例一(秒级定时):
    ```
    package main

    import (
        "fmt"
        "time"

        "github.com/shunhui19/goes/lib"
    )

    func main() {
        fmt.Printf("[%v]start...\n", time.Now())
        // 定义一圈格式数为 60, 每1秒转动一格
        timer := lib.NewTimer(60, 1*time.Second)

        // 添加一个定时任务， 5秒后执行, 只执行一次
        timer.Add(5*time.Second, func(v ...interface{}) {
            // 这里回调函数里写业务逻辑, 参数 v 是传递的 args 变量
            fmt.Printf("[%v]: %v\n", time.Now(), v)
        }, "hello, goes", false)

        // 开始执行
        timer.Start()

        select {}
    }
    ```
    实例二(毫秒级定时):
    ```
    package main

    import (
        "fmt"
        "time"

        "github.com/shunhui19/goes/lib"
    )

    func main() {
        // 毫秒定时器, 3600个格式，每100毫秒转动一格
        timerMillisecond := lib.NewTimer(3600, 100*time.Millisecond)

        // 添加一个定时任务, 每200毫秒执行一次, 并一直运行
        timerMillisecond.Add(200*time.Millisecond, func(v ...interface{}) {
            fmt.Printf("[%v]: %v\n", time.Now(), v)
        }, "milliSecond timer", true)

        timerMillisecond.Start()

        select {}
    }
    ```

- #### Add(timeInterval time.Duration, fn func(v ...interface{}), args interface{}, persitent bool) TaskID
    ```
    说明:增加一个定时任务
    参数：
    timerInterval   time.Duration           多久之后开始执行
    fn              func(v ...interface)    执行的回调函数
    args            interface{}             回调函数中的参数
    persistent      bool                    是否持久执行
    返回值:
    TaskID          TaskID                  当前添加定时任务的唯一ID标识，用于删除操作
    ```
    实例:
    ```
    package main

    import (
        "fmt"
        "time"

        "github.com/shunhui19/goes/lib"
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
    ```

- #### Del(taskID TaskID)
    ```
    说明:删除一个定时任务
    参数：
    TaskID          TaskID                  添加定时任务返回的唯一ID标识
    返回值:
    无
    ```

- #### Start()
    ```
    说明:启动定时服务
    参数：
    无
    返回值:
    无
    ```

- #### Stop()
    ```
    说明:停止定时服务
    参数：
    无
    返回值:
    无
    ```
