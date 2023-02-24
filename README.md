# Tiktok

基于gin框架实现了第五届字节跳动青训营互动方向的简易版抖音，完成了视频流、投稿、个人主页的基础功能和用户喜欢、评论的互动功能，且具有丰富的可观测性

地址：https://github.com/Ltd5552/Tiktok

## 技术栈

### 路由

> **Gin**

项目整体采用Go语言开发，选用了Gin作为HTTP框架

### 存储

> **Mysql+Gorm+MinIO**

考虑到应用场景和数据类型，

- 一般数据存储使用**mysql**关系型数据库，并使用**Gorm**访问
- 使用**ffmpeg**进行封面截取，上传的视频和封面使用**MinIO**进行对象存储

数据库均由docker容器化部署

### 鉴权与配置

> **JWT+Viper**

采用jwt进行用户鉴权，使用viper对配置文件进行全局、有效的管理

### 可观测性

> **Prometheus+zap+otelgin+Grafana+Jaeger**

代码基于gin中间件引入prom_goclient和otelgin实现metric和trace，使用再封装的zap实现traceID的日志输出，可视化采用Grafana组件并基于自定义格式对R.E.D等指标进行展示，Jaeger对trace进行捕捉分析，组件的部署均基于docker容器完成

## 设计

### 数据库设计

**ERD**

![ERD](/img/Tiktok_ERD.jpg)

### 整体架构

![JG](/img/Tiktok_JG.png)



## 代码

### 目录结构

- config用于viper初始化以及配置读取管理
- controller进行路由的初始化定义、接受和处理
- img存储一些图片
- internal存放sql代码
- model用于数据库初始化定义和操作
- pkg作为工具包，含hash、jwt、log、metric、minio、trace代码的定义和初始化
- main.go主要作为项目初始化

### 索引创建

基于使用场景

- 对用户登录场景的name、password创建索引
- 对查看用户发布视频场景的author_id创建索引
- 对根据时间倒序的feed流场景的created_at创建索引
- 对查看视频评论场景的video_id创建索引

```SQL
ALTER table users ADD INDEX name(name);
ALTER table users ADD INDEX password(password);
ALTER table videos ADD INDEX author_id(author_id);
ALTER table videos ADD INDEX created_at(created_at);
ALTER table comments ADD INDEX video_id(video_id);
```

## 存在问题

1.  视频播放较慢
2.  视频上传大小有限制（避免网络原因导致超时）
3.  jaeger接收器反应较慢

## 下一步计划

### 架构上

1. 使用高性能的hertz框架代替gin
2. 微服务化，考虑使用gRPC或Kitex
3. 全容器化处理，将程序打包为镜像
4. 使用微服务全链路透明化

### 优化上

1. 对于高访问的内容使用redis做缓存处理
2. 对于点赞、评论、发布视频等等操作可以引入消息队列，提高响应速度，优化用户体验
3. 部分代码存在冗余、重复结构体问题，有待更新
4. 部署较多，引入docker-compose
5. 代码修改再部署麻烦，引入cicd