## 项目简介

问答社区是一个基于Go语言和Gin框架开发的Web应用，旨在为用户提供一个在线问答平台。用户可以注册账号、登录、发布问题、回答问题、修改和删除问题和答案。

## 环境依赖

- Go 1.15 或以上版本
- Gin Web框架
- GORM 作为ORM工具
- MySQL 数据库

## 安装指南

1. 确保系统已安装Go环境和MySQL服务。
2. 已经自行获得源码
3. 进入项目目录：`cd [项目名称]`
4. 安装依赖：`go mod tidy`

## 配置文件

项目使用`.env`文件来管理环境变量。你需要在项目根目录下创建一个`.env`文件，并配置数据库连接信息，例如：

```
DB_USERNAME=root
DB_PASSWORD=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=问答社区
DB_TIMEOUT=30s
```

## 数据库迁移

项目使用GORM进行数据库迁移。首次运行前，需要执行数据库迁移以创建所需的数据表。

```
go run main.go
```

## 启动服务

执行以下命令启动问答社区服务：

```
go run main.go
```

服务将默认在`8080`端口启动。

## API文档

问答社区提供以下API接口：

- `POST /register` - 用户注册
- `POST /login` - 用户登录
- `POST /question/create` - 创建问题
- `POST /answer/create` - 创建回答
- `POST /answer/answer` - 回答指定问题
- `POST /modify/question` - 修改问题
- `POST /modify/answer` - 修改回答
- `POST /question/delete` - 删除问题
- `POST /answer/delete` - 删除回答
- `POST /follow/:id` - 关注用户
- `POST /unfollow/:id` - 取消关注用户

## 代码结构

- `configs` - 配置包，用于加载环境变量和数据库连接配置。
- `controllers` - 控制器包，处理HTTP请求和响应。
- `models` - 模型包，定义数据库模型和结构。
- `repositories` - 存储库包，用于数据库操作。
- `routes` - 路由包，定义应用的路由和控制器映射。
- `services` - 服务包，封装业务逻辑。

- - ### 1. 用户模块 (User)

    该模块用于管理用户身份。**登录后获取的 Token 是访问后续所有接口的凭证**。
  
    - **用户注册**
  
      - **URL**: `POST http://localhost:8080/user/register`
  
      - **Body**:
  
        JSON
  
        ```
        {
            "username": "岑曦",
            "password": "123456"
        }
        ```
  
    - **用户登录**
  
      - **URL**: `POST http://localhost:8080/user/login`
  
      - **Body**:
  
        JSON
  
        ```
        {
            "username": "岑曦",
            "password": "123456"
        }
        ```
  
      - **注意**: 登录成功后，请复制返回结果中的 `token` 字段。
  
    - **注销账号**
  
      - **URL**: `POST http://localhost:8080/user/delete`
      - **Header**: `Authorization: Bearer {{token}}`
  
    ------
  
    ### 2. 问题模块 (Question)
  
    管理社区内的提问。修改和删除操作必须由问题的作者执行。
  
    - **创建问题**
  
      - **URL**: `POST http://localhost:8080/question/create`
  
      - **Body**:
  
        JSON
  
        ```
        {
            "title": "如何评价 Go 语言？",
            "content": "请大家从并发模型、开发效率等角度谈谈看法。"
        }
        ```
  
    - **修改问题**
  
      - **URL**: `POST http://localhost:8080/modify/question`
  
      - **Body**: (需带上问题的 `ID`)
  
        JSON
  
        ```
        {
            "ID": 1,
            "title": "如何评价 Go 1.22 版本？",
            "content": "最新的 for 循环语义变化很大，大家怎么看？"
        }
        ```
  
    - **删除问题**
  
      - **URL**: `POST http://localhost:8080/question/delete?id=1`
      - **说明**: 通过 URL 参数 `id` 指定要删除的问题编号。
  
    ------
  
    ### 3. 回答与回复模块 (Answer & Reply)
  
    针对问题进行回答，或对他人的回答进行评论。
  
    - **创建回答**
  
      - **URL**: `POST http://localhost:8080/answer/create`
  
      - **Body**:
  
        JSON
  
        ```
        {
            "question_id": 1,
            "content": "我觉得 Go 最大的优势就是简洁和高并发。"
        }
        ```
  
    - **修改回答**
  
      - **URL**: `POST http://localhost:8080/modify/answer`
  
      - **Body**:
  
        JSON
  
        ```
        {
            "ID": 1,
            "content": "补充一下，标准库也非常强大。"
        }
        ```
  
    - **回复回答 (评论)**
  
      - **URL**: `POST http://localhost:8080/answer/reply`
  
      - **Body**:
  
        JSON
  
        ```
        {
            "answer_id": 1,
            "content": "完全同意你的看法！"
        }
        ```
  
    - **删除回答**
  
      - **URL**: `POST http://localhost:8080/answer/delete?id=1`
  
    ------
  
    ### 4. 社交模块 (Social)
  
    用户之间的关注与取消关注功能。
  
    - **关注用户**
      - **URL**: `POST http://localhost:8080/follow/2`
      - **说明**: 最后的 `/2` 代表你要关注的用户 ID。
    - **取消关注**
      - **URL**: `POST http://localhost:8080/unfollow/2`
      - **说明**: 撤销对 ID 为 2 的用户的关注。
