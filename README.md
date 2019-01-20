
## apiserver

基于Golang搭建 RESTful API 服务

### 准备工作
```shell
docker run -d -p 3306:3306 --name mymysql \
-v $PWD/conf:/etc/mysql/conf.d -v $PWD/logs:/logs -v $PWD/data:/var/lib/mysql \
-e MYSQL_ROOT_PASSWORD=123456 mysql:5.7

mysql -h 127.0.0.1 -P 3306 -u root -p

mysql>  source db.sql

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

### 注意
api身份验证用的是 JSON Web Token
```shell
# 登录
curl -XPOST -H "Content-Type: application/json" http://127.0.0.1:8080/login -d'{"username":"admin","password":"admin"}'

# 请求头带上token
curl -XPOST -H "Authorization: Bearer ${token}" -H "Content-Type: application/json" http://127.0.0.1:8080/v1/user -d'{"username":"user1","password":"user1234"}'
```
* 可用postman做接口测试，login之后把token保存到环境变量或者全局变量，请求头带上tokne  
* 也可用项目根目录下的test.py测试