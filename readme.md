# Simple Web Application

## 个人信息

|      |          |
| ---- | -------- |
| 姓名 | 白家栋   |
| 学号 | 18342001 |
| 专业 | 软件工程 |

## 目录结构

```
.
├── controllers         // 控制器，这里初始化了所有 http server 要处理的响应的逻辑
│   └── handler.go
├── go.mod
├── go.sum
├── main.go             // 入口文件
├── models              // 数据模型
│   └── user.go
├── readme.md           
└── views               // 视图，包含静态html模版，以及css和js文件
    ├── info.css
    ├── info.html
    ├── login.html
    ├── login.js
    └── style.css

```

## 运行
#### 运行环境
- golang 1.14+
- 开启 `GO111MODULE`(linux/mac下: `export GO111MODULE=on`)

#### 运行方法
- 进入项目根目录, 终端下输入:
```
go run main.go
```
此时会开始下载依赖包，并且服务将占用你的主机的3000端口。请确保您的电脑能够成功下载依赖包，如果不能下载成功，可以考虑使用golang的包代理，即（mac, linux下）在运行前:
```
export GO111MODULE=ON
export GOPROXY=https://goproxy.cn
```

## 效果

### 1. 界面
成功运行服务之后，在浏览器中输入: `http://localhost:3000`, 效果如下:
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gktc1z3pcjj31cx0u0auk.jpg)

填入任意的username和password之后，点击`click`, 效果如下:
- 弹出提示框：
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gktc53zjlyj31cx0u0x0l.jpg)
- 点击`确定`后跳转：
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gktc66cjclj31cx0u07p7.jpg)

根据这个用户id，我们也可以通过: `http://localhost:3000/users/147531` 再次访问到上述页面：
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gktc66cjclj31cx0u07p7.jpg)

### 2. curl 测试结果
1. curl http://localhost:3000
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gktcy3gwubj31r20u0x2a.jpg)

2. curl -d "user_id=123&username=baijiadong&password=123456" -X POST http://localhost:3000/
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gktcth17tij31r20u07lp.jpg)

3. curl http://localhost:3000/users/123
![](https://tva1.sinaimg.cn/large/0081Kckwgy1gktd0n2jsyj31r20u0qmr.jpg)

## 实现分析
在服务的路由上，我使用了 `mux`。mux 的特性在于能够匹配路由的模式，比如：
```go
	// 用户页面的处理
	mux.HandleFunc("/users/{user_id}", func(writer http.ResponseWriter, request *http.Request) {
		// 解析url中的参数
		vars := m.Vars(request)

		// 渲染模版
		t, _ := template.ParseFiles("./views/info.html")
		u := models.Users[vars["user_id"]]
		user := models.User{Username: u.Username, Password: u.Password, UserId: u.UserId}
		t.Execute(writer, user)

	})
```
这样，任何符合`/users/:id` 的url请求都会被路由到这里。
在返回静态页面的时候，我使用了 `http/template` 包来对于静态模版进行渲染，比如:
```go
// 渲染模版
		t, _ := template.ParseFiles("./views/info.html")
		u := models.Users[vars["user_id"]]
		user := models.User{Username: u.Username, Password: u.Password, UserId: u.UserId}
		t.Execute(writer, user)
```
这里，首先读取了 `views` 目录下的 info.html 这个 html 模版，同时从 users 数组中根据用户在 url 中传入的id进行查找，将查找到的 `User` 结构体渲染到该模版中，最后写入 `writer` 中。

还有一个实现要点是: 由于服务端在用户通过 POST 方法请求后会将请求的用户名密码写入到 `users` 中，因此面对高并发的情况会遇到多个进程同时修改这个变量的问题，因此这里要加上互斥锁，如下所示:
```go
// 将用户通过表单发送过来的字段存放到内存当中
			// 加锁，防止多个go程同时修改
			mutex.Lock()
			u := models.User{UserId: request.Form["user_id"][0], Username: request.Form["username"][0], Password: request.Form["password"][0]}
			models.Users[u.UserId] = u
			mutex.Unlock()
			
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusCreated)
			json.NewEncoder(writer).Encode(u)
```
## a/b test
- 工具: [apachebench](http://httpd.apache.org/download.cgi#apache24)

- 测试过程
  - 测试1：
    - 命令: ab -n 1000 -c 100 http://localhost:3000/
    - 参数解释: -n 表示请求数量，这里请求1000次；-c 表示并发客户端数，这里为100
    - 结果
      ![](https://tva1.sinaimg.cn/large/0081Kckwgy1gktds53lrjj31b90u01kx.jpg)
      
  - 测试2:
    - 命令: ab -n 1000 -c 100 -p post.txt  -T 'application/x-www-form-urlencoded; charset=UTF-8' -H "X-Requested-With: XMLHttpRequest"  http://localhost:3000/
    - 参数解释: -p 指 post, post.txt中存放了请求的表单, -T 指请求类型， -H包含了请求头部
    - 结果
    ![](https://tva1.sinaimg.cn/large/0081Kckwgy1gktff6cgcqj31gk0u0qtg.jpg)