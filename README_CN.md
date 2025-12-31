# Cato

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/NCUHOME/cato/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/NCUHOME/cato)](https://github.com/NCUHOME/cato/issues)
[![GitHub stars](https://img.shields.io/github/stars/NCUHOME/cato)](https://github.com/NCUHOME/cato/stargazers)

> 一个用于构建稳定、可维护的 Golang 服务端项目的强大代码生成工具链。

**⚠️ 重要通知：项目处于积极开发阶段**
此项目目前正在紧张开发中，**随时可能引入不兼容的更改**。请勿在生产环境中使用。

## ✨ Cato 能做什么？

Cato 当前的能力通过 `proto-gen-*` 系列的插件来提供。

### proto-gen-cato

#### 注意 ⚠️
`proto-gen-cato` 要求你的 `.proto` 文件内容符合 [Protocol Buffers 官方风格指南](https://protobuf.dev/programming-guides/style/)。否则，生成的结构体字段名可能会难以理解。

该插件通过在不同作用域下使用 [`proto`](https://github.com/NCUHOME/cato/tree/main/proto) 目录下的自定义选项来提供不同能力。具体使用姿势请参照 [cato-example-bms (简易图书管理系统)](https://github.com/NCUHOME/cato-example-bms/tree/main/proto) 中的 proto 文件。

`proto` 目录下各文件的职责划分如下：
+ **`extension.proto`**：插件选项的主要使用入口。
+ **`db.proto`**：与关系型数据库相关的选项。
+ **`defines.proto`**：控制文件生成的选项。
+ **`http.proto`**：用于生成 HTTP 服务的选项。
+ **`struct.proto`**：通用结构体相关的选项。

**核心功能：**
+ **`cato_package`**：根据 Protobuf 的 message 生成对应的 Go 结构体。
+ **`repo_package`**：如果一个 message 代表一个数据对象，则会根据定义的键生成相应的数据仓库操作接口，包含 `FindBy`, `UpdateBy`, `DeleteBy`, `Insert` 等方法。
+ **`rdb_repo_package`**：对于拥有仓库接口的 message，生成面向关系型数据库的具体实现方法。

## 🎯 Cato 的目标

Cato 致力于长久地维护使用者 Golang 服务端项目的**稳定性**与**清晰性**。Cato 的名字灵感来源于游戏 [《CATO: Butterered Cat/黄油猫》](https://store.steampowered.com/app/1999520/CATO_Buttered_Cat/)，其核心理念是 **“猫总是脚着地”**。

对于 Golang 服务端项目而言，使用者借助 Cato 也能确保项目总能平稳“落地”——保持稳定与可控，同时让项目内部的代码拥有如猫咪般的**灵活性与流畅度**。

## 🛣️ 开发路线图

我们还有**很多很多**功能计划和改进待完成。
+ [x] 建立通用包，抽象基准`interface`
+ [ ] 支持`DDL`语句生成和`DB Migrate`
+ [ ] 支持`enum`类型生成
+ [ ] 优化`Rdb`和`Repo`的生成方式
+ [ ] 优化`generator context`全局传递
+ [ ] ...
