# Cato

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/NCUHOME/cato/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/NCUHOME/cato)](https://github.com/NCUHOME/cato/issues)
[![GitHub stars](https://img.shields.io/github/stars/NCUHOME/cato)](https://github.com/NCUHOME/cato/stargazers)

English | [中文](./README_CN.md)

Cato is a protobuf-driven code generator for Go server projects. It keeps storage metadata, HTTP routing, request description, and model generation rules close to your `.proto` files, then generates the repetitive layers around them.

> [!WARNING]
> Cato is still under active development. Breaking changes can happen, and the current version should be evaluated carefully before production use.

## Why Cato

- Define server-side structure once in protobuf and reuse it across models, repo interfaces, RDB implementations, and HTTP services.
- Generate code that is repetitive but easy to drift by hand: table metadata, CRUD skeletons, route registration, request docs, and struct tags.
- Keep manual extension points safe. Files ending in `.cato.go` are regenerated every time, while `*_custom.go`, `*_extend.go`, and `extension.go` are only created when missing.

## What Cato Generates

- Model files with struct fields, tags, table helpers, column groups, JSON conversion helpers, and time-format helpers.
- Repo interfaces and constructors backed by the generic `core/rdb.Engine`.
- Default RDB implementations that currently target the `core/rdb` abstractions and ship with an `xorm` adapter.
- HTTP service interfaces, handler registration scaffolding, and OpenAPI 2.0 `swagger.json`.
- Shared protobuf option definitions under [`proto`](./proto) and their generated Go bindings under [`generated`](./generated).

## Quick Start

### Requirements

- Go 1.24 or newer
- `protoc`
- `protoc-gen-go`

### Install the Plugin

```sh
go install github.com/ncuhome/cato/cmd/protoc-gen-cato@latest
```

Or from this repository:

```sh
make install
```

If your generated project imports Cato runtime helpers, add the module dependency there as well:

```sh
go get github.com/ncuhome/cato@latest
```

### Make the Option Protos Available

Cato's custom options are imported as `cato/proto/*.proto`, so your `protoc` include path must point to the parent directory that contains the `cato` folder.

One workable layout is:

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

With that layout, pass `-I ./third_party` to `protoc`.

### Describe Output Packages in Your Proto File

At file level, Cato currently uses these package options:

- `cato_package`: package for generated model files and HTTP service scaffolding
- `repo_package`: package for generated repo interfaces
- `rdb_repo_package`: package for generated default RDB implementations

Example:

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

### Run `protoc`

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

`--cato_opt` currently supports:

| Flag | Description |
| --- | --- |
| `ext_out_dir` | Root directory used to check whether custom files already exist. In practice this should usually match `--cato_out` and point at your Go module root. |
| `swagger_path` | Optional output path for generated OpenAPI 2.0 JSON. |
| `api_host` | Optional `host` value written into `swagger.json`. |

## Generated Files and Overwrite Rules

Running Cato can produce files like:

- `<proto>.cato.go`: shared model output for the protobuf file
- `<message>_repo.cato.go`: generated repo interface and constructor
- `<message>_rdb.cato.go`: generated default RDB implementation
- `<message>_extend.go`: model extension file, created once and preserved
- `extension.go`: repo extension interface, created once and preserved
- `handlers.cato.go`: generated HTTP handler registration
- `handlers_custom.go`: custom HTTP extension file, created once and preserved
- `api.cato.go`: generated service interface bootstrap
- `api_custom.go`: custom service extension file, created once and preserved

The intended workflow is simple: regenerate `.cato.go` files whenever your proto changes, and put handwritten logic into the preserved extension files.

## Runtime Packages

Generated code currently imports helper packages from this repository:

- [`core/rdb`](./core/rdb): generic repo engine abstraction and a built-in `xorm` implementation
- [`core/httpc`](./core/httpc): HTTP service abstraction and handler container utilities

The repository also contains [`core/param`](./core/param), a small request/response binder abstraction that can be used by surrounding application code, although the current templates do not wire it automatically.

## Repository Layout

```text
cmd/protoc-gen-cato/   protoc plugin entry
proto/                 protobuf option definitions exposed to users
generated/             generated Go code for custom protobuf options
src/                   generator implementation and option handlers
core/                  runtime helpers used by generated code
config/templates/      code generation templates
```

## Current Status and Limitations

- Follow the protobuf official style guide. Field and message naming outside the usual conventions can lead to awkward generated identifiers.
- Cato already defines DDL-related protobuf options, but migration/DDL generation is not complete yet.
- Current runtime integrations are centered on repo generation, HTTP scaffolding, and the `xorm`-backed RDB adapter.
- `swagger.json` generation is OpenAPI 2.0, not OpenAPI 3.

## Example Project

For a fuller example of how these options are used in practice, see [cato-example-bms](https://github.com/NCUHOME/cato-example-bms/tree/main/proto).

## Roadmap

- Improve RDB and repo generation.
- Support enum-oriented generation.
- Complete DDL and migration support.
- Continue simplifying generator context and extensibility.
