# Cato

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/NCUHOME/cato/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/NCUHOME/cato)](https://github.com/NCUHOME/cato/issues)
[![GitHub stars](https://img.shields.io/github/stars/NCUHOME/cato)](https://github.com/NCUHOME/cato/stargazers)

> A powerful toolkit for building stable and maintainable Golang server-side projects.

**⚠️ Important Notice: Active Development Stage**
This project is currently under active development and **may introduce breaking changes at any time**. It is not recommended for production use.

## ✨ What Can Cato Do?

Cato's current capabilities are provided through a series of `proto-gen-*` plugins.

### proto-gen-cato

#### Note ⚠️
`proto-gen-cato` requires your `.proto` file content to adhere to the [Protocol Buffers Official Style Guide](https://protobuf.dev/programming-guides/style/). Otherwise, the generated struct field names may be difficult to understand.

This plugin provides various capabilities by using custom options from the [`proto`](https://github.com/NCUHOME/cato/tree/main/proto) directory within different scopes.
For specific usage examples, please refer to the proto files in [cato-example-bms (Simple Book Management System)](https://github.com/NCUHOME/cato-example-bms/tree/main/proto).

The responsibilities of the files in the `proto` directory are divided as follows:
+ **`extension.proto`**: Main entry point for plugin options.
+ **`db.proto`**: Options related to relational databases.
+ **`defines.proto`**: Options for controlling file generation.
+ **`http.proto`**: Options for generating HTTP service-related code.
+ **`struct.proto`**: Options for general struct-related generation.

#### Install
```shell
go install -v github.com/ncuhome/cato/cmd/protoc-gen-cato@latest
```

#### Usage
First, you need to reference the `*.proto` files from the `proto` directory. Download the proto files locally.
```shell
# Note: Replace /cato/proto/path with your local path to the `cato/proto` directory.
protoc -I=your/project/proto/path -I=/cato/proto/path --cato_out=../ --cato_opt=ext_out_dir=../,swagger_path=swagger.json,api_host=localhost path/to/your/file.proto
```
+ `ext_out_dir`: The base directory for Cato's generated files, which should be consistent with `cato_out`.
+ `swagger_path`: If you need to generate OpenAPI documentation, specify the path for the `swagger.json` file.
+ `api_host`: Specifies the host attribute in the `swagger.json` file.

**Core Features:**
+ Using `db.proto` options to generate:
    + Database-related `Table` models.
    + Database `Repo` methods: `Find`, `Delete`, `Insert`, `Update`.
    + Implementation of relational database `Rdb` methods.
+ Using `http.proto` options to generate:
    + `HTTP API` interface methods.
    + `HTTP params` definitions and parameter binding.
    + `OpenAPI 2.0` interface documentation.
+ Using `struct.proto` options to generate:
    + Structs with `tags`, typically used for generating `BO` (Business Objects).
    + `Mapper` methods for converting between different structs.

## 🎯 The Goal of Cato

Cato is committed to the long-term maintenance of **stability** and **clarity** in your Golang server-side projects. The name Cato is inspired by the game [《CATO: Butterered Cat》](https://store.steampowered.com/app/1999520/CATO_Buttered_Cat/), whose core concept is **"a cat always lands on its feet."**

For Golang server-side projects, using Cato helps ensure your project always lands smoothly—maintaining stability and controllability while enabling the codebase to possess the **flexibility and fluidity** of a cat.

## 🛣️ Development Roadmap

We have **many, many** features, plans, and improvements yet to complete.
+ [x] Establish common packages to abstract foundational `interfaces`.
+ [ ] Support `DDL` statement generation and `DB Migrate`.
+ [ ] Support `enum` type generation.
+ [ ] Optimize the generation method for `Rdb` and `Repo`.
+ [ ] Optimize the global passing of the `generator context`.
+ [ ] ...