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

- ### 0. 准备工作

  - **Base URL**: `http://127.0.0.1:8080`
  - **权限机制**:
    - 登录接口会返回一个 `token`。
    - 所有“受保护路由”均需在 Postman 的 **Headers** 中添加：
      - Key: `Authorization`
      - Value: `Bearer <你的TOKEN>`
  - **数据格式**: Body 选择 `raw` -> `JSON`。
  
  ------
  
  ### 1. 用户模块 (User)
  
  #### 1.1 用户注册
  
  - **接口**: `POST /user/register`
  - **功能**: 注册新账号。
  - **Body**:
  
  JSON
  
  ```
  {
      "username": "testuser",
      "password": "mysecretpassword"
  }
  ```
  
  #### 1.2 用户登录
  
  - **接口**: `POST /user/login`
  - **功能**: 验证身份并获取 Token。
  - **Body**: 同上。
  - **提示**: 复制返回结果中的 `token` 字段，用于后续接口。
  
  #### 1.3 注销账号 (需登录)
  
  - **接口**: `POST /user/delete`
  - **功能**: 永久删除当前登录的账号。
  - **Header**: `Authorization: Bearer <TOKEN>`
  
  ------
  
  ### 2. 问题模块 (Question)
  
  #### 2.1 发布问题 (需登录)
  
  - **接口**: `POST /question/create`
  - **Body**:
  
  JSON
  
  ```
  {
      "title": "如何学习 Golang?",
      "content": "求推荐一些好的学习路径和项目。"
  }
  ```
  
  #### 2.2 修改问题 (需登录 & 仅限作者)
  
  - **接口**: `POST /modify/question`
  - **说明**: 代码中会严格校验 `userID` 是否为该问题原作者。
  - **Body**:
  
  JSON
  
  ```
  {
      "id": 1, 
      "title": "如何学习 Golang (更新版)",
      "content": "我想找关于 Gin 框架的教程。"
  }
  ```
  
  #### 2.3 删除问题 (需登录 & 仅限作者)
  
  - **接口**: `POST /question/delete?id=1`
  - **方法**: `POST` (注意：代码中使用 `c.Query("id")` 获取参数)
  - **参数**: 在 URL 后拼接 `?id=问题ID`。
  
  ------
  
  ### 3. 回答与评论模块 (Answer)
  
  #### 3.1 提交回答 (需登录)
  
  - **接口**: `POST /answer/create`
  - **Body**:
  
  JSON
  
  ```
  {
      "question_id": 1,
      "content": "多写代码，多看标准库源码。"
  }
  ```
  
  #### 3.2 回复他人的回答 (需登录)
  
  - **接口**: `POST /answer/reply`
  - **功能**: 针对某个现有的回答进行追问或回复。
  - **Body**:
  
  JSON
  
  ```
  {
      "answer_id": 1,
      "content": "握手，我也是这么想的。"
  }
  ```
  
  #### 3.3 修改/删除回答 (需登录 & 仅限作者)
  
  - **修改**: `POST /modify/answer` (Body 需带回答 `id`)
  - **删除**: `POST /answer/delete?id=1` (Query 参数拼接)
  
  ------
  
  ### 4. 社交功能 (Social)
  
  #### 4.1 关注用户 (需登录)
  
  - **接口**: `POST /follow/:id`
  - **示例**: `POST http://127.0.0.1:8080/follow/5` (关注 ID 为 5 的用户)
  - **注意**: 不能关注自己，否则会报错。
  
  #### 4.2 取消关注 (需登录)
  
  - **接口**: `POST /unfollow/:id`
  - **示例**: `POST http://127.0.0.1:8080/unfollow/5`
