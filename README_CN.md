# Cato

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/NCUHOME/cato/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/NCUHOME/cato)](https://github.com/NCUHOME/cato/issues)
[![GitHub stars](https://img.shields.io/github/stars/NCUHOME/cato)](https://github.com/NCUHOME/cato/stargazers)

[English](./README.md) | 中文

Cato 是一个面向 Go 服务端项目的 protobuf 驱动代码生成器。你把存储映射、HTTP 路由、请求描述和结构体生成规则写在 `.proto` 里，Cato 再把这些约定扩展成模型、Repo、RDB 实现、HTTP 服务骨架和接口文档。

> [!WARNING]
> Cato 仍处于积极开发阶段。当前版本随时可能发生不兼容变更，投入生产前需要自行充分评估。

## 为什么用 Cato

- 以 protobuf 作为单一描述源，把服务端模型、存储信息和 HTTP 元数据放在一起维护。
- 自动生成那些重复、容易漂移、但又必须长期一致的代码：表字段信息、CRUD 骨架、路由注册、请求文档、结构体标签。
- 保留安全的手写扩展点。所有 `.cato.go` 文件都会在每次生成时覆盖；`*_custom.go`、`*_extend.go` 和 `extension.go` 只会在文件不存在时创建。

## Cato 会生成什么

- 带结构体字段、标签、表辅助方法、列分组、JSON 转换辅助方法和时间格式辅助方法的模型代码。
- 基于 `core/rdb.Engine` 的 Repo 接口和构造函数。
- 默认的 RDB 实现。目前运行时主要围绕 `core/rdb` 抽象，并内置了 `xorm` 适配器。
- HTTP 服务接口、处理器注册骨架，以及 OpenAPI 2.0 `swagger.json`。
- 对外暴露的 protobuf 自定义选项定义，位于 [`proto`](./proto)；对应的 Go 绑定位于 [`generated`](./generated)。

## 快速开始

### 环境要求

- Go 1.24 或更新版本
- `protoc`
- `protoc-gen-go`

### 安装插件

```sh
go install github.com/ncuhome/cato/cmd/protoc-gen-cato@latest
```

或者在本仓库内执行：

```sh
make install
```

如果你的目标项目会引用 Cato 的运行时辅助包，还需要把模块依赖加进去：

```sh
go get github.com/ncuhome/cato@latest
```

### 让 `protoc` 能找到 Cato 的选项文件

Cato 的自定义选项是通过 `cato/proto/*.proto` 引入的，所以 `protoc` 的 include path 需要指向“包含 `cato` 目录的父目录”，而不是直接指向 `proto` 目录本身。

一个可工作的目录结构示例如下：

```text
third_party/
  cato/
    proto/
      extension.proto
      db.proto
      defines.proto
      http.proto
      struct.proto
```

在这个布局下，`protoc` 需要带上 `-I ./third_party`。

### 在 Proto 文件里声明输出包

目前 Cato 在文件级会实际使用这些输出选项：

- `cato_package`：生成模型文件和 HTTP 服务骨架的目标包
- `repo_package`：生成 Repo 接口的目标包
- `rdb_repo_package`：生成默认 RDB 实现的目标包

示例：

```proto
syntax = "proto3";

package example.user.v1;
option go_package = "github.com/your-org/your-project/api/user/v1;userv1";

import "cato/proto/extension.proto";

option (cato.cato_opt) = {
  cato_package: "internal/model/user"
  repo_package: "internal/repo/user"
  rdb_repo_package: "internal/repo/user/rdb"
};

message User {
  option (cato.db_opt) = {
    db_type: CATO_DB_TYPE_MYSQL
  };
  option (cato.table_opt) = {
    name_option: { simple_name: "users" }
  };
  option (cato.struct_opt) = {
    field_default_tags: [
      {
        tag_name: "json"
        tag_value: "%s,omitempty"
        mapper: CATO_FIELD_MAPPER_SNAKE_CASE
      }
    ]
  };

  int64 id = 1 [(cato.column_opt) = {
    col_desc: {
      field_name: "id"
      comment: "primary key"
    }
    keys: [{ key_name: "PRIMARY", key_type: CATO_DB_KEY_TYPE_PRIMARY }]
  }];

  string nickname = 2 [(cato.column_opt) = {
    col_desc: {
      field_name: "nickname"
      comment: "display name"
    }
  }];
}

message CreateUserRequest {
  option (cato.http_param_opt) = {};

  string nickname = 1 [(cato.http_pf_opt) = {
    name: "nickname"
    must: true
    example: "neo"
  }];
}

message CreateUserResponse {
  User user = 1;
}

service UserService {
  option (cato.http_opt) = {
    group_prefix: "/v1/users"
    as_http_service: true
  };

  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse) {
    option (cato.router_opt) = {
      router: "/create"
      method: "POST"
    };
  }
}
```

### 执行 `protoc`

```sh
protoc \
  -I . \
  -I ./third_party \
  --go_out=. \
  --go_opt=paths=source_relative \
  --cato_out=. \
  --cato_opt=ext_out_dir=.,swagger_path=./swagger.json,api_host=localhost \
  api/user/v1/user.proto
```

目前 `--cato_opt` 支持：

| 参数 | 说明 |
| --- | --- |
| `ext_out_dir` | 用于检查自定义扩展文件是否已存在的根目录。通常应与 `--cato_out` 一致，并指向你的 Go 模块根目录。 |
| `swagger_path` | 可选，生成 OpenAPI 2.0 JSON 的输出路径。 |
| `api_host` | 可选，写入 `swagger.json` 中的 `host` 字段。 |

## 生成文件与覆盖规则

执行 Cato 后，常见输出包括：

- `<proto>.cato.go`：某个 proto 文件的共享模型输出
- `<message>_repo.cato.go`：生成的 Repo 接口与构造函数
- `<message>_rdb.cato.go`：生成的默认 RDB 实现
- `<message>_extend.go`：模型扩展文件，只创建一次，之后保留
- `extension.go`：Repo 扩展接口文件，只创建一次，之后保留
- `handlers.cato.go`：生成的 HTTP 处理器注册代码
- `handlers_custom.go`：自定义 HTTP 扩展文件，只创建一次，之后保留
- `api.cato.go`：生成的服务接口初始化代码
- `api_custom.go`：自定义服务扩展文件，只创建一次，之后保留

推荐的使用方式很直接：proto 变更后反复覆盖生成 `.cato.go`，把手写逻辑放进保留型扩展文件。

## 运行时辅助包

当前生成代码会依赖本仓库里的这些运行时包：

- [`core/rdb`](./core/rdb)：通用 Repo Engine 抽象，以及内置的 `xorm` 实现
- [`core/httpc`](./core/httpc)：HTTP 服务抽象和处理器容器工具

仓库里还提供了 [`core/param`](./core/param) 这组请求/响应绑定器抽象，但当前模板不会自动把它接进生成代码里，通常由上层应用自行接入。

## 仓库结构

```text
cmd/protoc-gen-cato/   protoc 插件入口
proto/                 对外暴露的 protobuf 选项定义
generated/             自定义 protobuf 选项对应的 Go 代码
src/                   生成器实现和选项处理逻辑
core/                  生成代码依赖的运行时辅助包
config/templates/      代码生成模板
```

## 当前状态与限制

- 请尽量遵循 protobuf 官方风格指南。不规范的字段或消息命名，容易导致生成后的 Go 标识符难看或难用。
- 仓库里已经定义了 DDL 相关的 protobuf 选项，但迁移/DDL 生成功能目前还没有完全接通。
- 当前运行时能力主要集中在 Repo 生成、HTTP 骨架和基于 `xorm` 的 RDB 适配。
- `swagger.json` 输出的是 OpenAPI 2.0，而不是 OpenAPI 3。

## 示例项目

如果你想看更完整的使用方式，可以参考 [cato-example-bms](https://github.com/NCUHOME/cato-example-bms/tree/main/proto)。

## 路线图

- 继续优化 RDB 和 Repo 的生成方式
- 支持更多面向 enum 的生成能力
- 完成 DDL 与迁移支持
- 继续简化生成上下文与扩展机制
