# Cato

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/NCUHOME/cato/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/NCUHOME/cato)](https://github.com/NCUHOME/cato/issues)
[![GitHub stars](https://img.shields.io/github/stars/NCUHOME/cato)](https://github.com/NCUHOME/cato/stargazers)

> 一个用于构建稳定、可维护的 Golang 服务端项目的强大工具。

**⚠️ 重要通知：项目处于积极开发阶段**  
此项目目前正在开发过程中，**随时可能引入不兼容的更改**。请勿在生产环境中使用。

## ✨ Cato 能做什么？

Cato 当前的能力通过 `proto-gen-*` 系列的插件来提供。

### proto-gen-cato

#### 注意 ⚠️
`proto-gen-cato` 要求你的 `.proto` 文件内容符合 [Protocol Buffers 官方风格指南](https://protobuf.dev/programming-guides/style/)。否则，生成的结构体字段名可能会难以理解。

该插件通过在不同作用域下使用 [`proto`](https://github.com/NCUHOME/cato/tree/main/proto) 目录下的自定义选项来提供不同能力。  
具体使用姿势请参照 [cato-example-bms (简易图书管理系统)](https://github.com/NCUHOME/cato-example-bms/tree/main/proto) 中的 proto 文件。

`proto` 目录下各文件的职责划分如下：

+ **`extension.proto`**：插件选项的主要使用入口。
+ **`db.proto`**：与关系型数据库相关的选项。
+ **`defines.proto`**：控制文件生成的选项。
+ **`http.proto`**：用于生成 HTTP 服务相关选项。
+ **`struct.proto`**：通用结构体相关的选项。

#### 安装
```shell
go install -v github.com/ncuhome/cato/cmd/protoc-gen-cato@latest
#### 使用
首先需要引用`proto`文件夹下的`*.proto`文件，将proto文件下载到本地。

```shell
# 注意：将 /cato/proto/path 替换为您本地的 `cato/proto` 目录路径。
protoc -I=your/project/proto/path -I=/cato/proto/path --cato_out=../ --cato_opt=ext_out_dir=../,swagger_path=swagger.json,api_host=localhost path/to/your/file.proto
```

+ `ext_out_dir`：是cato生成文件的基础目录，应与cato_out保持一致。
+ `swagger_path`：如果需要生成openapi文档，则指定为`swagger.json`文件路径。
+ `api_host`：指定swagger.json文件中的host属性。

**核心功能：**

+ 使用`db.proto`相关选项生成：
    + 数据库相关`Table`模型。
    + 数据库`Repo`相关`Find`、`Delete`、`Insert`、`Update`方法。
    + 关系型数据库`Rdb`相关方法的实现。
+ 使用`http.proto`相关选项生成：
    + `http api`接口方法。
    + `http params`定义与参数绑定。
    + `openapi2.0`接口文档。
+ 使用`struct.proto`相关选项生成：
    + 带`tag`的struct，通常用于`BO`生成。
    + 不同struct之间的`mapper`方法。

## 🎯 Cato 的目标

Cato 致力于长久地维护使用者 Golang 服务端项目的**稳定性**与**清晰性**。Cato 的名字灵感来源于游戏 [《CATO: Butterered Cat/黄油猫》](https://store.steampowered.com/app/1999520/CATO_Buttered_Cat/)，其核心理念是 **“猫总是脚着地”**。

对于 Golang 服务端项目而言，使用者借助 Cato 也能确保项目总能平稳“落地”——保持稳定与可控，同时让项目内部的代码拥有如猫咪般的**灵活性与流畅度**。

## 🛣️ 开发路线图

我们还有**很多很多**功能计划和改进待完成。

+ [x] 建立通用包，抽象基准`interface`。
+ [ ] 支持`DDL`语句生成和`DB Migrate`。
+ [ ] 支持`enum`类型生成。
+ [ ] 优化`Rdb`和`Repo`的生成方式。
+ [ ] 优化`generator context`全局传递。
+ [ ] ...
```