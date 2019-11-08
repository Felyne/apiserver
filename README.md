
## apiserver

基于Golang的gin框架搭建 RESTful API 服务

### 准备工作
```shell
docker run -d -p 3306:3306 --name mymysql \
-v $PWD/conf:/etc/mysql/conf.d -v $PWD/logs:/logs -v $PWD/data:/var/lib/mysql \
-e MYSQL_ROOT_PASSWORD=123456 mysql:5.7

mysql -h 127.0.0.1 -P 3306 -u root -p

mysql>  source db.sql

```

### 安装
```
# 有些第三方包下载不了，加个代理
export GOPROXY=https://athens.azurefd.net

# 编译
make

```

### 目录结构

        ├── admin.sh                     # 进程的start|stop|status|restart控制文件
        ├── conf                         # 配置文件统一存放目录
        │   ├── config.yaml              # 配置文件
        │   ├── server.crt               # TLS配置文件
        │   └── server.key
        ├── config                       # 专门用来处理配置和配置文件的Go package
        │   └── config.go                 
        ├── db.sql                       # 在部署新环境时，可以登录MySQL客户端，执行source db.sql创建数据库和表
        ├── docs                         # swagger文档，执行 swag init 生成的
        │   ├── docs.go
        │   └── swagger
        │       ├── swagger.json
        │       └── swagger.yaml
        ├── handler                      # 类似MVC架构中的C，用来读取输入，并将处理流程转发给实际的处理函数，最后返回结果
        │   ├── handler.go
        │   ├── sd                       # 健康检查handler
        │   │   └── check.go 
        │   └── user                     # 核心：用户业务逻辑handler
        │       ├── create.go            # 新增用户
        │       ├── delete.go            # 删除用户
        │       ├── get.go               # 获取指定的用户信息
        │       ├── list.go              # 查询用户列表
        │       ├── login.go             # 用户登录
        │       ├── update.go            # 更新用户
        │       └── user.go       
        ├── main.go                      # Go程序唯一入口
        ├── Makefile                     # Makefile文件，一般大型软件系统都是采用make来作为编译工具
        ├── model                        # 数据库相关的操作统一放在这里，包括数据库初始化和对表的增删改查
        │   ├── init.go                  # 初始化和连接数据库
        │   ├── model.go                 # 存放一些公用的go struct
        │   └── user.go                  # 用户相关的数据库CURD操作
        ├── pkg                          # 引用的包
        │   ├── auth                     # 认证包
        │   │   └── auth.go
        │   ├── constvar                 # 常量统一存放位置
        │   │   └── constvar.go
        │   ├── errno                    # 错误码存放位置
        │   │   ├── code.go
        │   │   └── errno.go
        │   ├── token
        │   │   └── token.go
        │   └── version                  # 版本包
        │       ├── base.go
        │       ├── doc.go
        │       └── version.go
        ├── README.md                    # API目录README
        ├── router                       # 路由相关处理
        │   ├── middleware               # API服务器用的是Gin Web框架，Gin中间件存放位置
        │   │   ├── auth.go 
        │   │   ├── header.go
        │   │   ├── logging.go
        │   │   └── requestid.go
        │   └── router.go
        ├── service                      # 实际业务处理函数存放位置
        │   └── service.go
        ├── util                         # 工具类函数存放目录
        │   ├── util.go 
        │   └── util_test.go
        └── vendor                       # vendor目录用来管理依赖包
            ├── github.com
            ├── golang.org
            ├── gopkg.in
            └── vendor.json

### 使用
api身份验证用的是 JSON Web Token
```shell
# 登录返回token信息
curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:8080/login -d'{"username":"admin","password":"admin"}'

# 请求头带上token
# 新增用户
curl -XPOST -H "Authorization: Bearer ${token}" -H "Content-Type: application/json" \
http://127.0.0.1:8080/v1/user -d'{"username":"user1","password":"user1234"}'

# 用户列表
curl -H "Authorization: Bearer ${token}" http://127.0.0.1:8080/v1/user?offset=0&limit=20
```

* 可用postman做接口测试，login之后把token保存到环境变量或者全局变量，请求头带上token
* 可用项目根目录下的test.py测试

### 生成API文档  
首先安装 [swag](https://github.com/swaggo/swag)   
```shell
# 在项目根目录下执行命令会生成docs目录
swag init
```
可以部署[YApi](https://github.com/fjc0k/docker-YApi)来管理项目的api文档，它可以导入`docs/swagger.json`文件

### 性能分析  
第一种:  
直接在浏览器访问 http://localhost:8080/debug/pprof 来查看当前 API 服务的状态，包括 CPU 占用情况和内存使用情况等

第二种:  
安装graphviz
```
apt install graphviz
```
通过 
go tool pprof http://127.0.0.1:8080/debug/pprof/profile，查看cpu使用情况  
go tool pprof http://127.0.0.1:8080/debug/pprof/heap   查看内存使用情况

```
top10
svg
web
```

### http压测

安装wrk
```
apt install libssl-dev
git clone https://github.com/wg/wrk
cd wrk
make
cp wrk /usr/local/bin/
```
简单测试
```
wrk --latency -t16 -d60s -T30s http://127.0.0.1:8080/sd/health -c 200
```

一个生成图表的工具
```
apt install gnuplot
```
真正测试
```
./wrktest.sh
```

### 高可用
[Keepalived+Nginx实现高可用](https://blog.csdn.net/xyang81/article/details/52556886)

### nginx常用命令
```
nginx -s stop       快速关闭 Nginx，可能不保存相关信息，并迅速终止 Web 服务
nginx -s quit       平稳关闭 Nginx，保存相关信息，有安排的结束 Web 服务
nginx -s reload     因改变了 Nginx 相关配置，需要重新加载配置而重载
nginx -s reopen     重新打开日志文件
nginx -c filename   为 Nginx 指定一个配置文件，来代替默认的
nginx -t            不运行，而仅仅测试配置文件。Nginx 将检查配置文件的语法的正确性，并尝试打开配置文件中所引用到的文件
nginx -v            显示 Nginx 的版本
nginx -V            显示 Nginx 的版本、编译器版本和配置参数
```